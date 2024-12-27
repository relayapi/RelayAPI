package crypto

import (
	"relayapi/server/internal/config"
	"testing"
)

func TestAESEncryptorCreation(t *testing.T) {
	tests := []struct {
		name        string
		keySize     int
		key         string
		ivSeed      string
		shouldError bool
	}{
		{
			name:        "Default Settings",
			keySize:     256,
			shouldError: false,
		},
		{
			name:        "Custom Key",
			keySize:     256,
			key:         "mysecretkey12345",
			shouldError: false,
		},
		{
			name:        "Custom IV Seed",
			keySize:     256,
			ivSeed:      "myivseed12345678",
			shouldError: false,
		},
		{
			name:        "Custom Key and IV",
			keySize:     256,
			key:         "mysecretkey12345",
			ivSeed:      "myivseed12345678",
			shouldError: false,
		},
		{
			name:        "Invalid Key Size",
			keySize:     123,
			shouldError: false, // Should not error as we adjust key size
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.Crypto.KeySize = tt.keySize
			cfg.Crypto.AESKey = tt.key
			cfg.Crypto.AESIVSeed = tt.ivSeed

			encryptor, err := NewAESEncryptor(cfg)
			if tt.shouldError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if encryptor == nil {
				t.Error("Expected encryptor, got nil")
			}

			// 验证密钥长度
			expectedKeySize := tt.keySize / 8
			if len(encryptor.key) != expectedKeySize {
				t.Errorf("Expected key size %d, got %d", expectedKeySize, len(encryptor.key))
			}

			// 验证 IV 种子长度
			if len(encryptor.ivSeed) != 16 {
				t.Errorf("Expected IV seed size 16, got %d", len(encryptor.ivSeed))
			}
		})
	}
}

func TestAESEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		keySize int
		key     string
		ivSeed  string
	}{
		{
			name:    "Short Text",
			data:    []byte("Hello, World!"),
			keySize: 256,
		},
		{
			name:    "Empty Text",
			data:    []byte(""),
			keySize: 256,
		},
		{
			name:    "Long Text",
			data:    []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
			keySize: 256,
		},
		{
			name:    "Binary Data",
			data:    []byte{0x00, 0x01, 0x02, 0x03, 0x04},
			keySize: 256,
		},
		{
			name:    "Custom Key",
			data:    []byte("Test with custom key"),
			keySize: 256,
			key:     "mysecretkey12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.Crypto.KeySize = tt.keySize
			cfg.Crypto.AESKey = tt.key
			cfg.Crypto.AESIVSeed = tt.ivSeed

			encryptor, err := NewAESEncryptor(cfg)
			if err != nil {
				t.Fatalf("Failed to create encryptor: %v", err)
			}

			// 加密
			encrypted, err := encryptor.Encrypt(tt.data)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			// 解密
			decrypted, err := encryptor.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			// 验证结果
			if string(decrypted) != string(tt.data) {
				t.Errorf("Decrypted data does not match original.\nGot: %s\nWant: %s",
					string(decrypted), string(tt.data))
			}
		})
	}
}

func TestAESGenerateIV(t *testing.T) {
	cfg := &config.Config{}
	cfg.Crypto.KeySize = 256

	encryptor, err := NewAESEncryptor(cfg)
	if err != nil {
		t.Fatalf("Failed to create encryptor: %v", err)
	}

	// 生成多个 IV 并验证它们是否不同
	iv1, err := encryptor.generateIV()
	if err != nil {
		t.Fatalf("Failed to generate IV: %v", err)
	}

	iv2, err := encryptor.generateIV()
	if err != nil {
		t.Fatalf("Failed to generate IV: %v", err)
	}

	// 验证 IV 长度
	if len(iv1) != 16 {
		t.Errorf("Expected IV length 16, got %d", len(iv1))
	}

	// 验证两个 IV 是否不同
	if string(iv1) == string(iv2) {
		t.Error("Generated IVs should be different")
	}
} 