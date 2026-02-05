package controllers

import (
	"department-eduvault-backend/services"
	"department-eduvault-backend/utils"
	"fmt"

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
	role := c.GetString("role")
	facultyID := c.GetString("faculty_id")

	fmt.Printf("[DASHBOARD] Role=%s, FacultyID=%s\n", role, facultyID)

	overview, err := dc.service.GetOverview(c.Request.Context(), role, facultyID)
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
	role := c.GetString("role")
	facultyID := c.GetString("faculty_id")

	sections, err := dc.service.GetSectionStats(c.Request.Context(), role, facultyID)
	if err != nil {
		_ = c.Error(utils.NewDatabaseError("failed to load section statistics", err))
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    sections,
	})
}

// GetRecentActivity handles GET /dashboard/recent-activity
func (dc *DashboardController) GetRecentActivity(c *gin.Context) {
	role := c.GetString("role")
	facultyID := c.GetString("faculty_id")

	activities, err := dc.service.GetRecentActivity(c.Request.Context(), role, facultyID)
	if err != nil {
		_ = c.Error(utils.NewDatabaseError("failed to load recent activity", err))
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"data":    activities,
	})
}

// ExportMyCertificates handles GET /dashboard/export/certificates
func (dc *DashboardController) ExportMyCertificates(c *gin.Context) {
	facultyID := c.GetString("faculty_id")
	if facultyID == "" {
		_ = c.Error(utils.NewValidationError("faculty_id is required", nil))
		return
	}

	filename, content, err := dc.service.ExportMyCertificates(c.Request.Context(), facultyID)
	if err != nil {
		_ = c.Error(utils.NewDatabaseError("failed to export certificates", err))
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", content)
}
