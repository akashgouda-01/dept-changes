package utils

import (
	"errors"
	"net/http"
)

// AppError represents a typed application error with an associated HTTP status code.
type AppError struct {
	Code    string
	Message string
	Status  int
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Factory helpers ------------------------------------------------------------

func NewValidationError(message string, err error) *AppError {
	return &AppError{Code: "VALIDATION_ERROR", Message: message, Status: http.StatusBadRequest, Err: err}
}

func NewAuthenticationError(message string, err error) *AppError {
	return &AppError{Code: "AUTHENTICATION_ERROR", Message: message, Status: http.StatusUnauthorized, Err: err}
}

func NewAuthorizationError(message string, err error) *AppError {
	return &AppError{Code: "AUTHORIZATION_ERROR", Message: message, Status: http.StatusForbidden, Err: err}
}

func NewNotFoundError(message string, err error) *AppError {
	return &AppError{Code: "NOT_FOUND", Message: message, Status: http.StatusNotFound, Err: err}
}

func NewDatabaseError(message string, err error) *AppError {
	return &AppError{Code: "DATABASE_ERROR", Message: message, Status: http.StatusInternalServerError, Err: err}
}

func NewInternalError(message string, err error) *AppError {
	return &AppError{Code: "INTERNAL_SERVER_ERROR", Message: message, Status: http.StatusInternalServerError, Err: err}
}

// IsAppError helps identify AppError via errors.As.
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}
