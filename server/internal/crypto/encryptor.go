package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"relayapi/server/internal/config"
)

// Encryptor 定义加密接口
type Encryptor interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

// NewEncryptor 创建加密器
func NewEncryptor(cfg *config.Config) (Encryptor, error) {
	switch cfg.Crypto.Method {
	case "aes":
		// 解码 AES 密钥
		key, err := hex.DecodeString(cfg.Crypto.AESKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode AES key: %v", err)
		}
		if len(key) != 32 {
			// 使用 SHA-256 调整密钥长度
			hash := sha256.Sum256(key)
			key = hash[:]
		}

		// 解码 IV 种子
		ivSeed := []byte(cfg.Crypto.AESIVSeed)
		if len(ivSeed) != 16 {
			// 使用 SHA-256 调整 IV 种子长度
			hash := sha256.Sum256(ivSeed)
			ivSeed = hash[:16]
		}

		return NewAESEncryptor(key, ivSeed)
	case "ecc":
		return NewECCEncryptor(cfg)
	default:
		return nil, fmt.Errorf("unsupported encryption method: %s", cfg.Crypto.Method)
	}
} 