package services

import (
	"bytes"
	"io"
	"net/http"
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