package config

import (
	"fmt"
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

	// 验证日志配置
	if cfg.Server.Log.Database.Enabled {
		if cfg.Server.Log.Database.ConnectionString == "" {
			return fmt.Errorf("database logging enabled but connection string is empty")
		}
	}
	if cfg.Server.Log.Web.Enabled {
		if cfg.Server.Log.Web.CallbackURL == "" {
			return fmt.Errorf("web logging enabled but callback URL is empty")
		}
	}
	if cfg.Server.Log.Parquet.Enabled {
		if cfg.Server.Log.Parquet.FilePath == "" {
			return fmt.Errorf("parquet logging enabled but file path is empty")
		}
	}

	// 验证速率限制配置
	if cfg.Server.RateLimit.RequestsPerSecond <= 0 {
		return fmt.Errorf("invalid requests per second")
	}
	if cfg.Server.RateLimit.Burst <= 0 {
		return fmt.Errorf("invalid burst size")
	}

	// 验证客户端配置
	if len(cfg.Clients) == 0 {
		return fmt.Errorf("no client configurations found")
	}

	for hash, clientCfg := range cfg.Clients {
		// 验证服务器配置
		if clientCfg.Server.Port <= 0 {
			return fmt.Errorf("invalid client server port for config %s", hash)
		}

		// 验证加密配置
		if clientCfg.Crypto.Method != "aes" {
			return fmt.Errorf("unsupported encryption method: %s for config %s", clientCfg.Crypto.Method, hash)
		}
		if len(clientCfg.Crypto.AESKey) != 64 {
			return fmt.Errorf("invalid AES key length for config %s", hash)
		}
		if len(clientCfg.Crypto.AESIVSeed) != 16 {
			return fmt.Errorf("invalid AES IV seed length for config %s", hash)
		}
	}

	return nil
}
