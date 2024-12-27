package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// Token 表示访问令牌
type Token struct {
	ID         string    `json:"id"`
	APIKey     string    `json:"api_key"`
	MaxCalls   int       `json:"max_calls"`
	UsedCalls  int       `json:"used_calls"`
	ExpireTime time.Time `json:"expire_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ExtInfo    string    `json:"ext_info,omitempty"`
}

// IsValid 检查令牌是否有效
func (t *Token) IsValid() bool {
	// 检查是否过期
	if time.Now().After(t.ExpireTime) {
		return false
	}

	// 检查使用次数
	if t.UsedCalls >= t.MaxCalls {
		return false
	}

	return true
}

// IncrementUsage 增加使用次数
func (t *Token) IncrementUsage() {
	t.UsedCalls++
	t.UpdatedAt = time.Now()
}

// Serialize 序列化令牌数据
func (t *Token) Serialize() ([]byte, error) {
	return json.Marshal(t)
}

// Deserialize 反序列化令牌数据
func (t *Token) Deserialize(data []byte) error {
	// 创建一个临时结构体来解析时间字符串
	type TempToken struct {
		ID         string `json:"id"`
		APIKey     string `json:"api_key"`
		MaxCalls   int    `json:"max_calls"`
		UsedCalls  int    `json:"used_calls"`
		ExpireTime string `json:"expire_time"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
		ExtInfo    string `json:"ext_info,omitempty"`
	}

	var temp TempToken
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal token data: %v (data: %s)", err, string(data))
	}

	// 解析时间字符串
	expireTime, err := time.Parse(time.RFC3339, temp.ExpireTime)
	if err != nil {
		return fmt.Errorf("failed to parse expire_time: %v", err)
	}

	createdAt, err := time.Parse(time.RFC3339, temp.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to parse created_at: %v", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, temp.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to parse updated_at: %v", err)
	}

	// 设置字段值
	t.ID = temp.ID
	t.APIKey = temp.APIKey
	t.MaxCalls = temp.MaxCalls
	t.UsedCalls = temp.UsedCalls
	t.ExpireTime = expireTime
	t.CreatedAt = createdAt
	t.UpdatedAt = updatedAt
	t.ExtInfo = temp.ExtInfo

	return nil
}
