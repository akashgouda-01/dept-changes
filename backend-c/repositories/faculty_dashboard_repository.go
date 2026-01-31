package repositories

import (
	"context"

	"gorm.io/gorm"
)

type FacultyDashboardStats struct {
	Total     int64
	Verified  int64
	Rejected  int64
	Pending   int64
}

type FacultyDashboardRepository struct {
	db *gorm.DB
}

func NewFacultyDashboardRepository(db *gorm.DB) *FacultyDashboardRepository {
	return &FacultyDashboardRepository{db}
}

func (r *FacultyDashboardRepository) GetStats(ctx context.Context, email string) (*FacultyDashboardStats, error) {
	var stats FacultyDashboardStats

	query := `
	SELECT
		COUNT(*) AS total,
		COUNT(*) FILTER (WHERE faculty_status='LEGIT') AS verified,
		COUNT(*) FILTER (WHERE faculty_status='NOT_LEGIT') AS rejected,
		COUNT(*) FILTER (WHERE faculty_status='PENDING') AS pending
	FROM certificates
	WHERE faculty_email = ?
	`

	err := r.db.WithContext(ctx).Raw(query, email).Scan(&stats).Error
	return &stats, err
}
