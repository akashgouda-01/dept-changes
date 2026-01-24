package middleware

import (
	"department-eduvault-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireRoles ensures the authenticated user has one of the allowed roles.
func RequireRoles(allowed ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := strings.ToLower(c.GetString("role"))
		for _, allowedRole := range allowed {
			if role == strings.ToLower(allowedRole) {
				c.Next()
				return
			}
		}
		_ = c.Error(utils.NewAuthorizationError("forbidden", nil))
		c.Abort()
	}
}
