package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const configFileName = "config.json"

// Config holds application configuration.
type Config struct {
	VaultPath string `json:"vault_path"`
}

// Dir returns the config directory path (~/.config/lazymd).
func Dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "lazymd")
}

// Load reads config from disk. Returns zero-value Config if file is missing.
// Returns error on parse failure.
func Load() (*Config, error) {
	path := filepath.Join(Dir(), configFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes config to disk atomically (write to tmp, then rename).
// Creates config directory (0755) if needed.
func Save(cfg *Config) error {
	dir := Dir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	tmp := filepath.Join(dir, configFileName+".tmp")
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, configFileName))
}

// NeedsSetup returns true if VaultPath is empty (string check only).
func NeedsSetup(cfg *Config) bool {
	return cfg.VaultPath == ""
}

// ExpandPath expands ~ to home dir and resolves relative paths to absolute.
func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[2:])
	} else if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = home
	}
	return filepath.Abs(path)
}
