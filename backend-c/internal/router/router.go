package router

import (
	"department-eduvault-backend/controllers"
	internalController "department-eduvault-backend/internal/controller"
	internalService "department-eduvault-backend/internal/service"
	"department-eduvault-backend/middleware"
	"department-eduvault-backend/repositories"
	"department-eduvault-backend/services"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// New constructs the HTTP router and wires routes to controllers.
func New(healthService internalService.HealthService, dashboardService services.DashboardService, adminRepo repositories.AdminRepository, logger *zap.Logger) *gin.Engine {
	engine := gin.New()
	engine.Use(
		middleware.CORSMiddleware(),
		gin.Recovery(),
		middleware.RequestLogger(logger),
		middleware.ErrorHandler(logger),
	)

	healthController := internalController.NewHealthController(healthService)
	engine.GET("/health", healthController.Health)

	dashboardController := controllers.NewDashboardController(dashboardService)
	dashboard := engine.Group("/dashboard")
	{
		dashboard.GET("/overview", dashboardController.GetOverview)
		dashboard.GET("/sections", dashboardController.GetSections)
	}

	adminController := controllers.NewAdminController(adminRepo)
	admin := engine.Group("/admin")
	{
		admin.POST("/seed", adminController.Seed)
	}

	return engine
}
