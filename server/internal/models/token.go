package models

import (
	"time"
)

// Token 表示访问令牌
type Token struct {
	ID           string    `json:"id"`
	APIKey       string    `json:"api_key"`
	MaxCalls     int       `json:"max_calls"`
	ExpireTime   time.Time `json:"expire_time"`
	ExtInfo      string    `json:"ext_info"`
}

// IsValid 检查令牌是否有效
func (t *Token) IsValid() bool {
	// 检查是否过期
	if time.Now().After(t.ExpireTime) {
		return false
	}

	// 检查调用次数
	if t.UsedCalls >= t.MaxCalls {
		return false
	}

	return true
}
