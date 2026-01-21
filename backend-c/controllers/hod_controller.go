package controllers

import (
	"net/http"

	"department-eduvault-backend/services"
	"department-eduvault-backend/utils"
	"github.com/gin-gonic/gin"
)

// HodController exposes HOD-facing APIs.
type HodController struct {
	service services.HodService
}

// NewHodController constructs a HodController.
func NewHodController(service services.HodService) *HodController {
	return &HodController{service: service}
}

// HodDashboard handles:
// GET /hod/dashboard
func (hc *HodController) HodDashboard(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "HOD Dashboard API",
	})
}

// GetStudentStats handles:
// GET /hod/faculty/students?faculty_id=FAC01
func (hc *HodController) GetStudentStats(c *gin.Context) {
	facultyID := c.Query("faculty_id")
	if facultyID == "" {
		_ = c.Error(utils.NewValidationError("faculty_id is required", nil))
		return
	}

	stats, err := hc.service.GetStudentStatsByFaculty(c.Request.Context(), facultyID)
	if err != nil {
		_ = c.Error(utils.NewDatabaseError("failed to load student statistics", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// ListStudentCertificates handles:
// GET /hod/student/certificates?reg_no=XXXX
func (hc *HodController) ListStudentCertificates(c *gin.Context) {
	regNo := c.Query("reg_no")
	if regNo == "" {
		_ = c.Error(utils.NewValidationError("reg_no is required", nil))
		return
	}

	certs, err := hc.service.ListStudentCertificates(c.Request.Context(), regNo)
	if err != nil {
		_ = c.Error(utils.NewDatabaseError("failed to load certificates for student", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    certs,
	})
}

// ExportCertificatesBySection handles:
// GET /hod/export/certificates/section?section=SECTION
func (hc *HodController) ExportCertificatesBySection(c *gin.Context) {
	section := c.Query("section")
	if section == "" {
		_ = c.Error(utils.NewValidationError("section is required", nil))
		return
	}

	filename, content, err := hc.service.ExportCertificatesBySection(c.Request.Context(), section)
	if err != nil {
		_ = c.Error(utils.NewDatabaseError("failed to export certificates by section", err))
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", content)
}

// ExportCertificatesByStudent handles:
// GET /hod/export/certificates/student?reg_no=XXXX
func (hc *HodController) ExportCertificatesByStudent(c *gin.Context) {
	regNo := c.Query("reg_no")
	if regNo == "" {
		_ = c.Error(utils.NewValidationError("reg_no is required", nil))
		return
	}

	filename, content, err := hc.service.ExportCertificatesByStudent(c.Request.Context(), regNo)
	if err != nil {
		_ = c.Error(utils.NewDatabaseError("failed to export certificates by student", err))
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", content)
}

