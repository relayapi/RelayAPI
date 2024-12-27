package middleware

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"relayapi/server/internal/config"
	"relayapi/server/internal/crypto"
	"relayapi/server/internal/models"
)

// TokenAuth 验证访问令牌的中间件
func TokenAuth(cfg *config.Config) gin.HandlerFunc {
	// 创建加密器
	encryptor, err := crypto.NewEncryptor(cfg)
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		// 从 URL 参数中获取令牌
		encryptedToken := c.Query("token")
		if encryptedToken == "" {
			// 尝试从 URL 路径中获取令牌（兼容某些 API 的路径格式）
			encryptedToken = c.Param("token")
		}

		if encryptedToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing API token",
				"message": "Please provide your API token as a URL parameter: ?token=your_token",
			})
			c.Abort()
			return
		}

		// 添加回 base64 padding
		if padding := len(encryptedToken) % 4; padding > 0 {
			encryptedToken += strings.Repeat("=", 4-padding)
		}

		// Base64 URL 安全解码
		tokenBytes, err := base64.URLEncoding.DecodeString(encryptedToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid token format",
				"message": "Token must be base64url encoded",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// 解密令牌
		decryptedBytes, err := encryptor.Decrypt(tokenBytes)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
				"message": "Failed to decrypt token",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// 反序列化令牌
		token := &models.Token{}
		if err := token.Deserialize(decryptedBytes); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
				"message": "Failed to parse token data",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// 验证令牌有效性
		if !token.IsValid() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token expired or exceeded usage limit",
				"message": "Please obtain a new token",
			})
			c.Abort()
			return
		}

		// 增加使用次数
		token.IncrementUsage()

		// TODO: 更新数据库中的令牌使用次数

		// 将令牌和 API Key 存储在上下文中
		c.Set("token", token)
		c.Set("api_token", token.APIKey)

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