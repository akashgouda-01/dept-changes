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

	var columns []struct {
		ColumnName string
		DataType   string
	}
	database.Raw("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'section_statistics';").Scan(&columns)
	fmt.Println("TABLE: section_statistics")
	for _, col := range columns {
		fmt.Printf("%s\n", col.ColumnName)
	}
}
