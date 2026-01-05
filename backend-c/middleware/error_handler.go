package middleware

import (
	"errors"

	"department-eduvault-backend/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler centralizes error-to-response mapping.
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		errs := c.Errors
		if len(errs) == 0 {
			return
		}

		// Use the first error for response; log all.
		appErr := normalizeError(errs[0].Err)

		logger.Error("request error",
			zap.Error(errs[0].Err),
			zap.String("request_id", c.GetString("request_id")),
			zap.String("code", appErr.Code),
		)

		c.JSON(appErr.Status, gin.H{
			"success": false,
			"error": gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
		c.Abort()
	}
}

func normalizeError(err error) *utils.AppError {
	if err == nil {
		return utils.NewInternalError("unknown error", nil)
	}
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return utils.NewInternalError("internal server error", err)
}
