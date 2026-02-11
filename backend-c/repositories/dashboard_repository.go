package repositories

import (
	"context"
	"department-eduvault-backend/models"
	"fmt"

	"gorm.io/gorm"
)

type DashboardOverview struct {
	TotalStudents        int64
	TotalCertificates    int64
	VerifiedCertificates int64
	RejectedCertificates int64
	PendingCertificates  int64
}

type SectionDashboardRow struct {
	Section              string
	TotalCertificates    int64
	VerifiedCertificates int64
	RejectedCertificates int64
	PendingCertificates  int64
}

type RecentActivityRow struct {
	ID          string
	StudentName string
	RegNo       string
	Section     string
	Action      string
	Timestamp   string
}

type DashboardRepository interface {
	GetOverview(ctx context.Context, facultyID string) (DashboardOverview, error)
	GetSectionStats(ctx context.Context, facultyID string) ([]SectionDashboardRow, error)
	GetRecentActivity(ctx context.Context, facultyID string) ([]RecentActivityRow, error)
	GetCertificatesByFacultyID(ctx context.Context, facultyID string) ([]models.Certificate, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) getFacultySections(facultyID string) []string {
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

func (r *dashboardRepository) GetOverview(ctx context.Context, facultyID string) (DashboardOverview, error) {
	// Aggregate from certificates table.
	type aggRow struct {
		TotalStudents     int64
		TotalCertificates int64
		VerifiedCount     int64
		RejectedCount     int64
		PendingCount      int64
	}

	var row aggRow
	args := []interface{}{}

	query := `
		SELECT
			COALESCE(COUNT(DISTINCT reg_no), 0) AS total_students,
			COALESCE(COUNT(*), 0) AS total_certificates,
			COALESCE(SUM(CASE WHEN faculty_status = 'LEGIT' THEN 1 ELSE 0 END), 0) AS verified_count,
			COALESCE(SUM(CASE WHEN faculty_status = 'NOT_LEGIT' THEN 1 ELSE 0 END), 0) AS rejected_count,
			COALESCE(SUM(CASE WHEN faculty_status = 'PENDING' THEN 1 ELSE 0 END), 0) AS pending_count
		FROM faculty_certificates
		WHERE archived = false
	`

	// If facultyID matches one of our known mapped faculties, filter by their sections.
	// Otherwise, fallback to faculty_id column filter (or no filter if empty).
	sections := r.getFacultySections(facultyID)
	if len(sections) > 0 {
		query += " AND section IN ?"
		args = append(args, sections)
	} else if facultyID != "" {
		query += " AND faculty_id = ?"
		args = append(args, facultyID)
	}

	query += ";"

	fmt.Printf("[REPO] GetOverview Query=%s, Args=%v\n", query, args)

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&row).Error; err != nil {
		return DashboardOverview{}, err
	}

	return DashboardOverview{
		TotalStudents:        row.TotalStudents,
		TotalCertificates:    row.TotalCertificates,
		VerifiedCertificates: row.VerifiedCount,
		RejectedCertificates: row.RejectedCount,
		PendingCertificates:  row.PendingCount,
	}, nil
}

func (r *dashboardRepository) GetSectionStats(ctx context.Context, facultyID string) ([]SectionDashboardRow, error) {
	var rows []SectionDashboardRow
	args := []interface{}{}

	query := `
		SELECT
			section AS section,
			COALESCE(COUNT(*), 0) AS total_certificates,
			COALESCE(SUM(CASE WHEN faculty_status = 'LEGIT' THEN 1 ELSE 0 END), 0) AS verified_certificates,
			COALESCE(SUM(CASE WHEN faculty_status = 'NOT_LEGIT' THEN 1 ELSE 0 END), 0) AS rejected_certificates,
			COALESCE(SUM(CASE WHEN faculty_status = 'PENDING' THEN 1 ELSE 0 END), 0) AS pending_certificates
		FROM faculty_certificates
		WHERE archived = false
	`

	sections := r.getFacultySections(facultyID)
	if len(sections) > 0 {
		query += " AND section IN ?"
		args = append(args, sections)
	} else if facultyID != "" {
		query += " AND faculty_id = ?"
		args = append(args, facultyID)
	}

	query += " GROUP BY section ORDER BY section;"

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *dashboardRepository) GetRecentActivity(ctx context.Context, facultyID string) ([]RecentActivityRow, error) {
	// We want activities from "Today".
	type dbRow struct {
		ID            string
		StudentName   string
		RegNo         string
		Section       string
		FacultyStatus string
		UploadedAt    string
		UpdatedAt     string
	}

	var results []dbRow
	args := []interface{}{}

	query := `
		SELECT 
			id, student_name, reg_no, section, faculty_status::text,
			uploaded_at::text, updated_at::text
		FROM faculty_certificates
		WHERE archived = false
		AND (DATE(updated_at) = CURRENT_DATE OR DATE(uploaded_at) = CURRENT_DATE)
	`

	sections := r.getFacultySections(facultyID)
	if len(sections) > 0 {
		query += " AND section IN ?"
		args = append(args, sections)
	} else if facultyID != "" {
		query += " AND faculty_id = ?"
		args = append(args, facultyID)
	}

	query += " ORDER BY updated_at DESC LIMIT 50;"

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&results).Error; err != nil {
		return nil, err
	}

	var activities []RecentActivityRow
	for _, row := range results {
		action := "UPLOADED"
		ts := row.UploadedAt

		if row.FacultyStatus == "LEGIT" {
			action = "VERIFIED"
			ts = row.UpdatedAt
		} else if row.FacultyStatus == "NOT_LEGIT" {
			action = "REJECTED"
			ts = row.UpdatedAt
		} else {
			action = "UPLOADED"
			ts = row.UploadedAt
		}

		activities = append(activities, RecentActivityRow{
			ID:          row.ID,
			StudentName: row.StudentName,
			RegNo:       row.RegNo,
			Section:     row.Section,
			Action:      action,
			Timestamp:   ts,
		})
	}

	return activities, nil
}

func (r *dashboardRepository) GetCertificatesByFacultyID(ctx context.Context, facultyID string) ([]models.Certificate, error) {
	if facultyID == "" {
		return nil, fmt.Errorf("faculty_id is required")
	}

	var certs []models.Certificate
	query := r.db.WithContext(ctx).Where("archived = false")

	// If we have mapped sections, filter by those.
	// Otherwise, fallback to database faculty_id.
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
