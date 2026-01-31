package middleware

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"department-eduvault-backend/utils"

	"github.com/gin-gonic/gin"
)

// MockAuthMiddleware validates a Google OAuth bearer token (placeholder) and enforces domain.
// In production, replace token validation with real Google token verification.
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

		// parseMockToken simulates token parsing; expected format: "email"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		// In previous logic, token might be "email|role" but client now sends just "email".
		// We handle purely by email now as requested.
		// If the token contains pipe, we extract just the email (first part).
		parts := strings.Split(token, "|")
		email := strings.TrimSpace(parts[0])

		if !strings.HasSuffix(strings.ToLower(email), "@"+strings.ToLower(allowedDomain)) {
			_ = c.Error(utils.NewAuthorizationError("email domain not allowed", nil))
			c.Abort()
			return
		}

		// Derive Identity
		role := "faculty"
		facultyID := ""

		lowerEmail := strings.ToLower(email)
		if lowerEmail == "hod@citchennai.net" {
			role = "hod"
		} else if strings.HasPrefix(lowerEmail, "faculty1") {
			facultyID = "FAC01" // Sections A
		} else if strings.HasPrefix(lowerEmail, "faculty2") {
			facultyID = "FAC02" // Sections B, C, D
		} else {
			// Default or unknown faculty
			facultyID = "FAC03" // fallback
		}

		// Debug logging
		logFile, _ := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if logFile != nil {
			fmt.Fprintf(logFile, "[AUTH] Email=%s, Role=%s, FacultyID=%s\n", email, role, facultyID)
			logFile.Close()
		}

		c.Set("email", email)
		c.Set("role", role)
		c.Set("faculty_id", facultyID)
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
