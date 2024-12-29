package services

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ProxyService 处理 API 代理请求
type ProxyService struct {
	client *http.Client
}

// NewProxyService 创建新的代理服务
func NewProxyService() *ProxyService {
	return &ProxyService{
		client: &http.Client{},
	}
}

// ProxyRequest 转发 API 请求
func (s *ProxyService) ProxyRequest(method, url string, headers map[string]string, body []byte) (*http.Response, error) {
	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ReadResponse 读取响应内容
func (s *ProxyService) ReadResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// HandleStreamResponse 处理流式响应
func (s *ProxyService) HandleStreamResponse(c *gin.Context, resp *http.Response) error {
	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 获取原始响应的 Content-Type
	contentType := resp.Header.Get("Content-Type")
	isEventStream := strings.Contains(contentType, "text/event-stream")

	// 创建一个 reader
	reader := bufio.NewReader(resp.Body)
	defer resp.Body.Close()

	// 刷新写入器以确保头信息被发送
	c.Writer.Flush()

	for {
		// 读取一行数据
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		// 如果不是 SSE 格式，转换为 SSE 格式
		if !isEventStream && len(line) > 0 {
			line = []byte("data: " + string(line) + "\n\n")
		}

		// 写入数据
		_, err = c.Writer.Write(line)
		if err != nil {
			return err
		}

		// 刷新写入器
		c.Writer.Flush()
	}
}
