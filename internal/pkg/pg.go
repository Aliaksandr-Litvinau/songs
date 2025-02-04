package pg

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// DB is a wrapper around gorm.DB that provides additional functionality
// and better control over the database connection
type DB struct {
	*gorm.DB
}

// Dial creates and configures new database connection to postgres
// It handles connection establishment, configuration and verification
// Returns custom DB type that wraps gorm.DB
func Dial(dsn string) (*DB, error) {
	// Validate input
	if dsn == "" {
		return nil, errors.New("no postgres DSN provided")
	}

	// Сreate GORM config for postgres with a custom logger
	gormConfig := &gorm.Config{
		// Added full SQL logging for debugging (analog bundebug)
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			},
		),
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("gorm.Open failed: %w", err)
	}

	// Get base *sql.DB for settings pooling.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(1 * time.Minute)

	// Ping the database to verify the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &DB{db}, nil
}
