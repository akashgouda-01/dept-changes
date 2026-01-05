package middleware

import (
	"net/http"
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
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
