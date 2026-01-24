package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"department-eduvault-backend/models"
	"department-eduvault-backend/repositories"
	"department-eduvault-backend/services"
	"department-eduvault-backend/utils"

	"github.com/gin-gonic/gin"
)

// CertificateController handles HTTP requests for certificate workflows.
type CertificateController struct {
	service services.CertificateService
}

// NewCertificateController creates a new controller instance.
func NewCertificateController(service services.CertificateService) *CertificateController {
	return &CertificateController{service: service}
}

// UploadCertificates handles POST /certificates/upload
func (cc *CertificateController) UploadCertificates(c *gin.Context) {
	if !cc.requireRole(c, "faculty") {
		return
	}

	var req uploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(utils.NewValidationError("invalid request payload", err))
		return
	}

	if len(req.Certificates) == 0 {
		_ = c.Error(utils.NewValidationError("certificates are required", nil))
		return
	}
	if len(req.Certificates) > 10 {
		_ = c.Error(utils.NewValidationError("cannot upload more than 10 certificates", nil))
		return
	}

	inputs := make([]services.CertificateInput, 0, len(req.Certificates))
	for _, item := range req.Certificates {
		if err := item.validate(); err != nil {
			_ = c.Error(utils.NewValidationError(err.Error(), err))
			return
		}
		inputs = append(inputs, services.CertificateInput{
			DriveLink:      item.DriveLink,
			RegisterNumber: item.RegisterNumber,
			Section:        item.Section,
			StudentName:    item.StudentName,
			UploadedBy:     item.UploadedBy,
			UploadedAt:     item.UploadedAt,
		})
	}

	if err := cc.service.UploadCertificates(c.Request.Context(), inputs); err != nil {
		_ = c.Error(mapServiceError(err))
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "certificates accepted for processing"})
}

// GetPendingReview handles GET /certificates/pending-review
func (cc *CertificateController) GetPendingReview(c *gin.Context) {
	if !cc.requireRole(c, "faculty") {
		return
	}

	limit := 50
	if v := c.Query("limit"); v != "" {
		parsed, err := parsePositiveInt(v)
		if err != nil {
			_ = c.Error(utils.NewValidationError("limit must be a positive integer", err))
			return
		}
		limit = parsed
	}

	certs, err := cc.service.GetPendingFacultyReview(c.Request.Context(), limit)
	if err != nil {
		_ = c.Error(mapServiceError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": certs})
}

// SubmitReview handles POST /certificates/review
func (cc *CertificateController) SubmitReview(c *gin.Context) {
	if !cc.requireRole(c, "faculty", "hod") {
		return
	}

	var req reviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(utils.NewValidationError("invalid request payload", err))
		return
	}
	if err := req.validate(); err != nil {
		_ = c.Error(utils.NewValidationError(err.Error(), err))
		return
	}

	if err := cc.service.SubmitFacultyDecision(c.Request.Context(), req.CertificateID, req.Status, req.IsLegit); err != nil {
		_ = c.Error(mapServiceError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review recorded"})
}

// TriggerMockVerification handles:
// POST /faculty/certificate/verify
// This endpoint only triggers the existing mock ML verification and returns
// a static response; no real ML integration is performed here.
func (cc *CertificateController) TriggerMockVerification(c *gin.Context) {
	if !cc.requireRole(c, "faculty", "hod") {
		return
	}

	var payload struct {
		CertificateID string `json:"certificate_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || strings.TrimSpace(payload.CertificateID) == "" {
		_ = c.Error(utils.NewValidationError("certificate_id is required", err))
		return
	}

	if err := cc.service.TriggerMockMLVerification(c.Request.Context(), payload.CertificateID); err != nil {
		_ = c.Error(mapServiceError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "mock ML verification triggered",
	})
}

// Helpers --------------------------------------------------------------------

func (cc *CertificateController) requireRole(c *gin.Context, allowedRoles ...string) bool {
	role := strings.ToLower(c.GetString("role"))
	for _, allowed := range allowedRoles {
		if role == strings.ToLower(allowed) {
			return true
		}
	}
	_ = c.Error(utils.NewAuthorizationError("forbidden", nil))
	c.Abort()
	return false
}

type uploadRequest struct {
	Certificates []uploadItem `json:"certificates" binding:"required"`
}

type uploadItem struct {
	DriveLink      string    `json:"drive_link" binding:"required"`
	RegisterNumber string    `json:"register_number" binding:"required"`
	Section        string    `json:"section" binding:"required"`
	StudentName    string    `json:"student_name" binding:"required"`
	UploadedBy     string    `json:"uploaded_by" binding:"required"`
	UploadedAt     time.Time `json:"uploaded_at"`
}

func (u uploadItem) validate() error {
	if u.DriveLink == "" || u.RegisterNumber == "" || u.Section == "" || u.StudentName == "" || u.UploadedBy == "" {
		return errors.New("all certificate fields are required")
	}
	return nil
}

type reviewRequest struct {
	CertificateID string               `json:"certificate_id" binding:"required"`
	Status        models.FacultyStatus `json:"status" binding:"required"`
	IsLegit       bool                 `json:"is_legit"`
}

func (r reviewRequest) validate() error {
	if r.CertificateID == "" {
		return errors.New("certificate_id is required")
	}
	if r.Status != models.FacultyStatusLegit && r.Status != models.FacultyStatusNotLegit {
		return errors.New("status must be either 'legit' or 'not_legit'")
	}
	return nil
}

func mapServiceError(err error) *utils.AppError {
	switch {
	case errors.Is(err, services.ErrInvalidDriveLink):
		return utils.NewValidationError(err.Error(), err)
	case errors.Is(err, services.ErrUploadLimitExceeded):
		return utils.NewValidationError(err.Error(), err)
	case errors.Is(err, services.ErrInvalidMLTransition):
		return utils.NewValidationError(err.Error(), err)
	case errors.Is(err, services.ErrInvalidFacultyState):
		return utils.NewValidationError(err.Error(), err)
	case errors.Is(err, services.ErrCertificateArchived):
		return utils.NewAuthorizationError(err.Error(), err)
	case errors.Is(err, repositories.ErrCertificateNotFound):
		return utils.NewNotFoundError(err.Error(), err)
	case errors.Is(err, repositories.ErrStatsNotFound):
		// Use 404 when stats not found, implying student/section not found for update
		return utils.NewNotFoundError("related statistics record not found", err)
	default:
		return utils.NewInternalError("internal server error", err)
	}
}

func parsePositiveInt(val string) (int, error) {
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	if n <= 0 {
		return 0, errors.New("value must be positive")
	}
	return n, nil
}
