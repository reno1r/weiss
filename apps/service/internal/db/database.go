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

func NewDatabase(cfg *config.Config) (*Database, error) {
	dsn := buildDSN(cfg)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(getLogLevel(cfg.AppDebug)),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if cfg.DatabaseMaxConnections > 0 {
		sqlDB.SetMaxOpenConns(cfg.DatabaseMaxConnections)
		sqlDB.SetMaxIdleConns(cfg.DatabaseMaxIdleConnections)
	}

	timeout := time.Duration(cfg.DatabaseConnectionTimeoutMs) * time.Millisecond
	if timeout > 0 {
		sqlDB.SetConnMaxLifetime(timeout)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{
		db:     db,
		config: cfg,
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

func buildDSN(cfg *config.Config) string {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s",
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabaseName,
	)

	if cfg.DatabasePassword != "" {
		dsn += fmt.Sprintf(" password=%s", cfg.DatabasePassword)
	}

	if cfg.DatabaseSSL != "" {
		dsn += fmt.Sprintf(" sslmode=%s", cfg.DatabaseSSL)
	} else {
		dsn += " sslmode=disable"
	}

	return dsn
}

func getLogLevel(debug bool) logger.LogLevel {
	if debug {
		return logger.Info
	}
	return logger.Error
}
