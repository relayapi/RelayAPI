package logger

import (
	"fmt"
)

// LogWriter 日志写入器接口
type LogWriter interface {
	Write(log map[string]interface{}) error
	Close() error
}

// AsyncLogWriter 异步日志写入器包装器
type AsyncLogWriter struct {
	writer  LogWriter
	logChan chan map[string]interface{}
	done    chan struct{}
}

// NewAsyncLogWriter 创建异步日志写入器
func NewAsyncLogWriter(writer LogWriter, bufferSize int) *AsyncLogWriter {
	w := &AsyncLogWriter{
		writer:  writer,
		logChan: make(chan map[string]interface{}, bufferSize),
		done:    make(chan struct{}),
	}
	go w.processLogs()
	return w
}

func (w *AsyncLogWriter) processLogs() {
	for log := range w.logChan {
		if err := w.writer.Write(log); err != nil {
			fmt.Printf("Failed to write log: %v\n", err)
		}
	}
	w.done <- struct{}{}
}

func (w *AsyncLogWriter) Write(log map[string]interface{}) error {
	select {
	case w.logChan <- log:
		return nil
	default:
		return fmt.Errorf("log buffer full")
	}
}

func (w *AsyncLogWriter) Close() error {
	close(w.logChan)
	<-w.done
	return w.writer.Close()
}

// CloseLogWriters 关闭所有日志写入器
func CloseLogWriters(writers []LogWriter) {
	for _, writer := range writers {
		if err := writer.Close(); err != nil {
			fmt.Printf("Failed to close log writer: %v\n", err)
		}
	}
}
