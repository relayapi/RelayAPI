package models

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Token 表示访问令牌
type Token struct {
	ID         string    `json:"id"`
	APIKey     string    `json:"api_key"`
	MaxCalls   int       `json:"max_calls"`
	ExpireTime time.Time `json:"expire_time"`
	CreatedAt  time.Time `json:"created_at"`
	Provider   string    `json:"provider"`  // API 提供商：openai, dashscope 等
	ExtInfo    string    `json:"ext_info,omitempty"`
}

// 使用计数器
var (
	usageCounters = make(map[string]int)  // token ID -> 使用次数
	usageMutex    sync.RWMutex
)

// IsValid 检查令牌是否有效
func (t *Token) IsValid() bool {
	// 检查是否过期
	if time.Now().After(t.ExpireTime) {
		return false
	}

	// 检查使用次数
	usageMutex.RLock()
	usedCalls := usageCounters[t.ID]
	usageMutex.RUnlock()

	if usedCalls >= t.MaxCalls {
		return false
	}

	return true
}

// IncrementUsage 增加使用次数
func (t *Token) IncrementUsage() {
	usageMutex.Lock()
	usageCounters[t.ID]++
	usageMutex.Unlock()
}

// GetUsage 获取使用次数
func (t *Token) GetUsage() int {
	usageMutex.RLock()
	count := usageCounters[t.ID]
	usageMutex.RUnlock()
	return count
}

// GetRemainingCalls 获取剩余调用次数
func (t *Token) GetRemainingCalls() int {
	usageMutex.RLock()
	count := usageCounters[t.ID]
	usageMutex.RUnlock()
	return t.MaxCalls - count
}

// ResetUsage 重置使用次数
func (t *Token) ResetUsage() {
	usageMutex.Lock()
	delete(usageCounters, t.ID)
	usageMutex.Unlock()
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
		ExpireTime string `json:"expire_time"`
		CreatedAt  string `json:"created_at"`
		Provider   string `json:"provider"`
		ExtInfo    string `json:"ext_info,omitempty"`
	}

	var temp TempToken
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal token data: %v (data: %s)", err, string(data))
	}

	// 验证必填字段
	if temp.ID == "" || temp.APIKey == "" || temp.Provider == "" {
		return fmt.Errorf("missing required fields")
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

	// 设置字段值
	t.ID = temp.ID
	t.APIKey = temp.APIKey
	t.MaxCalls = temp.MaxCalls
	t.ExpireTime = expireTime
	t.CreatedAt = createdAt
	t.Provider = temp.Provider
	t.ExtInfo = temp.ExtInfo

	return nil
}
