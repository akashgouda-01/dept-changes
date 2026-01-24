package repositories

import (
	"context"
	"fmt"

	"department-eduvault-backend/models"

	"gorm.io/gorm"
)

// StudentStatsRow represents aggregated certificate counts per student for a faculty.
type StudentStatsRow struct {
	RegisterNumber string
	StudentName    string
	Section        string
	Total          int64
	Verified       int64
	Rejected       int64
	Pending        int64
}

// HodRepository exposes queries used by HOD-facing APIs.
type HodRepository interface {
	GetStudentStatsByFaculty(ctx context.Context, facultyID string) ([]StudentStatsRow, error)
	GetCertificatesByStudent(ctx context.Context, regNo string) ([]models.Certificate, error)
	GetCertificatesBySection(ctx context.Context, section string) ([]models.Certificate, error)
}

type hodRepository struct {
	db *gorm.DB
}

// NewHodRepository creates a HOD repository instance.
func NewHodRepository(db *gorm.DB) HodRepository {
	return &hodRepository{db: db}
}

// GetStudentStatsByFaculty aggregates certificate counts per student for a faculty member.
func (r *hodRepository) GetStudentStatsByFaculty(ctx context.Context, facultyID string) ([]StudentStatsRow, error) {
	if facultyID == "" {
		return nil, fmt.Errorf("faculty id is required")
	}

	var rows []StudentStatsRow
	query := `
		SELECT
			c.reg_no AS register_number,
			CAST('' AS text) AS student_name,
			c.section,
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE c.faculty_status = 'LEGIT') AS verified,
			COUNT(*) FILTER (WHERE c.faculty_status = 'NOT_LEGIT') AS rejected,
			COUNT(*) FILTER (WHERE c.faculty_status = 'PENDING') AS pending
		FROM certificates c
		WHERE c.archived = false AND c.faculty_id = ?
		GROUP BY c.reg_no, c.section
		ORDER BY c.reg_no;
	`

	if err := r.db.WithContext(ctx).Raw(query, facultyID).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("query student stats: %w", err)
	}

	return rows, nil
}

// GetCertificatesByStudent returns certificates for a student by register number.
func (r *hodRepository) GetCertificatesByStudent(ctx context.Context, regNo string) ([]models.Certificate, error) {
	if regNo == "" {
		return nil, fmt.Errorf("reg_no is required")
	}

	var certs []models.Certificate
	if err := r.db.WithContext(ctx).
		Where("reg_no = ? AND archived = false", regNo).
		Order("uploaded_at DESC").
		Find(&certs).Error; err != nil {
		return nil, fmt.Errorf("query certificates by student: %w", err)
	}

	return certs, nil
}

// GetCertificatesBySection returns certificates for a section.
func (r *hodRepository) GetCertificatesBySection(ctx context.Context, section string) ([]models.Certificate, error) {
	if section == "" {
		return nil, fmt.Errorf("section is required")
	}

	var certs []models.Certificate
	if err := r.db.WithContext(ctx).
		Where("section = ? AND archived = false", section).
		Order("uploaded_at DESC").
		Find(&certs).Error; err != nil {
		return nil, fmt.Errorf("query certificates by section: %w", err)
	}

	return certs, nil
}
