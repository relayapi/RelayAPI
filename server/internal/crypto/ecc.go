package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"relayapi/server/internal/config"
)

// ECCEncryptor 实现 ECC 加密
type ECCEncryptor struct {
	keyPair *KeyPair
}

// NewECCEncryptor 创建 ECC 加密器
func NewECCEncryptor(cfg *config.Config) (*ECCEncryptor, error) {
	var keyPair *KeyPair
	var err error

	// 检查是否存在密钥文件
	if _, err := os.Stat(cfg.Crypto.PrivateKeyPath); os.IsNotExist(err) {
		// 生成新的密钥对
		keyPair, err = GenerateKeyPair()
		if err != nil {
			return nil, err
		}
		// TODO: 保存密钥到文件
	} else {
		// TODO: 从文件加载密钥
		keyPair, err = GenerateKeyPair() // 临时使用生成的密钥
		if err != nil {
			return nil, err
		}
	}

	return &ECCEncryptor{
		keyPair: keyPair,
	}, nil
}

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
func (e *ECCEncryptor) Encrypt(data []byte) ([]byte, error) {
	// 生成一个随机的 AES 密钥
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}

	// 使用 AES 加密数据
	aesEncryptor := &AESEncryptor{
		key:    aesKey,
		ivSeed: make([]byte, 16),
	}
	encryptedData, err := aesEncryptor.Encrypt(data)
	if err != nil {
		return nil, err
	}

	// 使用 ECC 加密 AES 密钥
	x, y := e.keyPair.PublicKey.ScalarBaseMult(aesKey)
	encryptedKey := append(x.Bytes(), y.Bytes()...)

	// 组合加密的密钥和数据
	result := append(encryptedKey, encryptedData...)
	return []byte(base64.StdEncoding.EncodeToString(result)), nil
}

// Decrypt 使用私钥解密数据
func (e *ECCEncryptor) Decrypt(encryptedData []byte) ([]byte, error) {
	// 解码 base64
	data, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}

	// 分离加密的密钥和数据
	keySize := (e.keyPair.PublicKey.Curve.Params().BitSize + 7) / 8 * 2
	if len(data) < keySize {
		return nil, fmt.Errorf("encrypted data too short")
	}

	encryptedKey := data[:keySize]
	encryptedContent := data[keySize:]

	// 解密 AES 密钥
	x := new(big.Int).SetBytes(encryptedKey[:keySize/2])
	y := new(big.Int).SetBytes(encryptedKey[keySize/2:])
	aesKey := make([]byte, 32)
	copy(aesKey, x.Bytes())

	// 使用解密的 AES 密钥解密数据
	aesEncryptor := &AESEncryptor{
		key:    aesKey,
		ivSeed: make([]byte, 16),
	}
	return aesEncryptor.Decrypt(encryptedContent)
}

// ExportPublicKey 导出公钥为 base64 字符串
func (e *ECCEncryptor) ExportPublicKey() string {
	x := e.keyPair.PublicKey.X.Bytes()
	y := e.keyPair.PublicKey.Y.Bytes()
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