package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
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

	// 检查是否存在密钥文件
	if _, fileErr := os.Stat(cfg.Crypto.PrivateKeyPath); os.IsNotExist(fileErr) {
		// 生成新的密钥对
		var genErr error
		keyPair, genErr = GenerateKeyPair()
		if genErr != nil {
			return nil, genErr
		}
		// TODO: 保存密钥到文件
	} else {
		// TODO: 从文件加载密钥
		var genErr error
		keyPair, genErr = GenerateKeyPair() // 临时使用生成的密钥
		if genErr != nil {
			return nil, genErr
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
	hash := sha256.Sum256(aesKey)
	r, s, err := ecdsa.Sign(rand.Reader, e.keyPair.PrivateKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt AES key: %v", err)
	}

	// 组合加密的密钥和数据
	encryptedKey := append(r.Bytes(), s.Bytes()...)
	encryptedKey = append(encryptedKey, aesKey...)
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

	// 分离签名、AES密钥和加密数据
	curve := e.keyPair.PublicKey.Curve
	keySize := (curve.Params().BitSize + 7) / 8
	signatureSize := keySize * 2
	if len(data) < signatureSize+32 {
		return nil, fmt.Errorf("encrypted data too short")
	}

	r := new(big.Int).SetBytes(data[:keySize])
	s := new(big.Int).SetBytes(data[keySize:signatureSize])
	aesKey := data[signatureSize : signatureSize+32]
	encryptedContent := data[signatureSize+32:]

	// 验证签名
	hash := sha256.Sum256(aesKey)
	if !ecdsa.Verify(e.keyPair.PublicKey, hash[:], r, s) {
		return nil, fmt.Errorf("invalid signature")
	}

	// 使用 AES 密钥解密数据
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