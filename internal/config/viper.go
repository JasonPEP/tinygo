package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// LoadWithViper loads configuration using viper with the following precedence:
// 1) Environment variables (with prefix TINYGO_)
// 2) YAML/JSON/TOML config file (configs/config.yaml)
// 3) Built-in defaults
func LoadWithViper() (Config, error) {
	// Set default values
	viper.SetDefault("addr", ":8080")
	viper.SetDefault("base_url", "http://localhost:8080")
	viper.SetDefault("data_file", "data/links.json")
	viper.SetDefault("code_length", 7)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("log_format", "text")
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.dsn", "data/tinygo.db")
	viper.SetDefault("database.log_level", "warn")

	// Set config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") // 支持 yaml, json, toml
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Environment variables
	viper.SetEnvPrefix("TINYGO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Also bind common environment variables without prefix
	viper.BindEnv("base_url", "BASE_URL")
	viper.BindEnv("addr", "ADDR")
	viper.BindEnv("log_level", "LOG_LEVEL")
	viper.BindEnv("log_format", "LOG_FORMAT")
	viper.BindEnv("database.driver", "DATABASE_DRIVER")
	viper.BindEnv("database.dsn", "DATABASE_DSN")
	viper.BindEnv("database.log_level", "DATABASE_LOG_LEVEL")

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return Config{}, fmt.Errorf("read config file: %w", err)
		}
		// Config file not found is OK, use defaults + env
	}

	// Unmarshal into struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	// Handle PORT environment variable (Railway specific)
	if port := os.Getenv("PORT"); port != "" && cfg.Addr == ":8080" {
		cfg.Addr = ":" + port
	}

	// Handle DATABASE_URL environment variable (Railway specific)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		cfg.Database.Driver = "postgres"
		cfg.Database.DSN = databaseURL
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// WatchConfig enables hot reloading of config file
func WatchConfig() {
	viper.WatchConfig()
}
