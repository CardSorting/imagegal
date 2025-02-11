package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig
	ModelsLab ModelsLabConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// ModelsLabConfig holds ModelsLab API configuration
type ModelsLabConfig struct {
	APIKey     string
	BaseURL    string
	MaxRetries int
}

// New creates a new Config instance with values from environment variables
func New() (*Config, error) {
	// Load .env file if it exists
	if err := loadEnvFile(); err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	port, err := strconv.Atoi(getEnvOrDefault("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid server port: %w", err)
	}

	readTimeout, err := time.ParseDuration(getEnvOrDefault("SERVER_READ_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid read timeout: %w", err)
	}

	writeTimeout, err := time.ParseDuration(getEnvOrDefault("SERVER_WRITE_TIMEOUT", "10s"))
	if err != nil {
		return nil, fmt.Errorf("invalid write timeout: %w", err)
	}

	maxRetries, err := strconv.Atoi(getEnvOrDefault("MODELSLAB_MAX_RETRIES", "3"))
	if err != nil {
		return nil, fmt.Errorf("invalid max retries: %w", err)
	}

	apiKey := os.Getenv("MODELSLAB_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("MODELSLAB_API_KEY environment variable is required")
	}

	return &Config{
		Server: ServerConfig{
			Port:         port,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
		ModelsLab: ModelsLabConfig{
			APIKey:     apiKey,
			BaseURL:    getEnvOrDefault("MODELSLAB_BASE_URL", "https://modelslab.com/api/v6"),
			MaxRetries: maxRetries,
		},
	}, nil
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile() error {
	// Look for .env file in current directory and parent directories
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	for {
		envFile := filepath.Join(dir, ".env")
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Load(envFile); err != nil {
				return fmt.Errorf("failed to load %s: %w", envFile, err)
			}
			return nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Don't return error if .env file is not found
	return nil
}

// getEnvOrDefault returns the value of an environment variable or a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
