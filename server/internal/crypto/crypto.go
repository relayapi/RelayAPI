package crypto

import (
	"fmt"
	"relayapi/server/internal/config"
)

// Encryptor 定义加密接口
type Encryptor interface {
	// Encrypt 加密数据
	Encrypt(data []byte) ([]byte, error)
	// Decrypt 解密数据
	Decrypt(encryptedData []byte) ([]byte, error)
}

// NewEncryptor 创建加密器实例
func NewEncryptor(cfg *config.Config) (Encryptor, error) {
	switch cfg.Crypto.Method {
	case "aes":
		return NewAESEncryptor(cfg)
	case "ecc":
		return NewECCEncryptor(cfg)
	default:
		return nil, fmt.Errorf("unsupported encryption method: %s", cfg.Crypto.Method)
	}
} 