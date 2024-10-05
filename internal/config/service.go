package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (c *Config) Load() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Failed to get user home directory: %w", err)
	}

	storageDir := filepath.Join(homeDir, ".tasks")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return fmt.Errorf("Failed to create storage directory: %w", err)
	}

	if err := c.ensureFile(storageDir); err != nil {
		return fmt.Errorf("Failed to create config file: %w", err)
	}

	data, err := os.ReadFile(c.filepath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, c); err != nil {
			return fmt.Errorf("failed to decode config: %w", err)
		}
	}

	if c.ServiceMode == "" {
		c.ServiceMode = ServiceModeSQL
	}

	updatedData, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.filepath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) UpdateServiceMode(mode ServiceMode) error {
	c.ServiceMode = mode

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(c.filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) ensureFile(homeDir string) error {
	p := filepath.Join(homeDir, "config.json")
	c.filepath = p
	if _, err := os.Stat(p); os.IsNotExist(err) {
		data, err := json.Marshal(c)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %w", err)
		}

		if err := os.WriteFile(p, data, 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	return nil
}
