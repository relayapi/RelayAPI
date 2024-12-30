package logger

import (
	"log"

	"github.com/hokaccha/go-prettyjson"
)

// ConsoleLogWriter 控制台日志写入器
type ConsoleLogWriter struct {
	formatter *prettyjson.Formatter
}

// NewConsoleLogWriter 创建一个新的控制台日志写入器
func NewConsoleLogWriter() *ConsoleLogWriter {
	formatter := prettyjson.NewFormatter()
	formatter.DisabledColor = false // 启用颜色输出
	formatter.Indent = 2            // 设置缩进

	return &ConsoleLogWriter{
		formatter: formatter,
	}
}

func (w *ConsoleLogWriter) Write(log1 map[string]interface{}) error {
	output, err := w.formatter.Marshal(log1)
	if err != nil {
		return err
	}

	log.Println(string(output))
	return nil
}

func (w *ConsoleLogWriter) Close() error {
	return nil
}
