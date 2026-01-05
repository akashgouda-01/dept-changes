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

	healthRepo := internalRepository.NewHealthRepository(database)
	healthService := internalService.NewHealthService(healthRepo)

	dashboardRepo := repositories.NewDashboardRepository(database)
	dashboardService := services.NewDashboardService(dashboardRepo)

	adminRepo := repositories.NewAdminRepository(database)

	engine := router.New(healthService, dashboardService, adminRepo, logger)

	srv := server.New(engine, cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
