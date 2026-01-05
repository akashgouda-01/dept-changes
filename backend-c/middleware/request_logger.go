package middleware

import (
	"time"

	"department-eduvault-backend/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger logs structured request metadata and injects a request ID into context.
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqID := utils.NewRequestID()
		c.Set("request_id", reqID)

		c.Next()

		duration := time.Since(start)
		logger.Info("request completed",
			zap.String("request_id", reqID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Float64("duration_ms", float64(duration.Milliseconds())),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}
