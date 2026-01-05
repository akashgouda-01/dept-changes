package services

import (
	"context"
	"errors"
	"regexp"
	"time"

	"department-eduvault-backend/models"
	"department-eduvault-backend/repositories"
)

var (
	ErrInvalidDriveLink            = errors.New("drive link must be a valid Google Drive URL")
	ErrUploadLimitExceeded         = errors.New("cannot upload more than 10 certificates in one request")
	ErrInvalidMLTransition         = errors.New("ml status transition is not allowed")
	ErrInvalidFacultyState         = errors.New("faculty decision is only allowed after ML verification and while pending")
	ErrCertificateArchived         = errors.New("archived certificates cannot be modified")
	driveLinkPattern               = regexp.MustCompile(`^https://drive\.google\.com/`)
	defaultMLScore         float64 = 95.0
)

// CertificateInput represents an upload payload.
type CertificateInput struct {
	DriveLink      string
	RegisterNumber string
	Section        string
	StudentName    string
	UploadedBy     string
	UploadedAt     time.Time
}

// CertificateService describes business operations for certificates.
type CertificateService interface {
	UploadCertificates(ctx context.Context, inputs []CertificateInput) error
	TriggerMockMLVerification(ctx context.Context, certificateID string) error
	GetPendingFacultyReview(ctx context.Context, limit int) ([]models.Certificate, error)
	SubmitFacultyDecision(ctx context.Context, certificateID string, status models.FacultyStatus, isLegit bool) error
}

type certificateService struct {
	repo repositories.CertificateRepository
}

// NewCertificateService constructs a CertificateService.
func NewCertificateService(repo repositories.CertificateRepository) CertificateService {
	return &certificateService{repo: repo}
}

// UploadCertificates validates input and delegates creation; kicks off mock ML verification.
func (s *certificateService) UploadCertificates(ctx context.Context, inputs []CertificateInput) error {
	if len(inputs) == 0 {
		return nil
	}
	if len(inputs) > 10 {
		return ErrUploadLimitExceeded
	}

	certs := make([]models.Certificate, 0, len(inputs))
	for _, in := range inputs {
		if !driveLinkPattern.MatchString(in.DriveLink) {
			return ErrInvalidDriveLink
		}
		uploadedAt := in.UploadedAt
		if uploadedAt.IsZero() {
			uploadedAt = time.Now().UTC()
		}
		certs = append(certs, models.Certificate{
			DriveLink:      in.DriveLink,
			RegisterNumber: in.RegisterNumber,
			Section:        in.Section,
			StudentName:    in.StudentName,
			UploadedBy:     in.UploadedBy,
			UploadedAt:     uploadedAt,
			MLStatus:       models.MLStatusPending,
			FacultyStatus:  models.FacultyStatusPending,
			Archived:       false,
		})
	}

	if err := s.repo.CreateCertificates(ctx, certs); err != nil {
		return err
	}

	// Trigger mock async ML verification.
	for _, cert := range certs {
		certID := cert.ID
		go func(id string) {
			_ = s.TriggerMockMLVerification(context.Background(), id)
		}(certID)
	}
	return nil
}

// TriggerMockMLVerification simulates an asynchronous ML process by marking verified.
func (s *certificateService) TriggerMockMLVerification(ctx context.Context, certificateID string) error {
	cert, err := s.repo.GetByID(ctx, certificateID)
	if err != nil {
		return err
	}
	if cert.Archived {
		return ErrCertificateArchived
	}
	if cert.MLStatus != models.MLStatusPending {
		return ErrInvalidMLTransition
	}

	return s.repo.UpdateMLStatus(ctx, certificateID, models.MLStatusVerified, &defaultMLScore)
}

// GetPendingFacultyReview fetches ML-verified certificates pending faculty action.
func (s *certificateService) GetPendingFacultyReview(ctx context.Context, limit int) ([]models.Certificate, error) {
	return s.repo.GetCertificatesPendingFacultyReview(ctx, limit)
}

// SubmitFacultyDecision records a faculty decision with state validation.
func (s *certificateService) SubmitFacultyDecision(ctx context.Context, certificateID string, status models.FacultyStatus, isLegit bool) error {
	if status != models.FacultyStatusLegit && status != models.FacultyStatusNotLegit {
		return ErrInvalidFacultyState
	}

	cert, err := s.repo.GetByID(ctx, certificateID)
	if err != nil {
		return err
	}
	if cert.Archived {
		return ErrCertificateArchived
	}
	if cert.MLStatus != models.MLStatusVerified || cert.FacultyStatus != models.FacultyStatusPending {
		return ErrInvalidFacultyState
	}

	return s.repo.UpdateFacultyDecision(ctx, certificateID, status, isLegit)
}

// Helpers (kept unexported) ---------------------------------------------------

func (s *certificateService) validateDriveLink(link string) error {
	if !driveLinkPattern.MatchString(link) {
		return ErrInvalidDriveLink
	}
	return nil
}

// For future extension where per-link validation might be more complex.
func (s *certificateService) validateInputs(inputs []CertificateInput) error {
	if len(inputs) > 10 {
		return ErrUploadLimitExceeded
	}
	for _, in := range inputs {
		if err := s.validateDriveLink(in.DriveLink); err != nil {
			return err
		}
	}
	return nil
}
