// Package config manages configuration via Viper.
// Priority: flag > env var > config file.
// Config file location: ~/.config/linear-cli/config.yaml
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	// Config keys.
	KeyAPIKey      = "api_key"
	KeyDefaultTeam = "default_team"

	// Env var names.
	EnvAPIKey = "LINEAR_API_KEY"
)

var configDir string

// Init initializes Viper configuration.
// Call this once during application startup.
func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir = filepath.Join(home, ".config", "linear-cli")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Environment variable bindings.
	viper.SetEnvPrefix("")
	_ = viper.BindEnv(KeyAPIKey, EnvAPIKey)

	// Read config file (ignore error if file doesn't exist).
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only return error if it's not a "file not found" error.
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	return nil
}

// GetAPIKey returns the API key from config or env var.
func GetAPIKey() string {
	return viper.GetString(KeyAPIKey)
}

// GetDefaultTeam returns the default team key from config.
func GetDefaultTeam() string {
	return viper.GetString(KeyDefaultTeam)
}

// Set saves a key-value pair to the config file.
func Set(key, value string) error {
	viper.Set(key, value)
	return save()
}

// Get returns the value for a given config key.
func Get(key string) string {
	return viper.GetString(key)
}

// GetAll returns all config settings as a map.
func GetAll() map[string]interface{} {
	return viper.AllSettings()
}

// ConfigDir returns the path to the config directory.
func ConfigDir() string {
	return configDir
}

// save writes the current config to disk, creating the directory if needed.
func save() error {
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")

	// If config file doesn't exist, create it.
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if _, err := os.Create(configPath); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
