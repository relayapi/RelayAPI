package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"relayapi/server/internal/config"
)

// AESEncryptor 实现 AES 加密
type AESEncryptor struct {
	key    []byte
	ivSeed []byte
}

// NewAESEncryptor 创建 AES 加密器
func NewAESEncryptor(cfg *config.Config) (*AESEncryptor, error) {
	var key []byte
	if cfg.Crypto.AESKey != "" {
		// 使用配置的密钥
		key = []byte(cfg.Crypto.AESKey)
	} else {
		// 生成随机密钥
		key = make([]byte, cfg.Crypto.KeySize/8)
		if _, err := rand.Read(key); err != nil {
			return nil, fmt.Errorf("failed to generate AES key: %v", err)
		}
	}

	var ivSeed []byte
	if cfg.Crypto.AESIVSeed != "" {
		// 使用配置的 IV 种子
		ivSeed = []byte(cfg.Crypto.AESIVSeed)
	} else {
		// 生成随机 IV 种子
		ivSeed = make([]byte, aes.BlockSize)
		if _, err := rand.Read(ivSeed); err != nil {
			return nil, fmt.Errorf("failed to generate IV seed: %v", err)
		}
	}

	// 确保密钥长度正确
	if len(key) != cfg.Crypto.KeySize/8 {
		// 使用 SHA-256 调整密钥长度
		hash := sha256.Sum256(key)
		key = hash[:cfg.Crypto.KeySize/8]
	}

	return &AESEncryptor{
		key:    key,
		ivSeed: ivSeed,
	}, nil
}

// generateIV 生成 IV
func (e *AESEncryptor) generateIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}
	// 使用 IV 种子进行混合
	for i := 0; i < len(iv); i++ {
		iv[i] ^= e.ivSeed[i]
	}
	return iv, nil
}

// Encrypt 加密数据
func (e *AESEncryptor) Encrypt(data []byte) ([]byte, error) {
	// 创建 cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// 生成 IV
	iv, err := e.generateIV()
	if err != nil {
		return nil, err
	}

	// 加密数据
	ciphertext := make([]byte, len(data))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, data)

	// 组合 IV 和加密数据
	result := append(iv, ciphertext...)
	return []byte(base64.StdEncoding.EncodeToString(result)), nil
}

// Decrypt 解密数据
func (e *AESEncryptor) Decrypt(encryptedData []byte) ([]byte, error) {
	// 解码 base64
	data, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}

	// 提取 IV
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("encrypted data too short")
	}
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	// 创建 cipher
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
} 