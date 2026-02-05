package services

import (
	"context"
	"department-eduvault-backend/internal/excel"
	"department-eduvault-backend/repositories"
	"fmt"
	"strings"
	"time"
)

type DashboardOverviewDTO struct {
	TotalStudents     int64 `json:"total_students"`
	TotalCertificates int64 `json:"total_certificates"`
	VerifiedCount     int64 `json:"verified_count"`
	RejectedCount     int64 `json:"rejected_count"`
	PendingCount      int64 `json:"pending_count"`
}

type SectionStatsDTO struct {
	Section           string  `json:"section"`
	TotalCertificates int64   `json:"total_certificates"`
	VerifiedCount     int64   `json:"verified_count"`
	RejectedCount     int64   `json:"rejected_count"`
	PendingCount      int64   `json:"pending_count"`
	VerificationRate  float64 `json:"verification_rate"`
}

type RecentActivityDTO struct {
	ID          string `json:"id"`
	StudentName string `json:"student_name"`
	RegNo       string `json:"reg_no"`
	Section     string `json:"section"`
	Action      string `json:"action"`
	Timestamp   string `json:"timestamp"`
}

type DashboardService interface {
	GetOverview(ctx context.Context, role, facultyID string) (DashboardOverviewDTO, error)
	GetSectionStats(ctx context.Context, role, facultyID string) ([]SectionStatsDTO, error)
	GetRecentActivity(ctx context.Context, role, facultyID string) ([]RecentActivityDTO, error)
	ExportMyCertificates(ctx context.Context, facultyID string) (string, []byte, error)
}

type dashboardService struct {
	repo repositories.DashboardRepository
}

func NewDashboardService(repo repositories.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetOverview(ctx context.Context, role, facultyID string) (DashboardOverviewDTO, error) {
	// If Role is HOD, we fetch all.
	// If Role is Faculty, we fetch only for their facultyID.
	var ov repositories.DashboardOverview
	var err error

	if role == "hod" {
		ov, err = s.repo.GetOverview(ctx, "") // empty means all
	} else {
		ov, err = s.repo.GetOverview(ctx, facultyID)
	}

	if err != nil {
		return DashboardOverviewDTO{}, err
	}
	return DashboardOverviewDTO{
		TotalStudents:     ov.TotalStudents,
		TotalCertificates: ov.TotalCertificates,
		VerifiedCount:     ov.VerifiedCertificates,
		RejectedCount:     ov.RejectedCertificates,
		PendingCount:      ov.PendingCertificates,
	}, nil
}

func (s *dashboardService) GetSectionStats(ctx context.Context, role, facultyID string) ([]SectionStatsDTO, error) {
	var rows []repositories.SectionDashboardRow
	var err error

	if role == "hod" {
		rows, err = s.repo.GetSectionStats(ctx, "")
	} else {
		rows, err = s.repo.GetSectionStats(ctx, facultyID)
	}

	if err != nil {
		return nil, err
	}
	result := make([]SectionStatsDTO, 0, len(rows))
	for _, r := range rows {
		total := r.TotalCertificates
		var rate float64
		if total > 0 {
			rate = float64(r.VerifiedCertificates) / float64(total)
		}
		result = append(result, SectionStatsDTO{
			Section:           r.Section,
			TotalCertificates: total,
			VerifiedCount:     r.VerifiedCertificates,
			RejectedCount:     r.RejectedCertificates,
			PendingCount:      r.PendingCertificates,
			VerificationRate:  rate,
		})
	}
	return result, nil
}

func (s *dashboardService) GetRecentActivity(ctx context.Context, role, facultyID string) ([]RecentActivityDTO, error) {
	var rows []repositories.RecentActivityRow
	var err error

	if role == "hod" {
		rows, err = s.repo.GetRecentActivity(ctx, "")
	} else {
		rows, err = s.repo.GetRecentActivity(ctx, facultyID)
	}

	if err != nil {
		return nil, err
	}

	result := make([]RecentActivityDTO, 0, len(rows))
	for _, r := range rows {
		result = append(result, RecentActivityDTO{
			ID:          r.ID,
			StudentName: r.StudentName,
			RegNo:       r.RegNo,
			Section:     r.Section,
			Action:      r.Action,
			Timestamp:   r.Timestamp,
		})
	}
	return result, nil
}

func (s *dashboardService) ExportMyCertificates(ctx context.Context, facultyID string) (string, []byte, error) {
	certs, err := s.repo.GetCertificatesByFacultyID(ctx, facultyID)
	if err != nil {
		return "", nil, err
	}

	dateStr := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s_VERIFIED_CERTIFICATES_%s.xlsx", strings.TrimSpace(facultyID), dateStr)

	bytes, err := excel.BuildCertificatesWorkbook(certs, fmt.Sprintf("Faculty-%s", facultyID))
	return filename, bytes, err
}
