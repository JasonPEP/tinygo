package database

import (
	"fmt"
	"os"
	"path/filepath"

	"tinygo/internal/config"
	"tinygo/internal/logger"
	"tinygo/internal/shortener"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Init initializes the database connection
func Init(cfg config.DatabaseConfig) error {
	// Ensure data directory exists
	if dir := filepath.Dir(cfg.DSN); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create data directory: %w", err)
		}
	}

	// Configure GORM logger
	var gormLogLevel gormlogger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		gormLogLevel = gormlogger.Silent
	case "error":
		gormLogLevel = gormlogger.Error
	case "warn":
		gormLogLevel = gormlogger.Warn
	case "info":
		gormLogLevel = gormlogger.Info
	default:
		gormLogLevel = gormlogger.Warn
	}

	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormLogLevel),
	}

	// Connect to database
	var err error
	switch cfg.Driver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(cfg.DSN), gormConfig)
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	// Auto migrate
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	logger.Log.Info("database initialized", "driver", cfg.Driver, "dsn", cfg.DSN)
	return nil
}

// autoMigrate runs database migrations
func autoMigrate() error {
	return DB.AutoMigrate(&shortener.Link{})
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
