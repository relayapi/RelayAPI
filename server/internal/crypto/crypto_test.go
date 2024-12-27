package crypto

import (
	"relayapi/server/internal/config"
	"testing"
)

func TestNewEncryptor(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		shouldError bool
	}{
		{
			name:        "AES Encryption",
			method:      "aes",
			shouldError: false,
		},
		{
			name:        "ECC Encryption",
			method:      "ecc",
			shouldError: false,
		},
		{
			name:        "Invalid Method",
			method:      "invalid",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.Crypto.Method = tt.method
			cfg.Crypto.KeySize = 256

			encryptor, err := NewEncryptor(cfg)
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
		})
	}
} 