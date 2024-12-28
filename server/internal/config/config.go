package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// ValidateConfig 验证配置是否有效
func ValidateConfig(cfg *Config) error {
	// 验证服务器配置
	if cfg.Server.Server.Port <= 0 {
		return fmt.Errorf("invalid server port")
	}
	if cfg.Server.Server.ReadTimeout <= 0 {
		return fmt.Errorf("invalid read timeout")
	}
	if cfg.Server.Server.WriteTimeout <= 0 {
		return fmt.Errorf("invalid write timeout")
	}

	// 验证数据库配置
	if cfg.Server.Database.Port <= 0 {
		return fmt.Errorf("invalid database port")
	}
	if cfg.Server.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("invalid max open connections")
	}
	if cfg.Server.Database.MaxIdleConns <= 0 {
		return fmt.Errorf("invalid max idle connections")
	}
	if cfg.Server.Database.ConnMaxLifetime <= 0 {
		return fmt.Errorf("invalid connection max lifetime")
	}

	// 验证速率限制配置
	if cfg.Server.RateLimit.RequestsPerSecond <= 0 {
		return fmt.Errorf("invalid requests per second")
	}
	if cfg.Server.RateLimit.Burst <= 0 {
		return fmt.Errorf("invalid burst size")
	}

	// 验证客户端配置
	if cfg.Client.Server.Port <= 0 {
		return fmt.Errorf("invalid client server port")
	}
	if cfg.Client.Crypto.Method != "aes" {
		return fmt.Errorf("unsupported encryption method: %s", cfg.Client.Crypto.Method)
	}
	if len(cfg.Client.Crypto.AESKey) != 64 {
		return fmt.Errorf("invalid AES key length")
	}
	if len(cfg.Client.Crypto.AESIVSeed) != 16 {
		return fmt.Errorf("invalid AES IV seed length")
	}

	return nil
}

// SaveConfig 保存配置到文件
func SaveConfig(cfg *Config, serverConfigPath string, clientConfigPath string) error {
	// 保存服务器配置
	serverData, err := json.MarshalIndent(cfg.Server, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal server config: %v", err)
	}
	if err := os.WriteFile(serverConfigPath, serverData, 0644); err != nil {
		return fmt.Errorf("failed to write server config: %v", err)
	}

	// 如果指定了客户端配置路径，保存客户端配置
	if clientConfigPath != "" {
		clientData, err := json.MarshalIndent(cfg.Client, "", "    ")
		if err != nil {
			return fmt.Errorf("failed to marshal client config: %v", err)
		}
		if err := os.WriteFile(clientConfigPath, clientData, 0644); err != nil {
			return fmt.Errorf("failed to write client config: %v", err)
		}
	}

	return nil
}
