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
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	// 创建停止通道
	stopChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("=== RelayAPI Server Starting ===")

	// 解析命令行参数
	serverConfig := flag.String("config", "config.json", "服务器配置文件路径")
	clientConfig := flag.String("rai", "", "客户端配置文件路径或目录 (.rai)")
	flag.Parse()

	log.Println("📚 Loading configuration files...")
	// 加载配置
	cfg, err := config.LoadConfig(*serverConfig, *clientConfig)
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}
	log.Println("✅ Configuration loaded successfully")

	// 验证配置
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("❌ Invalid config: %v", err)
	}
	log.Println("✅ Configuration validated")

	// 设置 Gin 模式
	if cfg.Server.Server.Debug {
		gin.SetMode(gin.DebugMode)
		log.Println("🔧 Running in DEBUG mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Println("🔧 Running in RELEASE mode")
	}

	// 创建 Gin 引擎
	router := gin.New()

	// 创建统计服务
	statsService := services.NewStats()

	// 启动统计信息显示
	go statsService.StartConsoleDisplay(stopChan)

	log.Println("🔧 Initializing middleware...")
	// 添加路径规范化中间件
	router.Use(middleware.PathNormalizationMiddleware())
	log.Println("✅ Path normalization middleware initialized")

	// 添加日志中间件
	router.Use(logger.Middleware(cfg))
	log.Println("✅ Logger middleware initialized")

	// 创建代理服务
	proxyService := services.NewProxyService()
	log.Println("✅ Proxy service initialized")

	// 创建 API 处理器
	apiHandler := handlers.NewAPIHandler(proxyService)
	log.Println("✅ API handler initialized")

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
		log.Println("🔧 Configuring rate limiters...")
		// 创建全局限流器和 IP 限流器
		globalLimiter := rate.NewLimiter(rate.Limit(cfg.Server.RateLimit.RequestsPerSecond), cfg.Server.RateLimit.Burst)
		ipLimiter := middleware.NewIPRateLimiter(
			rate.Limit(cfg.Server.RateLimit.IPLimit.RequestsPerSecond),
			cfg.Server.RateLimit.IPLimit.Burst,
		)
		log.Println("✅ Rate limiters configured")

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
		log.Println("✅ Rate limit middleware initialized")

		// 添加认证中间件
		api.Use(middleware.TokenAuth(cfg))
		log.Println("✅ Authentication middleware initialized")

		// 所有 API 请求通过统一入口处理
		api.Any("/*path", apiHandler.HandleRequest)
	}

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Server.Host, cfg.Server.Server.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.Server.Server.MaxHeaderBytes,
	}

	// 在新的 goroutine 中启动服务器
	go func() {
		log.Printf("🚀 Server starting on %s", addr)
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}
