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

	var result string
	database.Raw("SELECT enum_range(NULL::faculty_status_enum)::text").Scan(&result)
	fmt.Printf("FACULTY_STATUS_ENUM: %s\n", result)

	var res2 string
	database.Raw("SELECT enum_range(NULL::ml_status_enum)::text").Scan(&res2)
	fmt.Printf("ML_STATUS_ENUM: %s\n", res2)
}
