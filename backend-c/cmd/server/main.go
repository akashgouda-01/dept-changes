package main

import (
	"context"
	"log"
	"time"

	"department-eduvault-backend/internal/config"
	"department-eduvault-backend/internal/db"
	internalRepository "department-eduvault-backend/internal/repository"
	"department-eduvault-backend/internal/router"
	"department-eduvault-backend/internal/server"
	internalService "department-eduvault-backend/internal/service"
	"department-eduvault-backend/models"
	"department-eduvault-backend/repositories"
	"department-eduvault-backend/services"
	"department-eduvault-backend/utils"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := utils.NewLogger()
	defer logger.Sync() // flush buffered logs

	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect to database",
			zap.Error(err),
		)
	}

	// Quick connectivity check on startup.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.HealthCheck(ctx, database); err != nil {
		logger.Fatal("database ping failed",
			zap.Error(err),
		)
	}

	logger.Info("database connection established and ping successful")

	// Ensure ENUM types exist (idempotent-ish via exception handling in migration or ignored if error)
	// Note: basic Exec might fail if type exists, so we ignore error or use DO block.
	// Since we can't easily do complex blocks here without clutter, we'll try a simple create.
	// A better way is to rely on Gorm or Pre-migration.
	// Let's assume the user might restart often, so we wrap in DO block for safety if possible?
	// Go's Gorm Exec:
	database.Exec("DO $$ BEGIN CREATE TYPE ml_status_enum AS ENUM ('PENDING', 'VERIFIED', 'DUPLICATE'); EXCEPTION WHEN duplicate_object THEN null; END $$;")
	database.Exec("DO $$ BEGIN CREATE TYPE faculty_status_enum AS ENUM ('PENDING', 'LEGIT', 'NOT_LEGIT'); EXCEPTION WHEN duplicate_object THEN null; END $$;")

	// Auto-migrate tables for certificate workflow
	if err := database.AutoMigrate(
		&models.Student{},
		&models.Certificate{},
		&models.StudentStatistics{},
		&models.SectionStatistics{},
	); err != nil {
		logger.Fatal("failed to auto-migrate database", zap.Error(err))
	}

	// Backfill UpdatedAt if null, ensuring legacy data has timestamps
	database.Exec("UPDATE faculty_certificates SET updated_at = uploaded_at WHERE updated_at IS NULL")

	healthRepo := internalRepository.NewHealthRepository(database)
	healthService := internalService.NewHealthService(healthRepo)

	dashboardRepo := repositories.NewDashboardRepository(database)
	dashboardService := services.NewDashboardService(dashboardRepo)

	adminRepo := repositories.NewAdminRepository(database)

	engine := router.New(healthService, dashboardService, adminRepo, database, logger)

	srv := server.New(engine, cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
