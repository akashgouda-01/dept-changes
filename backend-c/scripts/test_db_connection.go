package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"department-eduvault-backend/internal/config"
	"department-eduvault-backend/internal/db"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Config error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📋 Configuration loaded:\n")
	fmt.Printf("   - Port: %s\n", cfg.Port)
	fmt.Printf("   - Database URL: %s\n", maskPassword(cfg.DatabaseURL))
	fmt.Printf("\n")

	fmt.Printf("🔌 Attempting database connection...\n")
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		fmt.Printf("❌ Connection failed: %v\n", err)
		fmt.Printf("\n")
		fmt.Printf("💡 Troubleshooting steps:\n")
		fmt.Printf("   1. Verify your Supabase project is active (not paused)\n")
		fmt.Printf("   2. Check DATABASE_URL in .env file\n")
		fmt.Printf("   3. Verify network connectivity to Supabase\n")
		fmt.Printf("   4. Try using the connection pooler URL (port 6543) instead of direct (5432)\n")
		os.Exit(1)
	}

	fmt.Printf("✅ Database connection established!\n")

	// Test ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("🏓 Testing database ping...\n")
	if err := db.HealthCheck(ctx, database); err != nil {
		fmt.Printf("❌ Ping failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Database ping successful!\n")
	fmt.Printf("\n🎉 Database connection is working correctly!\n")
}

func maskPassword(dsn string) string {
	// Simple password masking for display
	if idx := strings.Index(dsn, "@"); idx > 0 {
		if colonIdx := strings.LastIndex(dsn[:idx], ":"); colonIdx > 0 {
			return dsn[:colonIdx+1] + "***" + dsn[idx:]
		}
	}
	return dsn
}
