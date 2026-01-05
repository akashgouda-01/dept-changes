package utils

import (
	"log"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger returns a zap.Logger configured for production JSON output.
func NewLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	return logger
}

// NewRequestID generates a request identifier.
func NewRequestID() string {
	return uuid.NewString()
}
