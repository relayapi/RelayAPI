package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"relayapi/server/internal/config"
	"relayapi/server/internal/handlers"
	"relayapi/server/internal/middleware"
	"relayapi/server/internal/services"
)

func main() {
	// 解析命令行参数
	serverConfig := flag.String("config", "config.json", "服务器配置文件路径")
	clientConfig := flag.String("rai", "", "客户端配置文件路径 (.rai)")
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
	if cfg.Server.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 引擎
	router := gin.Default()

	// 创建代理服务
	proxyService := services.NewProxyService()

	// 创建 OpenAI 处理器
	openaiHandler := handlers.NewOpenAIHandler(proxyService)

	// 添加中间件
	router.Use(middleware.TokenAuth(&cfg.Client))

	// 注册路由
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
			"version": cfg.Client.Version,
		})
	})

	// OpenAI API 路由
	router.Any("/api/openai/v1/*path", openaiHandler.HandleRequest)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Server.Host, cfg.Server.Server.Port)
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.Server.Server.MaxHeaderBytes,
	}

	log.Printf("Server version %s starting on %s", cfg.Client.Version, addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 