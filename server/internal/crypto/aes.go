package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type AESEncryptor struct {
	key    []byte
	ivSeed []byte
}

func NewAESEncryptor(key []byte, ivSeed []byte) (*AESEncryptor, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("AES key must be 32 bytes (256 bits)")
	}
	if len(ivSeed) != aes.BlockSize {
		return nil, fmt.Errorf("IV seed must be %d bytes", aes.BlockSize)
	}
	return &AESEncryptor{
		key:    key,
		ivSeed: ivSeed,
	}, nil
}

func (e *AESEncryptor) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// 使用 CBC 模式
	mode := cipher.NewCBCEncrypter(block, e.ivSeed)

	// 填充数据
	paddedData := pkcs7Padding(data, aes.BlockSize)

	// 加密数据
	ciphertext := make([]byte, len(paddedData))
	mode.CryptBlocks(ciphertext, paddedData)

	return ciphertext, nil
}

func (e *AESEncryptor) Decrypt(data []byte) ([]byte, error) {
	// 检查数据长度
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// 提取 IV（前 16 字节）
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	// 创建解密器
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// 使用 CBC 模式
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充
	unpaddedData, err := pkcs7Unpadding(plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to remove padding: %v", err)
	}

	return unpaddedData, nil
}

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// PKCS7 去除填充
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("empty data")
	}
	
	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("invalid padding size")
	}
	
	// 验证所有填充字节
	for i := length - padding; i < length; i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding")
		}
	}
	
	return data[:length-padding], nil
} 