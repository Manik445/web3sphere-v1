package database

import (
	"fmt"
	"time"

	"github.com/web3sphere/backend/configs"
	"github.com/web3sphere/backend/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// New creates a new PostgreSQL database connection using GORM.
func New(cfg *configs.DatabaseConfig, log *logger.Logger) (*gorm.DB, error) {
	logLevel := gormlogger.Silent
	switch cfg.LogLevel {
	case "info":
		logLevel = gormlogger.Info
	case "warn":
		logLevel = gormlogger.Warn
	case "error":
		logLevel = gormlogger.Error
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// Health check
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connected successfully")
	return db, nil
}

// Close cleanly shuts down the database connection.
func Close(db *gorm.DB, log *logger.Logger) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("Failed to get underlying sql.DB for close: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Errorf("Failed to close database connection: %v", err)
	} else {
		log.Info("Database connection closed")
	}
}
