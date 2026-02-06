package router

import (
	"department-eduvault-backend/controllers"
	internalController "department-eduvault-backend/internal/controller"
	internalService "department-eduvault-backend/internal/service"
	"department-eduvault-backend/middleware"
	"department-eduvault-backend/repositories"
	"department-eduvault-backend/services"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// New constructs the HTTP router and wires routes to controllers.
func New(healthService internalService.HealthService, dashboardService services.DashboardService, adminRepo repositories.AdminRepository, db *gorm.DB, logger *zap.Logger) *gin.Engine {
	engine := gin.New()
	engine.Use(
		middleware.CORSMiddleware(),
		gin.Recovery(),
		middleware.RequestLogger(logger),
		middleware.ErrorHandler(logger),
	)

	healthController := internalController.NewHealthController(healthService)
	engine.GET("/health", healthController.Health)

	// Auth Routes
	authController := controllers.NewAuthController()
	authGroup := engine.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
	}

	dashboardController := controllers.NewDashboardController(dashboardService)
	dashboard := engine.Group("/dashboard")
	dashboard.Use(middleware.MockAuthMiddleware("citchennai.net"))
	{
		dashboard.GET("/overview", dashboardController.GetOverview)
		dashboard.GET("/sections", dashboardController.GetSections)
		dashboard.GET("/recent-activity", dashboardController.GetRecentActivity)
		dashboard.GET("/export/certificates", dashboardController.ExportMyCertificates)
	}

	adminController := controllers.NewAdminController(adminRepo)
	admin := engine.Group("/admin")
	{
		admin.POST("/seed", adminController.Seed)
	}

	// Certificate workflows (faculty & HOD)
	certRepo := repositories.NewCertificateRepository(db)
	certService := services.NewCertificateService(certRepo)
	certController := controllers.NewCertificateController(certService)

	certificates := engine.Group("/certificates")
	certificates.Use(middleware.MockAuthMiddleware("citchennai.net"))
	{
		certificates.POST("/upload", certController.UploadCertificates)
		certificates.GET("/pending-review", certController.GetPendingReview)
		certificates.POST("/review", certController.SubmitReview)
	}

	// Student Routes
	studentRepo := repositories.NewStudentRepository(db)
	studentService := services.NewStudentService(studentRepo)
	studentController := controllers.NewStudentController(studentService)

	faculty := engine.Group("/faculty")
	faculty.Use(middleware.MockAuthMiddleware("citchennai.net"))
	{
		faculty.GET("/students", studentController.GetStudentsBySection)
	}

	// Verify endpoint (Faculty/HOD)
	engine.POST("/faculty/certificate/verify",
		middleware.MockAuthMiddleware("citchennai.net"),
		certController.TriggerMockVerification,
	)

	// HOD-facing endpoints
	hodRepo := repositories.NewHodRepository(db)
	hodService := services.NewHodService(hodRepo)
	hodController := controllers.NewHodController(hodService)

	hod := engine.Group("/hod")
	hod.Use(
		middleware.MockAuthMiddleware("citchennai.net"),
		middleware.RequireRoles("HOD"),
	)
	{
		hod.GET("/dashboard", hodController.HodDashboard)
		hod.GET("/faculty/students", hodController.GetStudentStats)
		hod.GET("/student/certificates", hodController.ListStudentCertificates)
		hod.GET("/export/certificates/section", hodController.ExportCertificatesBySection)
		hod.GET("/export/certificates/student", hodController.ExportCertificatesByStudent)
		hod.GET("/export/certificates/faculty", hodController.ExportCertificatesByFaculty)
	}

	return engine
}
