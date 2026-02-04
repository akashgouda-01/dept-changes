package repositories

import (
	"context"
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

type DashboardRepository interface {
	GetOverview(ctx context.Context, facultyID string) (DashboardOverview, error)
	GetSectionStats(ctx context.Context, facultyID string) ([]SectionDashboardRow, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
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
			COALESCE(COUNT(CASE WHEN faculty_status = 'LEGIT' THEN 1 END), 0) AS verified_count,
			COALESCE(COUNT(CASE WHEN faculty_status = 'NOT_LEGIT' THEN 1 END), 0) AS rejected_count,
			COALESCE(COUNT(CASE WHEN faculty_status = 'PENDING' AND ml_status = 'VERIFIED' THEN 1 END), 0) AS pending_count
		FROM faculty_certificates
		WHERE archived = false
	`

	if facultyID != "" {
		// As per instructions, filtered by faculty_id
		query += " AND faculty_id = ?"
		args = append(args, facultyID)
	}

	query += ";"

	fmt.Printf("[REPO] Query=%s, Args=%v\n", query, args)

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&row).Error; err != nil {
		return DashboardOverview{}, err
	}

	fmt.Printf("[REPO] Result: TotalStudents=%d, TotalCerts=%d, Verified=%d, Rejected=%d, Pending=%d\n",
		row.TotalStudents, row.TotalCertificates, row.VerifiedCount, row.RejectedCount, row.PendingCount)

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
			COALESCE(COUNT(CASE WHEN faculty_status = 'LEGIT' THEN 1 END), 0) AS verified_certificates,
			COALESCE(COUNT(CASE WHEN faculty_status = 'NOT_LEGIT' THEN 1 END), 0) AS rejected_certificates,
			COALESCE(COUNT(CASE WHEN faculty_status = 'PENDING' AND ml_status = 'VERIFIED' THEN 1 END), 0) AS pending_certificates
		FROM faculty_certificates
		WHERE archived = false
	`

	if facultyID != "" {
		query += " AND faculty_id = ?"
		args = append(args, facultyID)
	}

	query += " GROUP BY section ORDER BY section;"

	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}
