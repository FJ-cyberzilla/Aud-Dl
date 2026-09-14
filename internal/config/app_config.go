package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	VaultDirectory  string `json:"vault_directory"`
	MaxConcurrency  int    `json:"max_concurrency"`
	DefaultTheme    string `json:"default_theme"`
	RequestTimeout  int    `json:"request_timeout_seconds"`
	AutoUpdateCheck bool   `json:"auto_update_check"`
}

func DefaultConfig() *AppConfig {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return &AppConfig{
		VaultDirectory:  filepath.Join(homeDir, "AudioCommandCenter", "Vault"),
		MaxConcurrency:  4,
		DefaultTheme:    "cyber_dark",
		RequestTimeout:  30,
		AutoUpdateCheck: true,
	}
}

func ensureConfigDir(configPath string) error {
	dir := filepath.Dir(configPath)
	return os.MkdirAll(dir, 0755)
}

func saveConfig(configPath string, cfg *AppConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func loadConfig(configPath string) (*AppConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadOrInitializeConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		configPath = "config/app_settings.json"
	}
	if err := ensureConfigDir(configPath); err != nil {
		return nil, err
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := DefaultConfig()
		if err := saveConfig(configPath, cfg); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	return loadConfig(configPath)
}
