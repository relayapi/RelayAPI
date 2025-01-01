package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"relayapi/server/internal/config"
	"relayapi/server/internal/handlers"
	"relayapi/server/internal/middleware"
	"relayapi/server/internal/middleware/logger"
	"relayapi/server/internal/services"
	"relayapi/server/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var debugMode bool
var logWriter io.Writer

func setupLogging(debug bool) {
	if debug {
		// 创建或打开debug.log文件
		logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("❌ Failed to open debug.log: %v", err)
		}
		// 设置日志输出到文件，并添加时间戳
		logWriter = logFile
		log.SetOutput(logWriter)
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
		// 设置 Gin 的日志输出到同一个文件
		gin.DefaultWriter = logWriter
		gin.SetMode(gin.DebugMode)
	} else {
		// 如果不是debug模式，禁用所有日志输出
		logWriter = io.Discard
		log.SetOutput(logWriter)
		log.SetFlags(0)
		gin.DefaultWriter = io.Discard
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	// 解析命令行参数
	var genConfig string
	serverConfig := flag.String("config", "config.json", "服务器配置文件路径")
	clientConfig := flag.String("rai", "default.rai", "客户端配置文件路径或目录 (.rai)")
	flag.StringVar(&genConfig, "gen", "", "生成客户端配置 (格式: [host:port] 或 help)")
	flag.BoolVar(&debugMode, "debug", false, "启用调试日志输出到debug.log")
	flag.BoolVar(&debugMode, "d", false, "启用调试日志输出到debug.log (简写)")

	// 自定义 Usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n使用 --gen help 查看配置生成的详细说明\n")
	}

	flag.Parse()

	// 检查是否有 --gen 标志
	if isFlagPassed("gen") {
		utils.OnceCMDGenerateClientConfig(genConfig)
	}

	// 设置日志
	setupLogging(debugMode)

	// 创建停止通道
	stopChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 渐变色数组
	gradientColors := []string{
		"\033[38;5;51m", // 浅青色
		"\033[38;5;45m", // 青色
		"\033[38;5;39m", // 深青色
		"\033[38;5;33m", // 蓝色
		"\033[38;5;27m", // 深蓝色
	}

	// 打印启动标题
	title := "=== RelayAPI Server Starting ==="
	colorIdx := 0
	for _, char := range title {
		fmt.Print(gradientColors[colorIdx%len(gradientColors)], string(char))
		colorIdx++
	}
	fmt.Print("\033[0m\n\n")

	log.Printf("📚 Loading configuration files %s", *serverConfig)
	log.Printf("📚 Loading rai file: %s", *clientConfig)

	// 加载配置
	cfg, err := config.LoadConfig(*serverConfig, *clientConfig)
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// 验证配置
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("❌ Invalid config: %v", err)
	}

	// 记录运行模式
	if debugMode {
		log.Println("🔧 Running in DEBUG mode")
	} else {
		log.Println("🔧 Running in RELEASE mode")
	}

	// 创建 Gin 引擎
	router := gin.New()

	// 创建统计服务
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Server.Host, cfg.Server.Server.Port)
	statsService := services.NewStats("1.0.0", serverAddr)

	// 启动统计信息显示
	go statsService.StartConsoleDisplay(stopChan)

	log.Println("🔧 Initializing middleware...")
	// 添加路径规范化中间件
	router.Use(middleware.PathNormalizationMiddleware())

	// 添加日志中间件
	router.Use(logger.Middleware(cfg))

	// 创建代理服务
	proxyService := services.NewProxyService()

	// 创建 API 处理器
	apiHandler := handlers.NewAPIHandler(proxyService)

	// 健康检查路由
	router.GET("/health", func(c *gin.Context) {
		uptime := statsService.GetUptime()
		totalReqs := atomic.LoadUint64(&statsService.TotalRequests)
		stats := map[string]interface{}{
			"uptime":              uptime.String(),
			"total_requests":      totalReqs,
			"successful_requests": atomic.LoadUint64(&statsService.SuccessfulRequests),
			"failed_requests":     atomic.LoadUint64(&statsService.FailedRequests),
			"bytes_received":      atomic.LoadUint64(&statsService.BytesReceived),
			"bytes_sent":          atomic.LoadUint64(&statsService.BytesSent),
			"tps":                 float64(totalReqs) / uptime.Seconds(),
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"time":    time.Now().Format(time.RFC3339),
			"version": "1.0.0",
			"stats":   stats,
		})
	})

	// API 路由组
	api := router.Group("/relayapi")
	{
		// 创建全局限流器和 IP 限流器
		globalLimiter := rate.NewLimiter(rate.Limit(cfg.Server.RateLimit.RequestsPerSecond), cfg.Server.RateLimit.Burst)
		ipLimiter := middleware.NewIPRateLimiter(
			rate.Limit(cfg.Server.RateLimit.IPLimit.RequestsPerSecond),
			cfg.Server.RateLimit.IPLimit.Burst,
		)

		// 添加统计中间件
		api.Use(func(c *gin.Context) {
			statsService.IncrementTotal()
			c.Next()
			status := c.Writer.Status()
			if status >= 400 {
				statsService.IncrementFailed()
				statsService.IncrementErrorStatus(status)
			} else {
				statsService.IncrementSuccess()
			}
			// 记录请求和响应大小
			statsService.AddBytesReceived(uint64(c.Request.ContentLength))
			statsService.AddBytesSent(uint64(c.Writer.Size()))
		})

		// 添加限流中间件（在认证之前）
		api.Use(middleware.RateLimit(globalLimiter, ipLimiter))

		// 添加认证中间件
		api.Use(middleware.TokenAuth(cfg))

		// 所有 API 请求通过统一入口处理
		api.Any("/*path", apiHandler.HandleRequest)
	}

	// 启动服务器
	server := &http.Server{
		Addr:           serverAddr,
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.Server.Server.MaxHeaderBytes,
	}

	// 在新的 goroutine 中启动服务器
	go func() {
		log.Printf("🚀 Server starting on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	<-sigChan
	log.Println("\n⚡ Shutting down server...")

	// 关闭统计显示
	close(stopChan)

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("❌ Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Server stopped gracefully")
}

// isFlagPassed 检查命令行参数是否被传递
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
