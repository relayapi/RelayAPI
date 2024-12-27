package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TokenAuth 验证访问令牌的中间件
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 URL 参数中获取令牌
		token := c.Query("token")
		if token == "" {
			// 尝试从 URL 路径中获取令牌（兼容某些 API 的路径格式）
			token = c.Param("token")
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing API token",
				"message": "Please provide your API token as a URL parameter: ?token=your_token",
			})
			c.Abort()
			return
		}

		// 将令牌存储在上下文中，供后续处理器使用
		c.Set("api_token", token)

		// TODO: 验证令牌
		// 1. 从数据库获取令牌信息
		// 2. 验证令牌有效性
		// 3. 更新令牌使用次数

		c.Next()
	}
}

// RateLimit 限制请求频率的中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现请求频率限制
		c.Next()
	}
} 