package models

import (
	"time"
)

// Token 表示访问令牌
type Token struct {
	ID           string    `json:"id"`
	EncryptedKey string    `json:"encrypted_key"`
	MaxCalls     int       `json:"max_calls"`
	UsedCalls    int       `json:"used_calls"`
	ExpireTime   time.Time `json:"expire_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

// IncrementUsage 增加令牌使用次数
func (t *Token) IncrementUsage() {
	t.UsedCalls++
	t.UpdatedAt = time.Now()
}

// RemainingCalls 获取剩余调用次数
func (t *Token) RemainingCalls() int {
	return t.MaxCalls - t.UsedCalls
} 