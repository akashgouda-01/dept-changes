package controller

import (
	"net/http"

	"department-eduvault-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type HealthController struct {
	service service.HealthService
}

func NewHealthController(service service.HealthService) *HealthController {
	return &HealthController{service: service}
}

func (h *HealthController) Health(c *gin.Context) {
	if err := h.service.Check(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
