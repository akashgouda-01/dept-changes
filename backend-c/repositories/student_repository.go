package repositories

import (
	"context"
	"department-eduvault-backend/models"

	"gorm.io/gorm"
)

type StudentRepository interface {
	GetStudentsBySection(ctx context.Context, section string) ([]models.Student, error)
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) GetStudentsBySection(ctx context.Context, section string) ([]models.Student, error) {
	var students []models.Student
	// Use LIKE to handle cases where input is "A" but DB has "Section A"
	// We match anything ending with the section letter/string or exact match
	if err := r.db.WithContext(ctx).Where("section LIKE ?", "%"+section).Find(&students).Error; err != nil {
		return nil, err
	}
	return students, nil
}
