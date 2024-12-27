package config

import (
	"encoding/json"
	"os"
)

// Config 表示服务器配置
type Config struct {
	Server struct {
		Port int `json:"port"`
	} `json:"server"`
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	} `json:"database"`
	Crypto struct {
		PrivateKeyPath string `json:"private_key_path"`
		PublicKeyPath  string `json:"public_key_path"`
	} `json:"crypto"`
}

var defaultConfig = Config{
	Server: struct {
		Port int `json:"port"`
	}{
		Port: 8080,
	},
	Database: struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
	}{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		DBName:   "relayapi",
	},
	Crypto: struct {
		PrivateKeyPath string `json:"private_key_path"`
		PublicKeyPath  string `json:"public_key_path"`
	}{
		PrivateKeyPath: "keys/private.pem",
		PublicKeyPath:  "keys/public.pem",
	},
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	config := defaultConfig

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果配置文件不存在，使用默认配置
			return &config, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
} 