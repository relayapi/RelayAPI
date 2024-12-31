package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
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
			Type             string `json:"type"`
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
	Server  ServerConfig
	Clients map[string]ClientConfig // key 是配置的 SHA256 hash
}

// GenerateConfigHash 根据 crypto 参数生成配置的 hash
func GenerateConfigHash(cfg *ClientConfig) string {
	data := cfg.Crypto.Method + cfg.Crypto.AESKey + cfg.Crypto.AESIVSeed
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// AddClientConfig 添加一个客户端配置
func (c *Config) AddClientConfig(cfg ClientConfig) string {
	hash := GenerateConfigHash(&cfg)
	c.Clients[hash] = cfg
	return hash
}

// GetClientConfig 根据 hash 获取客户端配置
func (c *Config) GetClientConfig(hash string) (ClientConfig, bool) {
	cfg, ok := c.Clients[hash]
	return cfg, ok
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
			Port:     8840,
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
	config := &Config{
		Clients: make(map[string]ClientConfig),
	}

	// 加载服务器配置
	serverData, err := os.ReadFile(serverConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read server config: %v", err)
	}
	if err := json.Unmarshal(serverData, &config.Server); err != nil {
		return nil, fmt.Errorf("failed to parse server config: %v", err)
	}

	// 如果客户端配置路径为空，加载默认配置
	if clientConfigPath == "" {
		defaultCfg := DefaultClientConfig()
		config.AddClientConfig(defaultCfg)
		return config, nil
	}

	// 检查是否是目录
	fileInfo, err := os.Stat(clientConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat client config path: %v", err)
	}

	if fileInfo.IsDir() {
		// 加载目录中的所有 .rai 文件
		files, err := os.ReadDir(clientConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read client config directory: %v", err)
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".rai") {
				filePath := filepath.Join(clientConfigPath, file.Name())
				if err := loadClientConfigFile(filePath, config); err != nil {
					log.Printf("Warning: failed to load client config %s: %v", filePath, err)
					continue
				}
				log.Println("load client config:", filePath)
			}
		}

		// 启动文件监控
		go watchConfigDirectory(clientConfigPath, config)
	} else {
		// 加载单个配置文件
		if err := loadClientConfigFile(clientConfigPath, config); err != nil {
			return nil, fmt.Errorf("failed to load client config: %v", err)
		}
	}

	return config, nil
}

// loadClientConfigFile 加载单个客户端配置文件
func loadClientConfigFile(filePath string, config *Config) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read client config file: %v", err)
	}

	var clientConfig ClientConfig
	if err := json.Unmarshal(data, &clientConfig); err != nil {
		log.Println("ERROR: clientConfig  unmarshal error ,filePath: ", filePath)
		return fmt.Errorf("failed to parse client config file: %v", err)
	}
	log.Println("clientConfig added :", clientConfig)
	config.AddClientConfig(clientConfig)
	return nil
}

// watchConfigDirectory 监控配置目录的变化
func watchConfigDirectory(dirPath string, config *Config) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Failed to create file watcher: %v", err)
		return
	}
	defer watcher.Close()

	if err := watcher.Add(dirPath); err != nil {
		log.Printf("Failed to watch directory: %v", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
				if strings.HasSuffix(event.Name, ".rai") {
					if err := loadClientConfigFile(event.Name, config); err != nil {
						log.Printf("Failed to load new client config %s: %v", event.Name, err)
					} else {
						log.Printf("Loaded new client config: %s", event.Name)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}
