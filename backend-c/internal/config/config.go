package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds environment-driven settings for the application.
type Config struct {
	AppEnv             string
	Port               string
	DatabaseURL        string
	AllowedEmailDomain string
}

// Load reads configuration from environment variables and optional .env file.
func Load() (*Config, error) {
	_ = godotenv.Load() // Best-effort; ignore missing file.

	cfg := &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", "8080"),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		AllowedEmailDomain: getEnv("ALLOWED_EMAIL_DOMAIN", "citchennai.net"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
