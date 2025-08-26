package config

import (
	"os"
	"testing"
)

// TestLoadWithMockEnvironment tests configuration loading with mock environment variables
func TestLoadWithMockEnvironment(t *testing.T) {
	// Set mock environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("DATA_FILE_PATH", "/mock/path/data.csv")
	os.Setenv("ENVIRONMENT", "testing")

	// Clean up after test
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("DATA_FILE_PATH")
		os.Unsetenv("ENVIRONMENT")
	}()

	cfg := Load()

	// Verify mock values are loaded correctly
	if cfg.Port != ":9090" {
		t.Errorf("Expected Port to be ':9090', got '%s'", cfg.Port)
	}
	if cfg.DataFilePath != "/mock/path/data.csv" {
		t.Errorf("Expected DataFilePath to be '/mock/path/data.csv', got '%s'", cfg.DataFilePath)
	}
	if cfg.Environment != "testing" {
		t.Errorf("Expected Environment to be 'testing', got '%s'", cfg.Environment)
	}
}

// TestLoadWithEmptyEnvironment tests configuration loading with no environment variables
func TestLoadWithEmptyEnvironment(t *testing.T) {
	// Ensure no environment variables are set
	os.Unsetenv("PORT")
	os.Unsetenv("DATA_FILE_PATH")
	os.Unsetenv("ENVIRONMENT")

	cfg := Load()

	// Verify default values
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

// TestConfigStructWithMockData tests configuration struct with hardcoded values
func TestConfigStructWithMockData(t *testing.T) {
	cfg := &Config{
		Port:         ":3000",
		DataFilePath: "/mock/test/path.csv",
		Environment:  "mock_test",
	}

	// Verify mock values
	if cfg.Port != ":3000" {
		t.Errorf("Expected Port to be ':3000', got '%s'", cfg.Port)
	}
	if cfg.DataFilePath != "/mock/test/path.csv" {
		t.Errorf("Expected DataFilePath to be '/mock/test/path.csv', got '%s'", cfg.DataFilePath)
	}
	if cfg.Environment != "mock_test" {
		t.Errorf("Expected Environment to be 'mock_test', got '%s'", cfg.Environment)
	}
}

// TestPortFormattingWithMockData tests port formatting with various mock inputs
func TestPortFormattingWithMockData(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"8080", ":8080"},
		{"3000", ":3000"},
		{"", ":"},
		{"9090", ":9090"},
	}

	for _, tc := range testCases {
		os.Setenv("PORT", tc.input)
		cfg := Load()

		if cfg.Port != tc.expected {
			t.Errorf("For input '%s', expected Port '%s', got '%s'", tc.input, tc.expected, cfg.Port)
		}

		os.Unsetenv("PORT")
	}
}

// TestEnvironmentOverrideWithMockData tests environment variable override behavior
func TestEnvironmentOverrideWithMockData(t *testing.T) {
	// Set initial values
	os.Setenv("PORT", "8080")
	os.Setenv("DATA_FILE_PATH", "/initial/path.csv")
	os.Setenv("ENVIRONMENT", "development")

	// Load first configuration
	cfg1 := Load()

	// Change environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("DATA_FILE_PATH", "/new/path.csv")
	os.Setenv("ENVIRONMENT", "production")

	// Load second configuration
	cfg2 := Load()

	// Verify configurations are different
	if cfg1.Port == cfg2.Port {
		t.Error("Expected different Port values after environment change")
	}
	if cfg1.DataFilePath == cfg2.DataFilePath {
		t.Error("Expected different DataFilePath values after environment change")
	}
	if cfg1.Environment == cfg2.Environment {
		t.Error("Expected different Environment values after environment change")
	}

	// Verify new values
	if cfg2.Port != ":9090" {
		t.Errorf("Expected new Port ':9090', got '%s'", cfg2.Port)
	}
	if cfg2.DataFilePath != "/new/path.csv" {
		t.Errorf("Expected new DataFilePath '/new/path.csv', got '%s'", cfg2.DataFilePath)
	}
	if cfg2.Environment != "production" {
		t.Errorf("Expected new Environment 'production', got '%s'", cfg2.Environment)
	}

	// Clean up
	os.Unsetenv("PORT")
	os.Unsetenv("DATA_FILE_PATH")
	os.Unsetenv("ENVIRONMENT")
}
