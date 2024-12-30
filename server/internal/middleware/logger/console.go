package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	// 全局日志缓冲区
	logBuffer     []string
	logBufferSize = 4 // 保留最近的100条日志
	logBufferMu   sync.RWMutex
	// 日志更新通知通道
	logUpdateChan = make(chan struct{}, 1)
)

// GetRecentLogs 获取最近的日志
func GetRecentLogs() string {
	logBufferMu.RLock()
	defer logBufferMu.RUnlock()
	return strings.Join(logBuffer, "\n")
}

// GetLogUpdateChan 获取日志更新通知通道
func GetLogUpdateChan() chan struct{} {
	return logUpdateChan
}

// ConsoleLogWriter 控制台日志写入器
type ConsoleLogWriter struct {
	mu sync.Mutex
}

// NewConsoleLogWriter 创建控制台日志写入器
func NewConsoleLogWriter() *ConsoleLogWriter {
	return &ConsoleLogWriter{}
}

func (w *ConsoleLogWriter) Write(log map[string]interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 获取日志类型和时间
	logType, _ := log["type"].(string)
	timeStr, _ := log["time"].(string)

	// 解析时间
	t, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		timeStr = t.Format("15:04:05.000")
	}

	// 构建日志行
	var logLine string

	// 如果是请求日志
	if logType == "request" {
		method, _ := log["method"].(string)
		path, _ := log["path"].(string)
		logLine = fmt.Sprintf("%s [%s] %s %s",
			timeStr,
			strings.ToUpper(logType),
			method,
			path)
	}

	// 如果是响应日志
	if logType == "response" {
		status, _ := log["status"].(float64)
		latency, _ := log["latency_ms"].(float64)
		logLine = fmt.Sprintf("%s [%s] Status: %d, Latency: %.2fms",
			timeStr,
			strings.ToUpper(logType),
			int(status),
			latency)
	}

	// 添加到日志缓冲区
	logBufferMu.Lock()
	logBuffer = append(logBuffer, logLine)
	if len(logBuffer) > logBufferSize {
		logBuffer = logBuffer[1:]
	}
	logBufferMu.Unlock()

	// 发送更新通知
	select {
	case logUpdateChan <- struct{}{}:
	default:
		// 如果通道已满，跳过通知
	}

	return nil
}

func (w *ConsoleLogWriter) Close() error {
	return nil
}
