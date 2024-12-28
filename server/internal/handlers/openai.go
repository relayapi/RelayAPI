package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"relayapi/server/internal/services"
	"relayapi/server/internal/models"
)

const (
	OpenAIBaseURL = "https://api.openai.com/v1"
	DashScopeBaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
)

// APIHandler 处理 API 请求
type APIHandler struct {
	proxyService *services.ProxyService
}

// NewAPIHandler 创建新的 API 处理器
func NewAPIHandler(proxyService *services.ProxyService) *APIHandler {
	return &APIHandler{
		proxyService: proxyService,
	}
}

// getBaseURL 根据提供者获取基础 URL
func (h *APIHandler) getBaseURL(provider string) string {
	switch provider {
	case "openai":
		return OpenAIBaseURL
	case "dashscope":
		return DashScopeBaseURL
	default:
		return DashScopeBaseURL // 默认使用 DashScope
	}
}

// HandleRequest 处理 API 请求
func (h *APIHandler) HandleRequest(c *gin.Context) {
	// 获取请求路径
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid API path",
		})
		return
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to read request body: %v", err),
		})
		return
	}

	// 从上下文中获取令牌和提供者信息
	token, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Token not found in context",
		})
		return
	}

	// 获取令牌中的 API Key 和提供者
	tokenObj := token.(*models.Token)
	apiKey := tokenObj.APIKey
	provider := tokenObj.Provider

	// 构建目标 URL
	baseURL := h.getBaseURL(provider)
	
	// 处理路径
	// 移除开头的斜杠
	path = strings.TrimPrefix(path, "/")
	// 如果路径包含版本号，则移除
	path = strings.TrimPrefix(path, "v1/")
	
	targetURL := fmt.Sprintf("%s/%s", baseURL, path)

	// 转发请求
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", apiKey),
	}

	fmt.Printf("Provider: %s, Target URL: %s\n", provider, targetURL)

	resp, err := h.proxyService.ProxyRequest(c.Request.Method, targetURL, headers, body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to proxy request: %v", err),
		})
		return
	}

	// 读取响应
	respBody, err := h.proxyService.ReadResponse(resp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to read response: %v", err),
		})
		return
	}

	// 设置响应头
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// 返回响应
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
} 