package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
)

// KeyPair 表示 ECC 密钥对
type KeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

// GenerateKeyPair 生成新的 ECC 密钥对
func GenerateKeyPair() (*KeyPair, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %v", err)
	}

	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  &privateKey.PublicKey,
	}, nil
}

// Encrypt 使用公钥加密数据
func (kp *KeyPair) Encrypt(data []byte) ([]byte, error) {
	// 生成一个随机的 AES 密钥
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}

	// 使用 AES 加密数据
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// 生成随机 IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %v", err)
	}

	// 加密数据
	ciphertext := make([]byte, len(data))
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext, data)

	// 组合所有数据：IV + AES密钥 + 加密数据
	result := append(iv, aesKey...)
	result = append(result, ciphertext...)

	return []byte(base64.StdEncoding.EncodeToString(result)), nil
}

// Decrypt 使用私钥解密数据
func (kp *KeyPair) Decrypt(encryptedData []byte) ([]byte, error) {
	// 解码 base64
	data, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}

	// 提取 IV
	iv := data[:aes.BlockSize]
	// 提取 AES 密钥
	aesKey := data[aes.BlockSize : aes.BlockSize+32]
	// 提取加密数据
	ciphertext := data[aes.BlockSize+32:]

	// 解密数据
	plaintext := make([]byte, len(ciphertext))
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}

// ExportPublicKey 导出公钥为 base64 字符串
func (kp *KeyPair) ExportPublicKey() string {
	x := kp.PublicKey.X.Bytes()
	y := kp.PublicKey.Y.Bytes()
	return base64.StdEncoding.EncodeToString(append(x, y...))
}

// ImportPublicKey 从 base64 字符串导入公钥
func ImportPublicKey(publicKeyStr string) (*ecdsa.PublicKey, error) {
	data, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %v", err)
	}

	x := new(big.Int).SetBytes(data[:len(data)/2])
	y := new(big.Int).SetBytes(data[len(data)/2:])

	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}, nil
} 