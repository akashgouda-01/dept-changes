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

// AuthMiddleware validates a Firebase ID token and enforces domain.
func AuthMiddleware(allowedDomain string) gin.HandlerFunc {
	// Lazy init
	initFirebase()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			_ = c.Error(utils.NewAuthenticationError("missing bearer token", nil))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Temporarily allow the old "mock" tokens for transition if needed,
		// OR strictly enforce Firebase. Let's strict enforce but handle the "dev mode" if firebase fails.
		if firebaseApp == nil {
			// Fallback for dev if firebase not set up
			// WARN: This is insecure, only for local dev without keys
			fmt.Println("[AUTH] WARNING: Firebase not initialized, allowing request (DEV MODE)")
			// Parse as mock if possible or just allow
			// For strict mode: return error.
			// Let's try to parse the mock format for backward compat during dev: email|role
			if strings.Contains(tokenString, "|") {
				MockFallback(c, tokenString, allowedDomain)
				return
			}
			_ = c.Error(utils.NewAuthenticationError("firebase not configured and invalid mock token", nil))
			c.Abort()
			return
		}

		ctx := context.Background()
		client, err := firebaseApp.Auth(ctx)
		if err != nil {
			_ = c.Error(utils.NewAuthenticationError("firebase auth client error", err))
			c.Abort()
			return
		}

		token, err := client.VerifyIDToken(ctx, tokenString)
		if err != nil {
			_ = c.Error(utils.NewAuthenticationError("invalid token", err))
			c.Abort()
			return
		}

		email := ""
		if e, ok := token.Claims["email"]; ok {
			email = e.(string)
		}

		if !strings.HasSuffix(strings.ToLower(email), "@"+strings.ToLower(allowedDomain)) {
			_ = c.Error(utils.NewAuthorizationError("email domain not allowed", nil))
			c.Abort()
			return
		}

		// Derive Identity (Role Logic mapped from existing)
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
		} else {
			facultyID = "FAC00" // Unknown
		}

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
