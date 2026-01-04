package db

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/reno1r/weiss/apps/service/internal/config"
)

type Database struct {
	db     *gorm.DB
	config *config.Config
}

func NewDatabase(config *config.Config) (*Database, error) {
	dsn := buildDSN(config)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(config.AppDebug)),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if config.DatabaseMaxConnections > 0 {
		sqlDB.SetMaxOpenConns(config.DatabaseMaxConnections)
		sqlDB.SetMaxIdleConns(config.DatabaseMaxIdleConnections)
	}

	timeout := time.Duration(config.DatabaseConnectionTimeoutMs) * time.Millisecond
	if timeout > 0 {
		sqlDB.SetConnMaxLifetime(timeout)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{
		db:     db,
		config: config,
	}, nil
}

func (d *Database) DB() *gorm.DB {
	return d.db
}

func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func buildDSN(config *config.Config) string {
	sslMode := normalizeSSLMode(config.DatabaseSSL)

	if config.DatabasePassword != "" {
		return fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			config.DatabaseUser,
			config.DatabasePassword,
			config.DatabaseHost,
			config.DatabasePort,
			config.DatabaseName,
			sslMode,
		)
	}

	return fmt.Sprintf(
		"postgres://%s@%s:%d/%s?sslmode=%s",
		config.DatabaseUser,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseName,
		sslMode,
	)
}

func normalizeSSLMode(sslMode string) string {
	switch sslMode {
	case "", "false", "0", "no", "off":
		return "disable"
	case "true", "1", "yes", "on":
		return "require"
	default:
		return sslMode
	}
}

func getLogLevel(debug bool) logger.LogLevel {
	if debug {
		return logger.Info
	}
	return logger.Error
}
