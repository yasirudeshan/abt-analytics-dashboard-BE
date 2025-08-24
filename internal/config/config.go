package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port         string
	DataFilePath string
	Environment  string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:         ":" + os.Getenv("PORT"),
		DataFilePath: os.Getenv("DATA_FILE_PATH"),
		Environment:  os.Getenv("ENVIRONMENT"),
	}
}
