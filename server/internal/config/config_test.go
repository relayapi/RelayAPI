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
			"private_key_path": "test/private.pem",
			"public_key_path": "test/public.pem"
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
}

func TestLoadConfigDefault(t *testing.T) {
	// 测试加载不存在的配置文件时使用默认值
	config, err := LoadConfig("nonexistent.json")
	if err != nil {
		t.Errorf("LoadConfig failed: %v", err)
	}

	// 验证默认值
	if config.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", config.Server.Port)
	}
	if config.Database.Host != "localhost" {
		t.Errorf("Expected default host localhost, got %s", config.Database.Host)
	}
} 