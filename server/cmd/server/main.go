package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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
	// 解析命令行参数
	serverConfig := flag.String("config", "config.json", "服务器配置文件路径")
	clientConfig := flag.String("rai", "", "客户端配置文件路径或目录 (.rai)")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*serverConfig, *clientConfig)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 验证配置
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	// 设置 Gin 模式
	if cfg.Server.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 引擎
	router := gin.Default()

	// 添加日志中间件
	router.Use(logger.Middleware(cfg))

	// 创建代理服务
	proxyService := services.NewProxyService()

	// 创建 API 处理器
	apiHandler := handlers.NewAPIHandler(proxyService)

	// 健康检查路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"time":    time.Now().Format(time.RFC3339),
			"version": "1.0.0",
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

		// 添加限流中间件（在认证之前）
		api.Use(middleware.RateLimit(globalLimiter, ipLimiter))

		// 添加认证中间件
		api.Use(middleware.TokenAuth(cfg))

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

	log.Printf("Server starting on %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
