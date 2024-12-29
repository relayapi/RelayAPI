package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// ClientConfig 客户端配置结构
type ClientConfig struct {
	Version string `json:"version"`
	Server  struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		BasePath string `json:"base_path"`
	} `json:"server"`
	Crypto struct {
		Method    string `json:"method"`
		AESKey    string `json:"aes_key"`
		AESIVSeed string `json:"aes_iv_seed"`
	} `json:"crypto"`
}

// ServerConfig 服务器配置结构
type ServerConfig struct {
	Server struct {
		Host           string `json:"host"`
		Port           int    `json:"port"`
		ReadTimeout    int    `json:"read_timeout"`
		WriteTimeout   int    `json:"write_timeout"`
		MaxHeaderBytes int    `json:"max_header_bytes"`
		Debug          bool   `json:"debug"`
	} `json:"server"`
	Log struct {
		Console  bool `json:"console"`
		Database struct {
			Enabled          bool   `json:"enabled"`
			ConnectionString string `json:"connection_string"`
		} `json:"database"`
		Web struct {
			Enabled     bool   `json:"enabled"`
			CallbackURL string `json:"callback_url"`
		} `json:"web"`
		Parquet struct {
			Enabled  bool   `json:"enabled"`
			FilePath string `json:"file_path"`
		} `json:"parquet"`
	} `json:"log"`
	RateLimit struct {
		RequestsPerSecond int `json:"requests_per_second"`
		Burst             int `json:"burst"`
		IPLimit           struct {
			RequestsPerSecond int `json:"requests_per_second"`
			Burst             int `json:"burst"`
		} `json:"ip_limit"`
	} `json:"rate_limit"`
}

// Config 完整配置结构
type Config struct {
	Client ClientConfig
	Server ServerConfig
}

// DefaultClientConfig 创建默认的客户端配置
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		Version: "1.0.0",
		Server: struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			BasePath string `json:"base_path"`
		}{
			Host:     "http://localhost",
			Port:     8080,
			BasePath: "/relayapi/",
		},
		Crypto: struct {
			Method    string `json:"method"`
			AESKey    string `json:"aes_key"`
			AESIVSeed string `json:"aes_iv_seed"`
		}{
			Method:    "aes",
			AESKey:    "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			AESIVSeed: "fedcba9876543210",
		},
	}
}

// LoadConfig 加载配置
func LoadConfig(serverConfigPath string, clientConfigPath string) (*Config, error) {
	config := &Config{}

	// 加载服务器配置
	serverData, err := os.ReadFile(serverConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read server config: %v", err)
	}
	if err := json.Unmarshal(serverData, &config.Server); err != nil {
		return nil, fmt.Errorf("failed to parse server config: %v", err)
	}

	// 尝试加载客户端配置
	if clientConfigPath != "" {
		clientData, err := os.ReadFile(clientConfigPath)
		if err == nil {
			if err := json.Unmarshal(clientData, &config.Client); err != nil {
				return nil, fmt.Errorf("failed to parse client config: %v", err)
			}
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read client config: %v", err)
		}
	}

	// 如果没有找到客户端配置，尝试加载默认配置
	if config.Client == (ClientConfig{}) {
		defaultPath := "default.rai"
		clientData, err := os.ReadFile(defaultPath)
		if err == nil {
			if err := json.Unmarshal(clientData, &config.Client); err != nil {
				return nil, fmt.Errorf("failed to parse default client config: %v", err)
			}
		} else if os.IsNotExist(err) {
			// 创建默认配置
			config.Client = DefaultClientConfig()
			defaultData, err := json.MarshalIndent(config.Client, "", "    ")
			if err != nil {
				return nil, fmt.Errorf("failed to marshal default client config: %v", err)
			}
			if err := os.WriteFile(defaultPath, defaultData, 0644); err != nil {
				return nil, fmt.Errorf("failed to write default client config: %v", err)
			}
			fmt.Printf("Created default client config at %s\n", defaultPath)
		} else {
			return nil, fmt.Errorf("failed to read default client config: %v", err)
		}
	}

	return config, nil
}
