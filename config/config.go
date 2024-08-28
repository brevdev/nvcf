package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	APIKey string `json:"api_key"`
}

var cfg Config

func Init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	// TODO: consider more robust config loading here
	configPath := filepath.Join(configDir, "nvcf", "config.json")
	data, err := os.ReadFile(configPath)
	if err == nil {
		json.Unmarshal(data, &cfg)
	}

	if v := os.Getenv("NGC_API_KEY"); v != "" {
		cfg.APIKey = v
	}
	if v := os.Getenv("NGC_CLI_API_KEY"); v != "" {
		cfg.APIKey = v
	}
}

func GetAPIKey() string {
	return cfg.APIKey
}

func SetAPIKey(apiKey string) error {
	cfg.APIKey = apiKey
	return saveConfig()
}

func ClearAPIKey() error {
	cfg.APIKey = ""
	return saveConfig()
}

func IsAuthenticated() bool {
	return cfg.APIKey != ""
}

// todo: consider reaching for more secret-specific storage.
func saveConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "nvcf", "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
