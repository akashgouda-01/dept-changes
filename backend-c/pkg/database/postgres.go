package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

)

// DB is the global database handle used across the application.
var DB *sql.DB

var (
	// ErrMissingDSN is returned when DATABASE_URL is not set.
	ErrMissingDSN = errors.New("DATABASE_URL is not set")
)

// Init initializes the global DB connection using the DATABASE_URL environment variable.
func Init() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return ErrMissingDSN
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	// Reasonable defaults for a small departmental app
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	log.Println("✅ Connected to PostgreSQL")
	return nil
}

// Close closes the global DB connection. It is safe to call multiple times.
func Close() {
	if DB != nil {
		_ = DB.Close()
	}
}