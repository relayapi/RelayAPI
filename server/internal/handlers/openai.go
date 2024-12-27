package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"relayapi/server/internal/services"
)

const (
	OpenAIBaseURL = "https://api.openai.com/v1"
)

// OpenAIHandler 处理 OpenAI API 请求
type OpenAIHandler struct {
	proxyService *services.ProxyService
}

// NewOpenAIHandler 创建新的 OpenAI 处理器
func NewOpenAIHandler(proxyService *services.ProxyService) *OpenAIHandler {
	return &OpenAIHandler{
		proxyService: proxyService,
	}
}

// HandleRequest 处理 OpenAI API 请求
func (h *OpenAIHandler) HandleRequest(c *gin.Context) {
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

	// 构建目标 URL
	targetURL := OpenAIBaseURL + path

	// 转发请求
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": c.GetHeader("Authorization"),
	}

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