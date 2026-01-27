package repositories

import (
	"context"
	"errors"
	"fmt"

	"department-eduvault-backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	// ErrTooManyCertificates is returned when bulk insert exceeds the allowed batch size.
	ErrTooManyCertificates = errors.New("cannot insert more than 10 certificates at once")
	// ErrCertificateNotFound is returned when a certificate lookup fails.
	ErrCertificateNotFound = errors.New("certificate not found")
	// ErrStatsNotFound indicates related stats rows are missing for updates.
	ErrStatsNotFound = errors.New("statistics record not found")
)

// CertificateRepository defines database operations for certificates and related statistics.
type CertificateRepository interface {
	GetByID(ctx context.Context, certificateID string) (*models.Certificate, error)
	CreateCertificates(ctx context.Context, certs []models.Certificate) error
	UpdateMLStatus(ctx context.Context, certificateID string, status models.MLStatus, mlScore *float64) error
	GetCertificatesPendingFacultyReview(ctx context.Context, limit int) ([]models.Certificate, error)
	UpdateFacultyDecision(ctx context.Context, certificateID string, status models.FacultyStatus, isLegit bool) error
}

type certificateRepository struct {
	db *gorm.DB
}

// NewCertificateRepository constructs a CertificateRepository.
func NewCertificateRepository(db *gorm.DB) CertificateRepository {
	return &certificateRepository{db: db}
}

// GetByID fetches a certificate by ID.
func (r *certificateRepository) GetByID(ctx context.Context, certificateID string) (*models.Certificate, error) {
	var cert models.Certificate
	if err := r.db.WithContext(ctx).First(&cert, "id = ?", certificateID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCertificateNotFound
		}
		return nil, fmt.Errorf("get certificate: %w", err)
	}
	return &cert, nil
}

// CreateCertificates inserts a batch of certificates (max 10) and updates stats in a transaction.
func (r *certificateRepository) CreateCertificates(ctx context.Context, certs []models.Certificate) error {
	if len(certs) == 0 {
		return nil
	}
	if len(certs) > 10 {
		return ErrTooManyCertificates
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&certs).Error; err != nil {
			return fmt.Errorf("insert certificates: %w", err)
		}

		// Update stats for each certificate; assumes rows already exist in stats tables.
		for _, cert := range certs {
			if err := r.incrementStats(ctx, tx, cert); err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateMLStatus sets ml_status (and optional score) and syncs stats counts.
func (r *certificateRepository) UpdateMLStatus(ctx context.Context, certificateID string, status models.MLStatus, mlScore *float64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cert models.Certificate
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&cert, "id = ?", certificateID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCertificateNotFound
			}
			return fmt.Errorf("fetch certificate: %w", err)
		}

		updates := map[string]interface{}{
			"ml_status": status,
		}
		// ml_score column is missing in DB
		// if mlScore != nil {
		// 	updates["ml_score"] = mlScore
		// }

		if err := tx.Model(&cert).Updates(updates).Error; err != nil {
			return fmt.Errorf("update ml status: %w", err)
		}

		// Update stats when ML verifies the certificate.
		if cert.MLStatus != models.MLStatusVerified && status == models.MLStatusVerified {
			if err := r.bumpMlVerified(ctx, tx, cert); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetCertificatesPendingFacultyReview returns ML-verified certificates awaiting faculty decision.
func (r *certificateRepository) GetCertificatesPendingFacultyReview(ctx context.Context, limit int) ([]models.Certificate, error) {
	if limit <= 0 {
		limit = 50
	}
	var certs []models.Certificate
	err := r.db.WithContext(ctx).
		Where("ml_status = ? AND faculty_status = ? AND archived = ?", models.MLStatusVerified, models.FacultyStatusPending, false).
		Order("uploaded_at ASC").
		Limit(limit).
		Find(&certs).Error
	if err != nil {
		return nil, fmt.Errorf("query pending faculty review: %w", err)
	}
	return certs, nil
}

// UpdateFacultyDecision records the faculty decision and updates stats in a transaction.
func (r *certificateRepository) UpdateFacultyDecision(ctx context.Context, certificateID string, status models.FacultyStatus, isLegit bool) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cert models.Certificate
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&cert, "id = ?", certificateID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrCertificateNotFound
			}
			return fmt.Errorf("fetch certificate: %w", err)
		}

		if err := tx.Model(&cert).Updates(map[string]interface{}{
			"faculty_status": status,
			// "is_legit":       isLegit, // Missing in DB
		}).Error; err != nil {
			return fmt.Errorf("update faculty decision: %w", err)
		}

		// Adjust stats only when transitioning from pending.
		if cert.FacultyStatus == models.FacultyStatusPending {
			if err := r.applyFacultyDecisionStats(ctx, tx, cert, status); err != nil {
				return err
			}
		}
		return nil
	})
}

// incrementStats bumps per-student and per-section totals for new certificates.
func (r *certificateRepository) incrementStats(ctx context.Context, tx *gorm.DB, cert models.Certificate) error {
	// Student statistics updates.
	if err := tx.WithContext(ctx).
		Model(&models.StudentStatistics{}).
		Where("reg_no = ?", cert.RegisterNumber).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Updates(map[string]interface{}{
			"total_uploaded": gorm.Expr("total_uploaded + 1"),
			// "pending_certificates": gorm.Expr("pending_certificates + 1"), // Missing
		}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrStatsNotFound
		}
		return fmt.Errorf("update student statistics: %w", err)
	}

	// Section statistics updates.
	if err := tx.WithContext(ctx).
		Model(&models.SectionStatistics{}).
		Where("section = ?", cert.Section).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Updates(map[string]interface{}{
			"total_uploaded": gorm.Expr("total_uploaded + 1"),
		}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrStatsNotFound
		}
		return fmt.Errorf("update section statistics: %w", err)
	}

	return nil
}

func (r *certificateRepository) bumpMlVerified(_ context.Context, _ *gorm.DB, _ models.Certificate) error {
	// Columns pending_certificates and ml_verified_certificates are missing in DB.
	// Skipping update.
	return nil
}

func (r *certificateRepository) applyFacultyDecisionStats(ctx context.Context, tx *gorm.DB, cert models.Certificate, status models.FacultyStatus) error {
	studentUpdates := map[string]interface{}{
		// "pending_certificates": gorm.Expr("pending_certificates - 1"), // Missing
	}
	switch status {
	case models.FacultyStatusLegit:
		studentUpdates["legit_count"] = gorm.Expr("legit_count + 1")
	case models.FacultyStatusNotLegit:
		studentUpdates["not_legit_count"] = gorm.Expr("not_legit_count + 1")
	}

	if len(studentUpdates) > 0 {
		if err := tx.WithContext(ctx).
			Model(&models.StudentStatistics{}).
			Where("reg_no = ?", cert.RegisterNumber).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Updates(studentUpdates).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrStatsNotFound
			}
			return fmt.Errorf("update student statistics: %w", err)
		}
	}

	sectionUpdates := map[string]interface{}{
		// "total_uploaded": gorm.Expr("total_uploaded"), // No change
	}
	switch status {
	case models.FacultyStatusLegit:
		sectionUpdates["legit_count"] = gorm.Expr("legit_count + 1")
	}

	if len(sectionUpdates) > 0 {
		if err := tx.WithContext(ctx).
			Model(&models.SectionStatistics{}).
			Where("section = ?", cert.Section).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Updates(sectionUpdates).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrStatsNotFound
			}
			return fmt.Errorf("update section statistics: %w", err)
		}
	}
	return nil
}
