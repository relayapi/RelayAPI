package main

import (
	"context"
	"flag"
	"fmt"
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

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func main() {
	// 禁用 Gin 的默认日志输出
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = os.Stdout
	log.SetFlags(0) // 移除默认的日志前缀

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

	// 解析命令行参数
	serverConfig := flag.String("config", "config.json", "服务器配置文件路径")
	clientConfig := flag.String("rai", "", "客户端配置文件路径或目录 (.rai)")
	flag.Parse()

	log.Printf("\033[36m📚 Loading configuration files %s\033[0m", *serverConfig)
	log.Printf("\033[36m📚 Loading rai file: %s\033[0m", *clientConfig)

	// 加载配置
	cfg, err := config.LoadConfig(*serverConfig, *clientConfig)
	if err != nil {
		log.Fatalf("\033[31m❌ Failed to load config: %v\033[0m", err)
	}
	log.Println("\033[32m✅ Configuration loaded successfully\033[0m")

	// 验证配置
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("\033[31m❌ Invalid config: %v\033[0m", err)
	}
	log.Println("\033[32m✅ Configuration validated\033[0m")

	// 设置 Gin 模式
	if cfg.Server.Server.Debug {
		gin.SetMode(gin.DebugMode)
		log.Println("\033[33m🔧 Running in DEBUG mode\033[0m")
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Println("\033[32m🔧 Running in RELEASE mode\033[0m")
	}

	// 创建 Gin 引擎
	router := gin.New()

	// 创建统计服务
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Server.Host, cfg.Server.Server.Port)
	statsService := services.NewStats("1.0.0", serverAddr)

	// 启动统计信息显示
	go statsService.StartConsoleDisplay(stopChan)

	log.Println("\033[36m🔧 Initializing middleware...\033[0m")
	// 添加路径规范化中间件
	router.Use(middleware.PathNormalizationMiddleware())
	log.Println("\033[32m✅ Path normalization middleware initialized\033[0m")

	// 添加日志中间件
	router.Use(logger.Middleware(cfg))
	log.Println("\033[32m✅ Logger middleware initialized\033[0m")

	// 创建代理服务
	proxyService := services.NewProxyService()
	log.Println("\033[32m✅ Proxy service initialized\033[0m")

	// 创建 API 处理器
	apiHandler := handlers.NewAPIHandler(proxyService)
	log.Println("\033[32m✅ API handler initialized\033[0m")

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
		log.Println("\033[36m🔧 Configuring rate limiters...\033[0m")
		// 创建全局限流器和 IP 限流器
		globalLimiter := rate.NewLimiter(rate.Limit(cfg.Server.RateLimit.RequestsPerSecond), cfg.Server.RateLimit.Burst)
		ipLimiter := middleware.NewIPRateLimiter(
			rate.Limit(cfg.Server.RateLimit.IPLimit.RequestsPerSecond),
			cfg.Server.RateLimit.IPLimit.Burst,
		)
		log.Println("\033[32m✅ Rate limiters configured\033[0m")

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
		log.Println("\033[32m✅ Rate limit middleware initialized\033[0m")

		// 添加认证中间件
		api.Use(middleware.TokenAuth(cfg))
		log.Println("\033[32m✅ Authentication middleware initialized\033[0m")

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
		log.Printf("\033[36m🚀 Server starting on %s\033[0m", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("\033[31m❌ Failed to start server: %v\033[0m", err)
		}
	}()

	// 等待中断信号
	<-sigChan
	log.Println("\n\033[33m⚡ Shutting down server...\033[0m")

	// 关闭统计显示
	close(stopChan)

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("\033[31m❌ Server forced to shutdown: %v\033[0m", err)
	}

	log.Println("\033[32m✅ Server stopped gracefully\033[0m")
}
