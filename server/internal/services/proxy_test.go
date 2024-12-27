package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyRequest(t *testing.T) {
	// 创建测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求头
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Content-Type header not set correctly")
		}
		if r.Header.Get("X-Test-Header") != "test-value" {
			t.Error("Custom header not set correctly")
		}

		// 返回测试响应
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	// 创建代理服务
	proxyService := NewProxyService()

	// 测试请求
	headers := map[string]string{
		"Content-Type":   "application/json",
		"X-Test-Header": "test-value",
	}
	body := []byte(`{"test":"data"}`)

	resp, err := proxyService.ProxyRequest("POST", ts.URL, headers, body)
	if err != nil {
		t.Fatalf("ProxyRequest failed: %v", err)
	}

	// 验证响应
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// 读取响应内容
	respBody, err := proxyService.ReadResponse(resp)
	if err != nil {
		t.Fatalf("ReadResponse failed: %v", err)
	}

	expectedBody := `{"status":"ok"}`
	if string(respBody) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, string(respBody))
	}
}

func TestProxyRequestError(t *testing.T) {
	proxyService := NewProxyService()

	// 测试无效 URL
	_, err := proxyService.ProxyRequest("GET", "invalid-url", nil, nil)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
} 