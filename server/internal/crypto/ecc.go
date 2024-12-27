package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
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
	// 使用 ECDSA 进行加密
	r, s, err := ecdsa.Sign(rand.Reader, kp.PrivateKey, data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}

	// 将签名转换为字节数组
	signature := append(r.Bytes(), s.Bytes()...)
	return []byte(base64.StdEncoding.EncodeToString(signature)), nil
}

// Decrypt 使用私钥解密数据
func (kp *KeyPair) Decrypt(encryptedData []byte) ([]byte, error) {
	// 解码 base64
	signature, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %v", err)
	}

	// 分离 r 和 s
	rBytes := signature[:len(signature)/2]
	sBytes := signature[len(signature)/2:]

	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)

	// 验证签名
	hash := sha256.Sum256(encryptedData)
	valid := ecdsa.Verify(kp.PublicKey, hash[:], r, s)
	if !valid {
		return nil, fmt.Errorf("invalid signature")
	}

	return encryptedData, nil
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