package service

import (
	"context"
	"time"

	"department-eduvault-backend/internal/repository"
)

type HealthService interface {
	Check(ctx context.Context) error
}

type healthService struct {
	repo repository.HealthRepository
}

// NewHealthService wires the health service to its repository.
func NewHealthService(repo repository.HealthRepository) HealthService {
	return &healthService{repo: repo}
}

func (s *healthService) Check(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return s.repo.Ping(ctx)
}
