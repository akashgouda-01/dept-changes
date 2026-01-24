package main

import (
	"fmt"
	"log"

	"department-eduvault-backend/internal/config"
	"department-eduvault-backend/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{DatabaseURL: "postgres://postgres:postgres@localhost:5432/eduvault?sslmode=disable"}
	}
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	var types []string
	database.Raw("SELECT typname FROM pg_type WHERE typtype = 'e'").Scan(&types)
	fmt.Println("ENUM TYPES FOUND:")
	for _, t := range types {
		fmt.Printf("- %s\n", t)
	}
}
