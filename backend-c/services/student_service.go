package services

import (
	"context"
	"department-eduvault-backend/models"
	"department-eduvault-backend/repositories"
)

type StudentService interface {
	GetStudentsBySection(ctx context.Context, section string) ([]models.Student, error)
}

type studentService struct {
	repo repositories.StudentRepository
}

func NewStudentService(repo repositories.StudentRepository) StudentService {
	return &studentService{repo: repo}
}

func (s *studentService) GetStudentsBySection(ctx context.Context, section string) ([]models.Student, error) {
	return s.repo.GetStudentsBySection(ctx, section)
}
