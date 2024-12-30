package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	tempConfig := `{
		"server": {
			"port": 9090
		},
		"database": {
			"host": "testhost",
			"port": 5432,
			"user": "testuser",
			"password": "testpass",
			"dbname": "testdb"
		},
		"crypto": {
			"method": "aes",
			"key_size": 192,
			"private_key_path": "test/private.pem",
			"public_key_path": "test/public.pem",
			"aes_key": "testkey123",
			"aes_iv_seed": "testiv456"
		}
	}`

	tmpfile, err := os.CreateTemp("", "config.*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(tempConfig)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// 测试加载配置
	config, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Errorf("LoadConfig failed: %v", err)
	}

	// 验证配置值
	if config.Server.Port != 9090 {
		t.Errorf("Expected port 9090, got %d", config.Server.Port)
	}
	if config.Database.Host != "testhost" {
		t.Errorf("Expected host testhost, got %s", config.Database.Host)
	}

	// 验证加密配置
	if config.Crypto.Method != "aes" {
		t.Errorf("Expected crypto method aes, got %s", config.Crypto.Method)
	}
	if config.Crypto.KeySize != 192 {
		t.Errorf("Expected key size 192, got %d", config.Crypto.KeySize)
	}
	if config.Crypto.AESKey != "testkey123" {
		t.Errorf("Expected AES key testkey123, got %s", config.Crypto.AESKey)
	}
	if config.Crypto.AESIVSeed != "testiv456" {
		t.Errorf("Expected AES IV seed testiv456, got %s", config.Crypto.AESIVSeed)
	}
}

func TestLoadConfigDefault(t *testing.T) {
	// 测试加载不存在的配置文件时使用默认值
	config, err := LoadConfig("nonexistent.json")
	if err != nil {
		t.Errorf("LoadConfig failed: %v", err)
	}

	// 验证默认值
	if config.Server.Port != 8840 {
		t.Errorf("Expected default port 8840, got %d", config.Server.Port)
	}
	if config.Database.Host != "localhost" {
		t.Errorf("Expected default host localhost, got %s", config.Database.Host)
	}

	// 验证默认加密配置
	if config.Crypto.Method != "aes" {
		t.Errorf("Expected default crypto method aes, got %s", config.Crypto.Method)
	}
	if config.Crypto.KeySize != 256 {
		t.Errorf("Expected default key size 256, got %d", config.Crypto.KeySize)
	}
	if config.Crypto.AESKey != "" {
		t.Errorf("Expected empty default AES key, got %s", config.Crypto.AESKey)
	}
	if config.Crypto.AESIVSeed != "" {
		t.Errorf("Expected empty default AES IV seed, got %s", config.Crypto.AESIVSeed)
	}
}

func TestSaveConfig(t *testing.T) {
	config := &Config{}
	config.Server.Port = 9090
	config.Database.Host = "testhost"
	config.Crypto.Method = "aes"
	config.Crypto.KeySize = 192
	config.Crypto.AESKey = "testkey123"
	config.Crypto.AESIVSeed = "testiv456"

	tmpfile, err := os.CreateTemp("", "config.*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// 保存配置
	if err := SaveConfig(config, tmpfile.Name()); err != nil {
		t.Errorf("SaveConfig failed: %v", err)
	}

	// 重新加载配置并验证
	loadedConfig, err := LoadConfig(tmpfile.Name())
	if err != nil {
		t.Errorf("LoadConfig failed: %v", err)
	}

	// 验证所有字段
	if loadedConfig.Server.Port != config.Server.Port {
		t.Errorf("Port mismatch: expected %d, got %d", config.Server.Port, loadedConfig.Server.Port)
	}
	if loadedConfig.Database.Host != config.Database.Host {
		t.Errorf("Host mismatch: expected %s, got %s", config.Database.Host, loadedConfig.Database.Host)
	}
	if loadedConfig.Crypto.Method != config.Crypto.Method {
		t.Errorf("Crypto method mismatch: expected %s, got %s", config.Crypto.Method, loadedConfig.Crypto.Method)
	}
	if loadedConfig.Crypto.KeySize != config.Crypto.KeySize {
		t.Errorf("Key size mismatch: expected %d, got %d", config.Crypto.KeySize, loadedConfig.Crypto.KeySize)
	}
	if loadedConfig.Crypto.AESKey != config.Crypto.AESKey {
		t.Errorf("AES key mismatch: expected %s, got %s", config.Crypto.AESKey, loadedConfig.Crypto.AESKey)
	}
	if loadedConfig.Crypto.AESIVSeed != config.Crypto.AESIVSeed {
		t.Errorf("AES IV seed mismatch: expected %s, got %s", config.Crypto.AESIVSeed, loadedConfig.Crypto.AESIVSeed)
	}
}
