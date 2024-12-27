package models

import (
	"testing"
	"time"
)

func TestTokenValidity(t *testing.T) {
	// 创建一个有效的令牌
	validToken := &Token{
		ID:           "test-token",
		EncryptedKey: "encrypted-key",
		MaxCalls:     100,
		UsedCalls:    50,
		ExpireTime:   time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 测试有效令牌
	if !validToken.IsValid() {
		t.Error("Token should be valid")
	}

	// 测试过期令牌
	expiredToken := &Token{
		ID:           "expired-token",
		EncryptedKey: "encrypted-key",
		MaxCalls:     100,
		UsedCalls:    50,
		ExpireTime:   time.Now().Add(-24 * time.Hour),
		CreatedAt:    time.Now().Add(-48 * time.Hour),
		UpdatedAt:    time.Now(),
	}

	if expiredToken.IsValid() {
		t.Error("Expired token should be invalid")
	}

	// 测试超出调用次数的令牌
	exhaustedToken := &Token{
		ID:           "exhausted-token",
		EncryptedKey: "encrypted-key",
		MaxCalls:     100,
		UsedCalls:    100,
		ExpireTime:   time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if exhaustedToken.IsValid() {
		t.Error("Exhausted token should be invalid")
	}
}

func TestTokenUsage(t *testing.T) {
	token := &Token{
		ID:           "test-token",
		EncryptedKey: "encrypted-key",
		MaxCalls:     100,
		UsedCalls:    0,
		ExpireTime:   time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 测试初始剩余调用次数
	if remaining := token.RemainingCalls(); remaining != 100 {
		t.Errorf("Expected 100 remaining calls, got %d", remaining)
	}

	// 测试增加使用次数
	token.IncrementUsage()
	if token.UsedCalls != 1 {
		t.Errorf("Expected 1 used call, got %d", token.UsedCalls)
	}

	// 测试更新时间
	if !token.UpdatedAt.After(token.CreatedAt) {
		t.Error("UpdatedAt should be after CreatedAt")
	}

	// 测试剩余调用次数
	if remaining := token.RemainingCalls(); remaining != 99 {
		t.Errorf("Expected 99 remaining calls, got %d", remaining)
	}
} 