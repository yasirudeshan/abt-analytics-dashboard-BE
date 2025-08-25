package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test with default environment (no env vars set)
	cfg := Load()

	if cfg.Port != ":" {
		t.Errorf("Expected Port to be ':', got '%s'", cfg.Port)
	}

	if cfg.DataFilePath != "" {
		t.Errorf("Expected DataFilePath to be empty, got '%s'", cfg.DataFilePath)
	}

	if cfg.Environment != "" {
		t.Errorf("Expected Environment to be empty, got '%s'", cfg.Environment)
	}
}

func TestLoadWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("DATA_FILE_PATH", "/path/to/data.csv")
	os.Setenv("ENVIRONMENT", "production")

	// Clean up after test
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATA_FILE_PATH")
		os.Unsetenv("ENVIRONMENT")
	}()

	cfg := Load()

	if cfg.Port != ":8080" {
		t.Errorf("Expected Port to be ':8080', got '%s'", cfg.Port)
	}

	if cfg.DataFilePath != "/path/to/data.csv" {
		t.Errorf("Expected DataFilePath to be '/path/to/data.csv', got '%s'", cfg.DataFilePath)
	}

	if cfg.Environment != "production" {
		t.Errorf("Expected Environment to be 'production', got '%s'", cfg.Environment)
	}
}

func TestLoadWithEmptyPort(t *testing.T) {
	// Set empty PORT
	os.Setenv("PORT", "")
	defer os.Unsetenv("PORT")

	cfg := Load()

	if cfg.Port != ":" {
		t.Errorf("Expected Port to be ':' when PORT is empty, got '%s'", cfg.Port)
	}
}

func TestConfigStruct(t *testing.T) {
	cfg := &Config{
		Port:         ":3000",
		DataFilePath: "/test/path.csv",
		Environment:  "test",
	}

	if cfg.Port != ":3000" {
		t.Errorf("Expected Port to be ':3000', got '%s'", cfg.Port)
	}

	if cfg.DataFilePath != "/test/path.csv" {
		t.Errorf("Expected DataFilePath to be '/test/path.csv', got '%s'", cfg.DataFilePath)
	}

	if cfg.Environment != "test" {
		t.Errorf("Expected Environment to be 'test', got '%s'", cfg.Environment)
	}
}

