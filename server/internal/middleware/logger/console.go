package logger

import (
	"encoding/json"
	"fmt"
)

// ConsoleLogWriter 控制台日志写入器
type ConsoleLogWriter struct{}

func (w *ConsoleLogWriter) Write(log map[string]interface{}) error {
	logJSON, _ := json.Marshal(log)
	fmt.Println(string(logJSON))
	return nil
}

func (w *ConsoleLogWriter) Close() error {
	return nil
}
