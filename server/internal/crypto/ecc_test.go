package crypto

import (
	"relayapi/server/internal/config"
	"testing"
)

func TestECCEncryptorCreation(t *testing.T) {
	cfg := &config.Config{}
	cfg.Crypto.PrivateKeyPath = "test/private.pem"
	cfg.Crypto.PublicKeyPath = "test/public.pem"

	encryptor, err := NewECCEncryptor(cfg)
	if err != nil {
		t.Fatalf("Failed to create ECC encryptor: %v", err)
	}

	if encryptor == nil {
		t.Error("Expected encryptor, got nil")
	}

	if encryptor.keyPair == nil {
		t.Error("Expected key pair, got nil")
	}
}

func TestECCEncryptDecrypt(t *testing.T) {
	cfg := &config.Config{}
	cfg.Crypto.PrivateKeyPath = "test/private.pem"
	cfg.Crypto.PublicKeyPath = "test/public.pem"

	encryptor, err := NewECCEncryptor(cfg)
	if err != nil {
		t.Fatalf("Failed to create ECC encryptor: %v", err)
	}

	// 测试数据
	testData := []byte("Hello, World!")

	// 加密
	encrypted, err := encryptor.Encrypt(testData)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	// 解密
	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %v", err)
	}

	// 验证解密结果
	if string(decrypted) != string(testData) {
		t.Errorf("Decrypted data does not match original. Got %s, want %s",
			string(decrypted), string(testData))
	}
}

func TestECCPublicKeyExport(t *testing.T) {
	cfg := &config.Config{}
	cfg.Crypto.PrivateKeyPath = "test/private.pem"
	cfg.Crypto.PublicKeyPath = "test/public.pem"

	encryptor, err := NewECCEncryptor(cfg)
	if err != nil {
		t.Fatalf("Failed to create ECC encryptor: %v", err)
	}

	// 导出公钥
	publicKeyStr := encryptor.ExportPublicKey()
	if publicKeyStr == "" {
		t.Error("Exported public key is empty")
	}

	// 导入公钥
	importedKey, err := ImportPublicKey(publicKeyStr)
	if err != nil {
		t.Fatalf("Failed to import public key: %v", err)
	}

	// 验证导入的公钥
	if importedKey.X.Cmp(encryptor.keyPair.PublicKey.X) != 0 ||
		importedKey.Y.Cmp(encryptor.keyPair.PublicKey.Y) != 0 {
		t.Error("Imported public key does not match original")
	}
} 