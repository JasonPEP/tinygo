package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Config holds runtime configuration for the server.
// Environment variables override file values when provided.
type Config struct {
	Addr       string `json:"addr" yaml:"addr" mapstructure:"addr"`
	BaseURL    string `json:"base_url" yaml:"base_url" mapstructure:"base_url"`
	DataFile   string `json:"data_file" yaml:"data_file" mapstructure:"data_file"`
	CodeLength int    `json:"code_length" yaml:"code_length" mapstructure:"code_length"`
	LogLevel   string `json:"log_level" yaml:"log_level" mapstructure:"log_level"`
	LogFormat  string `json:"log_format" yaml:"log_format" mapstructure:"log_format"`

	// Database configuration
	Database DatabaseConfig `json:"database" yaml:"database" mapstructure:"database"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string `json:"driver" yaml:"driver" mapstructure:"driver"`
	DSN      string `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
	LogLevel string `json:"log_level" yaml:"log_level" mapstructure:"log_level"`
}

// Default returns sane defaults for local development.
func Default() Config {
	return Config{
		Addr:       ":8080",
		BaseURL:    "http://localhost:8080",
		DataFile:   filepath.Join("data", "links.json"),
		CodeLength: 7,
		LogLevel:   "info",
		LogFormat:  "text",
		Database: DatabaseConfig{
			Driver:   "sqlite",
			DSN:      "data/tinygo.db",
			LogLevel: "warn",
		},
	}
}

// Load loads configuration with the following precedence (high -> low):
// 1) Environment variables
// 2) JSON file at configs/config.json (if exists)
// 3) Built-in defaults
// Note: We use JSON to avoid third-party YAML dependency.
func Load() (Config, error) {
	cfg := Default()

	// Load from file if present (configs/config.json)
	filePath := filepath.Join("configs", "config.json")
	if b, err := os.ReadFile(filePath); err == nil {
		var fc Config
		if err := json.Unmarshal(b, &fc); err != nil {
			return Config{}, err
		}
		merge(&cfg, fc)
	} else if !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}

	// Env overrides
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	if v := os.Getenv("BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("DATA_FILE"); v != "" {
		cfg.DataFile = v
	}
	if v := os.Getenv("CODE_LENGTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.CodeLength = n
		}
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("LOG_FORMAT"); v != "" {
		cfg.LogFormat = v
	}

	// Ensure data dir exists
	if dir := filepath.Dir(cfg.DataFile); dir != "." && dir != "" {
		_ = os.MkdirAll(dir, 0o755)
	}

	return cfg, nil
}

func merge(dst *Config, src Config) {
	if src.Addr != "" {
		dst.Addr = src.Addr
	}
	if src.BaseURL != "" {
		dst.BaseURL = src.BaseURL
	}
	if src.DataFile != "" {
		dst.DataFile = src.DataFile
	}
	if src.CodeLength > 0 {
		dst.CodeLength = src.CodeLength
	}
	if src.LogLevel != "" {
		dst.LogLevel = src.LogLevel
	}
	if src.LogFormat != "" {
		dst.LogFormat = src.LogFormat
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Addr == "" {
		return fmt.Errorf("addr cannot be empty")
	}
	if c.BaseURL == "" {
		return fmt.Errorf("base_url cannot be empty")
	}
	if c.CodeLength < 3 || c.CodeLength > 32 {
		return fmt.Errorf("code_length must be between 3 and 32")
	}

	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log_level: %s", c.LogLevel)
	}

	validLogFormats := map[string]bool{
		"text": true, "json": true,
	}
	if !validLogFormats[c.LogFormat] {
		return fmt.Errorf("invalid log_format: %s", c.LogFormat)
	}

	return nil
}
