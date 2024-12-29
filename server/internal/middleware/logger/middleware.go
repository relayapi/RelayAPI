package logger

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"strings"
	"time"

	"relayapi/server/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// bodyLogWriter 响应体写入器
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// decompressBody 根据Content-Encoding解压缩body
func decompressBody(encoding string, body []byte) ([]byte, error) {
	switch strings.ToLower(encoding) {
	case "gzip":
		reader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		return io.ReadAll(reader)
	case "deflate":
		reader, err := zlib.NewReader(bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		return io.ReadAll(reader)
	default:
		return body, nil
	}
}

// isTextContent 判断内容类型是否为文本
func isTextContent(contentType string) bool {
	if contentType == "" {
		return true
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return true
	}
	return strings.HasPrefix(mediaType, "text/") ||
		strings.HasPrefix(mediaType, "application/json") ||
		strings.HasPrefix(mediaType, "application/xml") ||
		strings.HasPrefix(mediaType, "application/x-www-form-urlencoded")
}

// 获取响应体内容
func getResponseBody(c *gin.Context, blw *bodyLogWriter) string {
	responseBody := blw.body.String()

	// 如果是 SSE 响应，使用特定的格式化
	if strings.Contains(c.Writer.Header().Get("Content-Type"), "text/event-stream") {
		return formatSSEResponse(responseBody)
	}

	// 处理压缩的响应
	if encoding := c.Writer.Header().Get("Content-Encoding"); encoding != "" {
		if decompressed, err := decompressBody(encoding, []byte(responseBody)); err == nil {
			responseBody = string(decompressed)
		}
	}

	// 判断是否需要 base64 编码
	contentType := c.Writer.Header().Get("Content-Type")
	if !isTextContent(contentType) {
		return base64.StdEncoding.EncodeToString([]byte(responseBody))
	}

	return responseBody
}

// Middleware 创建日志中间件
func Middleware(cfg *config.Config) gin.HandlerFunc {
	var writers []LogWriter

	// 根据配置创建日志写入器
	if cfg.Server.Log.Console {
		writers = append(writers, NewAsyncLogWriter(NewConsoleLogWriter(), 1000))
	}

	if cfg.Server.Log.Database.Enabled {
		if dbWriter, err := NewDatabaseLogWriter(
			cfg.Server.Log.Database.Type,
			cfg.Server.Log.Database.ConnectionString,
		); err == nil {
			writers = append(writers, NewAsyncLogWriter(dbWriter, 1000))
		} else {
			fmt.Printf("Failed to create database log writer: %v\n", err)
		}
	}

	if cfg.Server.Log.Web.Enabled {
		writers = append(writers, NewAsyncLogWriter(NewWebLogWriter(cfg.Server.Log.Web.CallbackURL), 1000))
	}

	if cfg.Server.Log.Parquet.Enabled {
		if parquetWriter, err := NewParquetLogWriter(cfg.Server.Log.Parquet.FilePath); err == nil {
			writers = append(writers, NewAsyncLogWriter(parquetWriter, 1000))
		} else {
			fmt.Printf("Failed to create parquet log writer: %v\n", err)
		}
	}

	return func(c *gin.Context) {
		// 生成请求ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// 记录开始时间
		startTime := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			// 处理压缩的内容
			if encoding := c.Request.Header.Get("Content-Encoding"); encoding != "" {
				if decompressed, err := decompressBody(encoding, requestBody); err == nil {
					requestBody = decompressed
				}
			}
		}

		// 判断是否需要base64编码
		var requestBodyStr string
		contentType := c.Request.Header.Get("Content-Type")
		if isTextContent(contentType) {
			requestBodyStr = string(requestBody)
		} else {
			requestBodyStr = base64.StdEncoding.EncodeToString(requestBody)
		}

		// 记录请求日志
		requestLog := map[string]interface{}{
			"request_id":   requestID,
			"type":         "request",
			"time":         startTime.Format(time.RFC3339),
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"client_ip":    c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"request_body": requestBodyStr,
			"headers":      c.Request.Header,
		}

		// 写入请求日志到所有写入器
		for _, writer := range writers {
			if err := writer.Write(requestLog); err != nil {
				fmt.Printf("Failed to write request log: %v\n", err)
			}
		}

		// 包装响应写入器以捕获响应体
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 记录结束时间
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 获取并处理响应体
		responseBody := getResponseBody(c, blw)

		// 记录响应日志
		responseLog := map[string]interface{}{
			"request_id":    requestID,
			"type":          "response",
			"time":          endTime.Format(time.RFC3339),
			"status":        c.Writer.Status(),
			"latency_ms":    latency.Milliseconds(),
			"response_body": responseBody,
			"headers":       c.Writer.Header(),
			"errors":        c.Errors.Errors(),
		}

		// 写入响应日志到所有写入器
		for _, writer := range writers {
			if err := writer.Write(responseLog); err != nil {
				fmt.Printf("Failed to write response log: %v\n", err)
			}
		}
	}
}

// formatSSEResponse 格式化SSE响应数据
func formatSSEResponse(response string) string {
	lines := strings.Split(response, "\n")
	var content strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				continue
			}

			var event struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &event); err == nil {
				if len(event.Choices) > 0 && event.Choices[0].Delta.Content != "" {
					content.WriteString(event.Choices[0].Delta.Content)
				}
			} else {
				// JSON 解析失败时，保留原始数据
				content.WriteString(data)
			}
		}
	}

	return content.String()
}
