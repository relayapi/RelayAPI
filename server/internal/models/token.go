package models

import (
	"encoding/json"
	"time"
)

// Token 表示访问令牌
type Token struct {
	ID           string    `json:"id"`            // 令牌ID
	APIKey       string    `json:"api_key"`       // 实际的 API Key
	MaxCalls     int       `json:"max_calls"`     // 最大调用次数
	UsedCalls    int       `json:"used_calls"`    // 已使用的调用次数
	ExpireTime   time.Time `json:"expire_time"`   // 过期时间
	CreatedAt    time.Time `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`    // 更新时间
	ExtInfo      string    `json:"ext_info"`      // 扩展信息
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

// Serialize 序列化令牌为 JSON 字节数组
func (t *Token) Serialize() ([]byte, error) {
	return json.Marshal(t)
}

// Deserialize 从 JSON 字节数组反序列化令牌
func (t *Token) Deserialize(data []byte) error {
	return json.Unmarshal(data, t)
}
