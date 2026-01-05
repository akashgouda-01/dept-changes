package db

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect opens a GORM connection to PostgreSQL using the provided DSN.
func Connect(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database DSN is empty")
	}

	normalized, err := ensureSSLMode(dsn)
	if err != nil {
		return nil, fmt.Errorf("normalize DSN: %w", err)
	}

	// Configure GORM with connection settings
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(postgres.Open(normalized), config)
	if err != nil {
		// Provide more context about the connection failure
		return nil, fmt.Errorf("opening database connection: %w (DSN format: postgresql://user:***@host:port/db)", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("getting sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(15)

	// Set connection timeout
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}

// HealthCheck performs a simple ping against the database.
func HealthCheck(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}
	return db.WithContext(ctx).Exec("SELECT 1").Error
}

// ensureSSLMode makes sure the DSN has sslmode=require (needed for Supabase)
// without hardcoding any credentials.
func ensureSSLMode(dsn string) (string, error) {
	// If sslmode is already present, leave as-is.
	if strings.Contains(dsn, "sslmode=") {
		return dsn, nil
	}

	u, err := url.Parse(dsn)
	if err != nil || u.Scheme == "" {
		// Fallback for DSN strings that are not URL-form (let postgres driver handle details).
		// Append sslmode=require conservatively.
		sep := "?"
		if strings.Contains(dsn, "?") {
			sep = "&"
		}
		return dsn + sep + "sslmode=require", nil
	}

	// Parse and add SSL mode
	q := u.Query()
	q.Set("sslmode", "require")
	// Add connect_timeout for better error handling
	if q.Get("connect_timeout") == "" {
		q.Set("connect_timeout", "10")
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
