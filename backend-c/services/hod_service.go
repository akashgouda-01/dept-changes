package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"department-eduvault-backend/internal/excel"
	"department-eduvault-backend/models"
	"department-eduvault-backend/repositories"
)

// StudentStatsDTO is returned to clients for per-student aggregates.
type StudentStatsDTO struct {
	RegisterNumber string `json:"register_number"`
	StudentName    string `json:"student_name"`
	Section        string `json:"section"`
	Total          int64  `json:"total_certificates"`
	Verified       int64  `json:"verified_count"`
	Rejected       int64  `json:"rejected_count"`
	Pending        int64  `json:"pending_count"`
}

// HodService defines HOD-facing operations.
type HodService interface {
	GetStudentStatsByFaculty(ctx context.Context, facultyID string) ([]StudentStatsDTO, error)
	ListStudentCertificates(ctx context.Context, regNo string) ([]models.Certificate, error)
	ExportCertificatesBySection(ctx context.Context, section string) (string, []byte, error)
	ExportCertificatesByStudent(ctx context.Context, regNo string) (string, []byte, error)
}

type hodService struct {
	repo repositories.HodRepository
}

// NewHodService constructs a HodService.
func NewHodService(repo repositories.HodRepository) HodService {
	return &hodService{repo: repo}
}

func (s *hodService) GetStudentStatsByFaculty(ctx context.Context, facultyID string) ([]StudentStatsDTO, error) {
	rows, err := s.repo.GetStudentStatsByFaculty(ctx, strings.TrimSpace(facultyID))
	if err != nil {
		return nil, err
	}

	stats := make([]StudentStatsDTO, 0, len(rows))
	for _, r := range rows {
		stats = append(stats, StudentStatsDTO{
			RegisterNumber: r.RegisterNumber,
			StudentName:    r.StudentName,
			Section:        r.Section,
			Total:          r.Total,
			Verified:       r.Verified,
			Rejected:       r.Rejected,
			Pending:        r.Pending,
		})
	}
	return stats, nil
}

func (s *hodService) ListStudentCertificates(ctx context.Context, regNo string) ([]models.Certificate, error) {
	return s.repo.GetCertificatesByStudent(ctx, strings.TrimSpace(regNo))
}

func (s *hodService) ExportCertificatesBySection(ctx context.Context, section string) (string, []byte, error) {
	section = strings.TrimSpace(section)
	certs, err := s.repo.GetCertificatesBySection(ctx, section)
	if err != nil {
		return "", nil, err
	}
	filename := fmt.Sprintf("certificates_section_%s_%d.xlsx", sanitizeForFilename(section), time.Now().Unix())
	bytes, err := excel.BuildCertificatesWorkbook(certs, fmt.Sprintf("Section-%s", section))
	return filename, bytes, err
}

func (s *hodService) ExportCertificatesByStudent(ctx context.Context, regNo string) (string, []byte, error) {
	regNo = strings.TrimSpace(regNo)
	certs, err := s.repo.GetCertificatesByStudent(ctx, regNo)
	if err != nil {
		return "", nil, err
	}
	filename := fmt.Sprintf("certificates_student_%s_%d.xlsx", sanitizeForFilename(regNo), time.Now().Unix())
	bytes, err := excel.BuildCertificatesWorkbook(certs, fmt.Sprintf("Student-%s", regNo))
	return filename, bytes, err
}

// sanitizeForFilename is a minimal helper to keep filenames readable.
func sanitizeForFilename(val string) string {
	if val == "" {
		return "all"
	}
	return strings.ReplaceAll(val, " ", "_")
}
