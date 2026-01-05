package controllers

import (
	"net/http"

	"department-eduvault-backend/repositories"
	"department-eduvault-backend/utils"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminRepo repositories.AdminRepository
}

func NewAdminController(adminRepo repositories.AdminRepository) *AdminController {
	return &AdminController{adminRepo: adminRepo}
}

// Seed inserts sample data into the certificates table if it is empty.
func (ac *AdminController) Seed(c *gin.Context) {
	if err := ac.adminRepo.SeedCertificatesIfEmpty(c.Request.Context()); err != nil {
		_ = c.Error(utils.NewDatabaseError("failed to seed certificates", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "seed completed (no-op if data already present)",
	})
}
