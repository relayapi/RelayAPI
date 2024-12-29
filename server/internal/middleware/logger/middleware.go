package logger

import (
	"bytes"
	"fmt"
	"io"
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
			"request_body": string(requestBody),
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

		// 获取响应体
		responseBody := blw.body.String()
		if strings.Contains(c.Writer.Header().Get("Content-Type"), "text/event-stream") {
			responseBody = formatSSEResponse(responseBody)
		}

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
	var events []string
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			events = append(events, strings.TrimPrefix(line, "data: "))
		}
	}
	return "[" + strings.Join(events, ",") + "]"
}
