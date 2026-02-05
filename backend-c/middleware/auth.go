package middleware

import (
	"strings"

	"department-eduvault-backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates the JWT token signed by the backend.
func AuthMiddleware(allowedDomain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			_ = c.Error(utils.NewAuthenticationError("missing bearer token", nil))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify JWT using the backend's secret key
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			_ = c.Error(utils.NewAuthenticationError("invalid or expired token", err))
			c.Abort()
			return
		}

		email := claims.Email
		if !strings.HasSuffix(strings.ToLower(email), "@"+strings.ToLower(allowedDomain)) {
			_ = c.Error(utils.NewAuthorizationError("email domain not allowed", nil))
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("faculty_id", claims.FacultyID)
		c.Next()
	}
}

// Keep MockAuthMiddleware for backward compatibility if needed, calling AuthMiddleware
func MockAuthMiddleware(allowedDomain string) gin.HandlerFunc {
	return AuthMiddleware(allowedDomain)
}
