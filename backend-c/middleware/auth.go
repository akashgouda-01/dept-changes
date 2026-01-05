package middleware

import (
	"errors"
	"strings"

	"department-eduvault-backend/utils"
	"github.com/gin-gonic/gin"
)

// MockAuthMiddleware validates a Google OAuth bearer token (placeholder) and enforces domain.
// In production, replace token validation with real Google token verification.
func MockAuthMiddleware(allowedDomain string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			_ = c.Error(utils.NewAuthenticationError("missing bearer token", nil))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		// TODO: Replace this mock parsing with real Google token introspection.
		email, role, err := parseMockToken(token)
		if err != nil {
			_ = c.Error(utils.NewAuthenticationError("invalid token", err))
			c.Abort()
			return
		}

		if !strings.HasSuffix(strings.ToLower(email), "@"+strings.ToLower(allowedDomain)) {
			_ = c.Error(utils.NewAuthorizationError("email domain not allowed", nil))
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("role", strings.ToLower(role))
		c.Next()
	}
}

// parseMockToken simulates token parsing; expected format: "email|role".
func parseMockToken(token string) (string, string, error) {
	parts := strings.Split(token, "|")
	if len(parts) != 2 {
		return "", "", errInvalidToken
	}
	email := strings.TrimSpace(parts[0])
	role := strings.TrimSpace(parts[1])
	if email == "" || role == "" {
		return "", "", errInvalidToken
	}
	return email, role, nil
}

var errInvalidToken = errors.New("invalid token")
