package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("YOUR_SUPER_SECRET_KEY_CHANGE_IN_PROD") // In real app, load from env

type Claims struct {
	Email     string `json:"email"`
	Role      string `json:"role"`
	FacultyID string `json:"faculty_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a signed JWT for the user
func GenerateToken(email, role, facultyID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email:     email,
		Role:      role,
		FacultyID: facultyID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

// ValidateToken parses and validates the signed JWT
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
