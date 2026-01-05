package repositories

import (
	"context"

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
	GetOverview(ctx context.Context) (DashboardOverview, error)
	GetSectionStats(ctx context.Context) ([]SectionDashboardRow, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetOverview(ctx context.Context) (DashboardOverview, error) {
	// Aggregate solely from certificates table using COUNT + CASE expressions.
	type aggRow struct {
		TotalStudents     int64
		TotalCertificates int64
		VerifiedCount     int64
		RejectedCount     int64
		PendingCount      int64
	}

	var row aggRow

	query := `
		SELECT
			COALESCE(COUNT(DISTINCT reg_no), 0) AS total_students,
			COALESCE(COUNT(*), 0) AS total_certificates,
			COALESCE(COUNT(CASE WHEN faculty_status = 'LEGIT' THEN 1 END), 0) AS verified_count,
			COALESCE(COUNT(CASE WHEN faculty_status = 'NOT_LEGIT' THEN 1 END), 0) AS rejected_count,
			COALESCE(COUNT(CASE WHEN faculty_status = 'PENDING' THEN 1 END), 0) AS pending_count
		FROM certificates
		WHERE archived = false;
	`

	if err := r.db.WithContext(ctx).Raw(query).Scan(&row).Error; err != nil {
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

func (r *dashboardRepository) GetSectionStats(ctx context.Context) ([]SectionDashboardRow, error) {
	var rows []SectionDashboardRow

	query := `
		SELECT
			section AS section,
			COALESCE(COUNT(*), 0) AS total_certificates,
			COALESCE(COUNT(CASE WHEN faculty_status = 'LEGIT' THEN 1 END), 0) AS verified_certificates,
			COALESCE(COUNT(CASE WHEN faculty_status = 'NOT_LEGIT' THEN 1 END), 0) AS rejected_certificates,
			COALESCE(COUNT(CASE WHEN faculty_status = 'PENDING' THEN 1 END), 0) AS pending_certificates
		FROM certificates
		WHERE archived = false
		GROUP BY section
		ORDER BY section;
	`

	if err := r.db.WithContext(ctx).Raw(query).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}
