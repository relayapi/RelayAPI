package crypto

import (
	"testing"
)

func TestKeyPairGeneration(t *testing.T) {
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	if keyPair.PrivateKey == nil {
		t.Error("Private key is nil")
	}
	if keyPair.PublicKey == nil {
		t.Error("Public key is nil")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// 生成密钥对
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// 测试数据
	testData := []byte("Hello, World!")

	// 加密
	encrypted, err := keyPair.Encrypt(testData)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// 解密
	decrypted, err := keyPair.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// 验证解密结果
	if string(decrypted) != string(testData) {
		t.Errorf("Decrypted data does not match original. Got %s, want %s",
			string(decrypted), string(testData))
	}
}

func TestPublicKeyExportImport(t *testing.T) {
	// 生成密钥对
	keyPair, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// 导出公钥
	publicKeyStr := keyPair.ExportPublicKey()
	if publicKeyStr == "" {
		t.Error("Exported public key is empty")
	}

	// 导入公钥
	importedKey, err := ImportPublicKey(publicKeyStr)
	if err != nil {
		t.Fatalf("Failed to import public key: %v", err)
	}

	// 验证导入的公钥
	if importedKey.X.Cmp(keyPair.PublicKey.X) != 0 || 
	   importedKey.Y.Cmp(keyPair.PublicKey.Y) != 0 {
		t.Error("Imported public key does not match original")
	}
} 