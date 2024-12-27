package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"relayapi/server/internal/config"
	"relayapi/server/internal/handlers"
	"relayapi/server/internal/middleware"
	"relayapi/server/internal/services"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 创建服务
	proxyService := services.NewProxyService()
	openaiHandler := handlers.NewOpenAIHandler(proxyService)

	// 创建路由
	r := gin.Default()

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API 代理路由组
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.TokenAuth(cfg))
	apiGroup.Use(middleware.RateLimit())
	{
		// OpenAI API 代理
		apiGroup.Any("/openai/*path", openaiHandler.HandleRequest)
	}

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 