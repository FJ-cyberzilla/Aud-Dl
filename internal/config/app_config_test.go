package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOrInitializeConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "settings.json")

	// Test case 1: Initialize default config
	cfg, err := LoadOrInitializeConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load/initialize config: %v", err)
	}
	if cfg.VaultDirectory == "" {
		t.Error("VaultDirectory should not be empty")
	}

	// Test case 2: Load existing config
	cfg2, err := LoadOrInitializeConfig(configPath)
	if err != nil {
		t.Fatalf("Failed to load existing config: %v", err)
	}
	if cfg2.VaultDirectory != cfg.VaultDirectory {
		t.Errorf("Expected VaultDirectory %s, got %s", cfg.VaultDirectory, cfg2.VaultDirectory)
	}
}
