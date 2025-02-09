package pg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresDB is a wrapper around gorm.DB that provides additional functionality
// and better control over the database connection.
// Instead of using embedding, we explicitly define methods used in repositories,
// which gives us several advantages:
// 1. Full control over the public API of our type
// 2. Ability to add additional functionality (logging, metrics, tracing)
// 3. Flexibility to modify method behavior when needed
// 4. Cleaner and more understandable code without name conflicts
// 5. Looser coupling with the specific gorm implementation
//
// In the future, if we need to add a new method or change existing behavior,
// we can do it in one place without affecting client code.
// For example, we can easily add logging or tracing:
//
//	func (db *PostgresDB) WithContext(ctx context.Context) *gorm.DB {
//	    span, ctx := tracer.StartSpan(ctx, "db.query")
//	    defer span.End()
//	    return db.gorm.WithContext(ctx)
//	}
type PostgresDB struct {
	gorm *gorm.DB
}

func NewPostgresDB(db *gorm.DB) *PostgresDB {
	return &PostgresDB{gorm: db}
}

// Close closes the database connection
func (db *PostgresDB) Close() error {
	sqlDB, err := db.gorm.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// WithContext returns gorm DB with context
func (db *PostgresDB) WithContext(ctx context.Context) *gorm.DB {
	return db.gorm.WithContext(ctx)
}

// Create inserts value into database
func (db *PostgresDB) Create(value interface{}) *gorm.DB {
	return db.gorm.Create(value)
}

// First finds first record that match given conditions
func (db *PostgresDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return db.gorm.First(dest, conds...)
}

// Model specify the model you would like to run db operations
func (db *PostgresDB) Model(value interface{}) *gorm.DB {
	return db.gorm.Model(value)
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (db *PostgresDB) Save(value interface{}) *gorm.DB {
	return db.gorm.Save(value)
}

// Delete delete value match given conditions
func (db *PostgresDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return db.gorm.Delete(value, conds...)
}

// Dial creates and configures new database connection to postgres
// It handles connection establishment, configuration and verification
// Returns custom DB type that wraps gorm.DB
func Dial(dsn string) (*PostgresDB, error) {
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

	return &PostgresDB{db}, nil
}
