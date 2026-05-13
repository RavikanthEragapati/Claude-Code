package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Name         string `json:"name"`
	Team         string `json:"team"`
	Framework    string `json:"framework"`
	SprintNumber int    `json:"sprintNumber"`
	SprintLength int    `json:"sprintLength"`
}

// ConfigDir returns the base directory for pm config and data.
func ConfigDir() string {
	return filepath.Join("~", ".pm")
}

// ConfigFile returns the config file path.
func ConfigFile() string {
	return filepath.Join(ConfigDir(), "config.json")
}

// SaveConfig writes the config to disk.
func SaveConfig(cfg *Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFile(), data, 0644)
}

// LoadConfig reads config from disk; returns nil if file does not exist.
func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigFile())
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}