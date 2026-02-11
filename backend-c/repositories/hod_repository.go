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
	GetCertificatesByFacultyID(ctx context.Context, facultyID string) ([]models.Certificate, error)
	GetAllCertificates(ctx context.Context) ([]models.Certificate, error)
}

type hodRepository struct {
	db *gorm.DB
}

// NewHodRepository creates a HOD repository instance.
func NewHodRepository(db *gorm.DB) HodRepository {
	return &hodRepository{db: db}
}

// getFacultySections maps new Faculty IDs to their assigned sections.
// verified duplicates from dashboard_repository.go
func (r *hodRepository) getFacultySections(facultyID string) []string {
	switch facultyID {
	case "CSE245":
		return []string{"L", "M", "N", "O", "P", "Q"}
	case "CSE086":
		return []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q"}
	case "CSE262":
		// C,D,E,F,G,I,H -> Sorted: C, D, E, F, G, H, I
		return []string{"C", "D", "E", "F", "G", "H", "I"}
	case "CSE345":
		// A,B,I,J,K
		return []string{"A", "B", "I", "J", "K"}
	default:
		return nil
	}
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
			MAX(c.student_name) AS student_name,
			c.section,
			COUNT(*) AS total,
			COUNT(*) FILTER (WHERE c.faculty_status = 'LEGIT') AS verified,
			COUNT(*) FILTER (WHERE c.faculty_status = 'NOT_LEGIT') AS rejected,
			COUNT(*) FILTER (WHERE c.faculty_status = 'PENDING') AS pending
		FROM faculty_certificates c
		WHERE c.archived = false
	`

	args := []interface{}{}

	// Use section mapping to filter active students by section
	sections := r.getFacultySections(facultyID)
	if len(sections) > 0 {
		query += " AND c.section IN ?"
		args = append(args, sections)
	} else {
		// Fallback to strict faculty_id match if no section mapping exists
		query += " AND c.faculty_id = ?"
		args = append(args, facultyID)
	}

	query += `
		GROUP BY c.reg_no, c.section
		ORDER BY c.reg_no;
	`

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
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

// GetCertificatesByFacultyID returns all certificates uploaded by a specific faculty member.
func (r *hodRepository) GetCertificatesByFacultyID(ctx context.Context, facultyID string) ([]models.Certificate, error) {
	if facultyID == "" {
		return nil, fmt.Errorf("faculty_id is required")
	}

	var certs []models.Certificate
	query := r.db.WithContext(ctx).Where("archived = false")

	// Use section mapping if available
	sections := r.getFacultySections(facultyID)
	if len(sections) > 0 {
		query = query.Where("section IN ?", sections)
	} else {
		query = query.Where("faculty_id = ?", facultyID)
	}

	if err := query.Order("section ASC, reg_no ASC").Find(&certs).Error; err != nil {
		return nil, fmt.Errorf("query certificates by faculty: %w", err)
	}

	return certs, nil
}

// GetAllCertificates returns all active certificates in the system.
func (r *hodRepository) GetAllCertificates(ctx context.Context) ([]models.Certificate, error) {
	var certs []models.Certificate
	if err := r.db.WithContext(ctx).
		Where("archived = false").
		Order("section ASC, reg_no ASC").
		Find(&certs).Error; err != nil {
		return nil, fmt.Errorf("query all certificates: %w", err)
	}
	return certs, nil
}
