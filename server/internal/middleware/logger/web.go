package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// WebLogWriter Web回调日志写入器
type WebLogWriter struct {
	callbackURL string
	client      *http.Client
}

// NewWebLogWriter 创建Web回调日志写入器
func NewWebLogWriter(callbackURL string) *WebLogWriter {
	return &WebLogWriter{
		callbackURL: callbackURL,
		client:      &http.Client{Timeout: 5 * time.Second},
	}
}

func (w *WebLogWriter) Write(log map[string]interface{}) error {
	logJSON, _ := json.Marshal(log)
	_, err := w.client.Post(w.callbackURL, "application/json", bytes.NewBuffer(logJSON))
	return err
}

func (w *WebLogWriter) Close() error {
	return nil
}
