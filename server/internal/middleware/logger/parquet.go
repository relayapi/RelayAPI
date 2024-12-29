package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/source"
	"github.com/xitongsys/parquet-go/writer"
)

// LogRecord Parquet日志记录结构
type LogRecord struct {
	RequestID   string `parquet:"name=request_id, type=BYTE_ARRAY, convertedtype=UTF8"`
	Type        string `parquet:"name=type, type=BYTE_ARRAY, convertedtype=UTF8"`
	Time        string `parquet:"name=time, type=BYTE_ARRAY, convertedtype=UTF8"`
	Method      string `parquet:"name=method, type=BYTE_ARRAY, convertedtype=UTF8"`
	Path        string `parquet:"name=path, type=BYTE_ARRAY, convertedtype=UTF8"`
	Query       string `parquet:"name=query, type=BYTE_ARRAY, convertedtype=UTF8"`
	ClientIP    string `parquet:"name=client_ip, type=BYTE_ARRAY, convertedtype=UTF8"`
	UserAgent   string `parquet:"name=user_agent, type=BYTE_ARRAY, convertedtype=UTF8"`
	RequestBody string `parquet:"name=request_body, type=BYTE_ARRAY, convertedtype=UTF8"`
	Status      int32  `parquet:"name=status, type=INT32"`
	LatencyMS   int64  `parquet:"name=latency_ms, type=INT64"`
	Errors      string `parquet:"name=errors, type=BYTE_ARRAY, convertedtype=UTF8"`
}

// ParquetLogWriter Parquet文件日志写入器
type ParquetLogWriter struct {
	baseDir     string
	currentFile string
	currentDate string
	fw          source.ParquetFile
	pw          *writer.ParquetWriter
	mutex       sync.Mutex
}

// NewParquetLogWriter 创建支持日期轮转的Parquet日志写入器
func NewParquetLogWriter(baseDir string) (*ParquetLogWriter, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}

	w := &ParquetLogWriter{
		baseDir: baseDir,
	}
	if err := w.rotateIfNeeded(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *ParquetLogWriter) rotateIfNeeded() error {
	currentDate := time.Now().Format("2006-01-02")
	if currentDate == w.currentDate && w.pw != nil {
		return nil
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 关闭当前文件
	if w.pw != nil {
		w.pw.WriteStop()
		w.fw.Close()
	}

	// 创建新文件
	filename := filepath.Join(w.baseDir, fmt.Sprintf("logs_%s.parquet", currentDate))
	var err error
	w.fw, err = local.NewLocalFileWriter(filename)
	if err != nil {
		return err
	}

	w.pw, err = writer.NewParquetWriter(w.fw, new(LogRecord), 4)
	if err != nil {
		w.fw.Close()
		return err
	}

	// 设置压缩方式
	w.pw.CompressionType = parquet.CompressionCodec_SNAPPY

	w.currentDate = currentDate
	w.currentFile = filename
	return nil
}

func (w *ParquetLogWriter) Write(log map[string]interface{}) error {
	if err := w.rotateIfNeeded(); err != nil {
		return err
	}

	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 转换日志格式
	record := &LogRecord{
		RequestID:   getString(log, "request_id"),
		Type:        getString(log, "type"),
		Time:        getString(log, "time"),
		Method:      getString(log, "method"),
		Path:        getString(log, "path"),
		Query:       getString(log, "query"),
		ClientIP:    getString(log, "client_ip"),
		UserAgent:   getString(log, "user_agent"),
		RequestBody: getString(log, "request_body"),
		Status:      getInt32(log, "status"),
		LatencyMS:   getInt64(log, "latency_ms"),
		Errors:      getString(log, "errors"),
	}

	return w.pw.Write(record)
}

func (w *ParquetLogWriter) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.pw != nil {
		if err := w.pw.WriteStop(); err != nil {
			return err
		}
	}
	if w.fw != nil {
		return w.fw.Close()
	}
	return nil
}

// 辅助函数：从map中安全获取字符串值
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// 辅助函数：从map中安全获取int32值
func getInt32(m map[string]interface{}, key string) int32 {
	if v, ok := m[key]; ok {
		switch t := v.(type) {
		case int:
			return int32(t)
		case int32:
			return t
		case float64:
			return int32(t)
		}
	}
	return 0
}

// 辅助函数：从map中安全获取int64值
func getInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		switch t := v.(type) {
		case int:
			return int64(t)
		case int64:
			return t
		case float64:
			return int64(t)
		}
	}
	return 0
}
