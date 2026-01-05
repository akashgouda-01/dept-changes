package repositories

import (
	"context"
	"time"

	"department-eduvault-backend/models"
	"gorm.io/gorm"
)

type AdminRepository interface {
	SeedCertificatesIfEmpty(ctx context.Context) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

// SeedCertificatesIfEmpty inserts a small set of sample certificates if the table is empty.
// It is idempotent and safe to call multiple times.
func (r *adminRepository) SeedCertificatesIfEmpty(ctx context.Context) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Certificate{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Now().UTC()
	samples := []models.Certificate{
		{
			DriveLink:      "https://drive.google.com/file/d/sample-a-1",
			RegisterNumber: "RA2211003010",
			Section:        "A",
			StudentName:    "John Doe",
			UploadedBy:     "faculty@citchennai.net",
			UploadedAt:     now,
			MLStatus:       models.MLStatusVerified,
			FacultyStatus:  models.FacultyStatusLegit,
			Archived:       false,
		},
		{
			DriveLink:      "https://drive.google.com/file/d/sample-a-2",
			RegisterNumber: "RA2211003011",
			Section:        "A",
			StudentName:    "Jane Smith",
			UploadedBy:     "faculty@citchennai.net",
			UploadedAt:     now,
			MLStatus:       models.MLStatusVerified,
			FacultyStatus:  models.FacultyStatusPending,
			Archived:       false,
		},
		{
			DriveLink:      "https://drive.google.com/file/d/sample-b-1",
			RegisterNumber: "RA2211003020",
			Section:        "B",
			StudentName:    "Mike Johnson",
			UploadedBy:     "faculty@citchennai.net",
			UploadedAt:     now,
			MLStatus:       models.MLStatusVerified,
			FacultyStatus:  models.FacultyStatusNotLegit,
			Archived:       false,
		},
	}

	return r.db.WithContext(ctx).Create(&samples).Error
}
