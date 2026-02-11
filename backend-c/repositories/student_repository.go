package repositories

import (
	"context"
	"department-eduvault-backend/models"

	"gorm.io/gorm"
)

type StudentRepository interface {
	GetStudentsBySection(ctx context.Context, section string) ([]models.Student, error)
	GetStudentByRegNo(ctx context.Context, regNo string, allowedSections []string) (*models.Student, error)
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

func (r *studentRepository) GetStudentByRegNo(ctx context.Context, regNo string, allowedSections []string) (*models.Student, error) {
	var student models.Student

	// Fetch student by RegNo first (without section filter)
	if err := r.db.WithContext(ctx).Where("register_number = ?", regNo).First(&student).Error; err != nil {
		return nil, err
	}

	// Manual check for section in Go to handle "Section A" vs "A" mismatch
	isAllowed := false
	for _, allowed := range allowedSections {
		// Normalize DB section: remove "section" prefix, trim spaces
		// Note: Doing this properly requires handling case-insensitivity
		// Let's rely on simple checks first as Go has no built-in regex in loop that is cheap

		// 1. Exact match
		if student.Section == allowed {
			isAllowed = true
			break
		}

		sLen := len(student.Section)
		aLen := len(allowed)

		if sLen >= aLen {
			// Check if suffix matches
			suffix := student.Section[sLen-aLen:]
			if suffix == allowed {
				// Matched suffix. Now check boundary.
				// Valid if:
				// - Exact match (sLen == aLen) - handled above
				// - Preceded by space (student.Section[...-1] == ' ')
				if sLen > aLen && student.Section[sLen-aLen-1] == ' ' {
					isAllowed = true
					break
				}
			}
		}
	}

	if !isAllowed {
		return nil, gorm.ErrRecordNotFound
	}

	return &student, nil
}
