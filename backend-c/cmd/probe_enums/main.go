package main

import (
	"fmt"
	"log"

	"department-eduvault-backend/internal/config"
	"department-eduvault-backend/internal/db"
)

func main() {
	cfg, _ := config.Load()
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	var values []string
	database.Raw("SELECT unnest(enum_range(NULL::faculty_status))::text").Scan(&values)
	fmt.Println("ENUM: faculty_status")
	for _, v := range values {
		fmt.Printf("'%s'\n", v)
	}

	var mlValues []string
	database.Raw("SELECT unnest(enum_range(NULL::ml_status))::text").Scan(&mlValues)
	fmt.Println("ENUM: ml_status")
	for _, v := range mlValues {
		fmt.Printf("'%s'\n", v)
	}
}
