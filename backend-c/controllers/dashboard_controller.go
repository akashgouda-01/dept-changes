package controllers

import (
	"department-eduvault-backend/services"
	"department-eduvault-backend/utils"
	"github.com/gin-gonic/gin"
)

// DashboardController exposes read-only dashboard endpoints.
type DashboardController struct {
	service services.DashboardService
}

func NewDashboardController(service services.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

// GetOverview handles GET /dashboard/overview
func (dc *DashboardController) GetOverview(c *gin.Context) {
	overview, err := dc.service.GetOverview(c.Request.Context())
	if err != nil {
		_ = c.Error(utils.NewDatabaseError("failed to load dashboard overview", err))
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    overview,
	})
}

// GetSections handles GET /dashboard/sections
func (dc *DashboardController) GetSections(c *gin.Context) {
	sections, err := dc.service.GetSectionStats(c.Request.Context())
	if err != nil {
		_ = c.Error(utils.NewDatabaseError("failed to load section statistics", err))
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    sections,
	})
}
