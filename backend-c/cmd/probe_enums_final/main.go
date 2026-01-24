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
		// Fallback if config fails
		cfg = &config.Config{DatabaseURL: "postgres://postgres:postgres@localhost:5432/eduvault?sslmode=disable"}
	}
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	var result string
	// Attempt to get one value to see case
	err = database.Raw("SELECT enum_range(null::faculty_status)::text").Scan(&result).Error
	if err != nil {
		fmt.Printf("Error probing faculty_status: %v\n", err)
	} else {
		fmt.Printf("FACULTY_STATUS_RAW: %s\n", result)
	}

	err = database.Raw("SELECT enum_range(null::ml_status)::text").Scan(&result).Error
	if err != nil {
		fmt.Printf("Error probing ml_status: %v\n", err)
	} else {
		fmt.Printf("ML_STATUS_RAW: %s\n", result)
	}
}
