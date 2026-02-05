package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"department-eduvault-backend/utils"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App

func initFirebase() {
	if firebaseApp != nil {
		return
	}
	ctx := context.Background()
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	// If file doesn't exist, it might fallback to env GOOGLE_APPLICATION_CREDENTIALS
	// We try initializing.
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		// Fallback: Try without specific file (env var)
		app, err = firebase.NewApp(ctx, nil)
		if err != nil {
			fmt.Printf("[AUTH] Failed to initialize Firebase App: %v\n", err)
			return
		}
	}
	firebaseApp = app
	fmt.Println("[AUTH] Firebase App Initialized")
}

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

		// parseMockToken simulates token parsing; expected format: "email|role"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		parts := strings.Split(token, "|")
		if len(parts) < 1 {
			_ = c.Error(utils.NewAuthenticationError("invalid mock token format", nil))
			c.Abort()
			return
		}
		email := strings.TrimSpace(parts[0])
		lowerEmail := strings.ToLower(email)

		// Define Allowlist & Role Mapping
		// Map email -> struct{Role, FacultyID}
		type UserConfig struct {
			Role      string
			FacultyID string
		}

		accessMap := map[string]UserConfig{
			"harjeetp.cse2024@citchennai.net":        {Role: "faculty", FacultyID: "FAC01"},
			"hemanm.cse2024@citchennai.net":          {Role: "faculty", FacultyID: "FAC02"},
			"akashkumargouda.cse2024@citchennai.net": {Role: "faculty", FacultyID: "FAC03"},
			"aadhishs.cse2024@citchennai.net":        {Role: "hod", FacultyID: ""},
		}

		config, allowed := accessMap[lowerEmail]
		if !allowed {
			_ = c.Error(utils.NewAuthorizationError("email is not authorized", nil))
			c.Abort()
			return
		}

		// Derive Identity
		role := config.Role
		facultyID := config.FacultyID

		c.Set("email", email)
		c.Set("role", role)
		c.Set("faculty_id", facultyID)
		c.Next()
	}
}

func MockFallback(c *gin.Context, token string, allowedDomain string) {
	// Re-implementation of the old mock logic for continuity if Firebase fails
	parts := strings.Split(token, "|")
	email := strings.TrimSpace(parts[0])

	if !strings.HasSuffix(strings.ToLower(email), "@"+strings.ToLower(allowedDomain)) {
		_ = c.Error(utils.NewAuthorizationError("email domain not allowed", nil))
		c.Abort()
		return
	}

	role := "faculty"
	facultyID := ""
	lowerEmail := strings.ToLower(email)
	if strings.HasPrefix(lowerEmail, "hod") {
		role = "hod"
	} else if strings.HasPrefix(lowerEmail, "faculty1") {
		facultyID = "FAC01"
	} else if strings.HasPrefix(lowerEmail, "faculty2") {
		facultyID = "FAC02"
	} else if strings.HasPrefix(lowerEmail, "faculty3") {
		facultyID = "FAC03"
	}

	c.Set("email", email)
	c.Set("role", role)
	c.Set("faculty_id", facultyID)
	c.Next()
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
