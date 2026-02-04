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

// SeedCertificatesIfEmpty inserts a small set of sample certificates and STUDENTS if the table is empty.
// It is idempotent and safe to call multiple times.
func (r *adminRepository) SeedCertificatesIfEmpty(ctx context.Context) error {
	// 1. Seed Students if empty
	var studentCount int64
	if err := r.db.WithContext(ctx).Model(&models.Student{}).Count(&studentCount).Error; err != nil {
		return err
	}

	if studentCount == 0 {
		sections := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q"}
		// Logic: 6 sections per faculty
		// Faculty 1: A-F
		// Faculty 2: G-L
		// Faculty 3: M-Q

		students := []models.Student{}

		for _, section := range sections {
			// Determine faculty
			var facultyEmail string
			switch {
			case section >= "A" && section <= "F":
				facultyEmail = "faculty1@citchennai.net"
			case section >= "G" && section <= "L":
				facultyEmail = "faculty2@citchennai.net"
			default:
				facultyEmail = "faculty3@citchennai.net"
			}

			// Create 2 dummy students per section
			for i := 1; i <= 2; i++ {
				students = append(students, models.Student{
					RegisterNumber: "REG" + section + "00" + string(rune('0'+i)), // e.g., REGA001
					Name:           "Student " + section + "-" + string(rune('0'+i)),
					Email:          "student" + section + string(rune('0'+i)) + "@citchennai.net",
					Section:        section,
					Semester:       5,
					IsPresent:      true,
					FacultyEmail:   facultyEmail,
				})
			}
		}

		if err := r.db.WithContext(ctx).Create(&students).Error; err != nil {
			return err
		}
	}

	// 2. Seed Certificates
	var certCount int64
	if err := r.db.WithContext(ctx).Model(&models.Certificate{}).Count(&certCount).Error; err != nil {
		return err
	}
	if certCount > 0 {
		return nil
	}

	// We'll attach certificates to the first few seeded students if possible,
	// or just use generic ones matching the schema.
	// Matching valid reg numbers ensures constraints (if any) are met.

	now := time.Now().UTC()
	samples := []models.Certificate{
		{
			DriveLink:      "https://drive.google.com/file/d/sample-a-1",
			RegisterNumber: "REGA001", // Matches seeded student
			Section:        "A",
			StudentName:    "Student A-1",
			UploadedBy:     "faculty1@citchennai.net",
			UploadedAt:     now,
			MLStatus:       models.MLStatusVerified,
			FacultyStatus:  models.FacultyStatusLegit,
			Archived:       false,
		},
		{
			DriveLink:      "https://drive.google.com/file/d/sample-a-2",
			RegisterNumber: "REGA002",
			Section:        "A",
			StudentName:    "Student A-2",
			UploadedBy:     "faculty1@citchennai.net",
			UploadedAt:     now,
			MLStatus:       models.MLStatusVerified,
			FacultyStatus:  models.FacultyStatusPending,
			Archived:       false,
		},
		{
			DriveLink:      "https://drive.google.com/file/d/sample-g-1",
			RegisterNumber: "REGG001",
			Section:        "G",
			StudentName:    "Student G-1",
			UploadedBy:     "faculty2@citchennai.net",
			UploadedAt:     now,
			MLStatus:       models.MLStatusVerified,
			FacultyStatus:  models.FacultyStatusPending,
			Archived:       false,
		},
	}

	return r.db.WithContext(ctx).Create(&samples).Error
}
