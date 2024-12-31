package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"relayapi/server/internal/config"
	"relayapi/server/internal/crypto"
	"relayapi/server/internal/models"

	"github.com/gin-gonic/gin"
)

func splitStringByFirstSlash(input string) (string, string) {
	index := strings.Index(input, "/")
	if index == -1 {
		return input, "" // 如果没有找到斜杠，返回原字符串和空字符串
	}
	return input[:index], input[index+1:]
}

// TokenAuth 验证访问令牌的中间件
func TokenAuth(cfg *config.Config) gin.HandlerFunc {
	// 创建加密器映射
	encryptors := make(map[string]crypto.Encryptor)

	return func(c *gin.Context) {
		// 从 URL 参数中获取令牌
		encryptedToken := c.Query("token")
		if encryptedToken == "" {
			// 尝试从 URL 路径中获取令牌（兼容某些 API 的路径格式）
			encryptedToken = c.Param("token")
		}

		if encryptedToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Missing API token",
				"message": "Please provide your API token as a URL parameter: ?token=your_token",
			})
			c.Abort()
			return
		}
		encryptedToken, extPath := splitStringByFirstSlash(encryptedToken)
		if extPath != "" {
			c.Set("ext_path", strings.TrimRight(extPath, "="))
		}

		// 获取配置 hash
		raiHash := c.Query("rai_hash")
		if raiHash == "" {
			raiHash = c.Param("rai_hash")
		}

		raiHash, extPath = splitStringByFirstSlash(raiHash)
		if extPath != "" {
			c.Set("ext_path", strings.TrimRight(extPath, "="))
		}

		if raiHash == "" {
			// 如果没有指定 hash，使用第一个可用的配置
			for hash := range cfg.Clients {
				raiHash = hash
				break
			}
		}
		log.Println("raiHash", raiHash)
		log.Println("cfg.Clients", cfg.Clients)
		clientCfg, ok := cfg.GetClientConfig(raiHash)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid configuration hash",
				"message": "The provided configuration hash is not valid",
			})
			c.Abort()
			return
		}

		// 获取或创建加密器
		encryptor, ok := encryptors[raiHash]
		if !ok {
			var err error
			encryptor, err = crypto.NewEncryptor(&clientCfg)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Encryptor initialization failed",
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			encryptors[raiHash] = encryptor
		}

		// 清理令牌字符串
		encryptedToken = strings.TrimSpace(encryptedToken)

		// 添加回 base64 padding
		if padding := len(encryptedToken) % 4; padding > 0 {
			encryptedToken += strings.Repeat("=", 4-padding)
		}

		// Base64 URL 安全解码
		tokenBytes, err := base64.URLEncoding.DecodeString(encryptedToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":        "Invalid token format",
				"message":      "Token must be base64url encoded",
				"details":      err.Error(),
				"token_length": len(encryptedToken),
				"token_start":  encryptedToken[:10],
			})
			c.Abort()
			return
		}

		// 解密令牌
		decryptedBytes, err := encryptor.Decrypt(tokenBytes)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
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
				"error":   "Invalid token",
				"message": "Failed to parse token data",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// 验证令牌有效性
		if !token.IsValid() {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Token expired or exceeded usage limit",
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
		c.Set("api_key", token.APIKey)
		c.Set("rai_hash", raiHash)

		c.Next()
	}
}
