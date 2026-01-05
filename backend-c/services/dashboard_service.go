package services

import (
	"context"

	"department-eduvault-backend/repositories"
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

type DashboardService interface {
	GetOverview(ctx context.Context) (DashboardOverviewDTO, error)
	GetSectionStats(ctx context.Context) ([]SectionStatsDTO, error)
}

type dashboardService struct {
	repo repositories.DashboardRepository
}

func NewDashboardService(repo repositories.DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetOverview(ctx context.Context) (DashboardOverviewDTO, error) {
	ov, err := s.repo.GetOverview(ctx)
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

func (s *dashboardService) GetSectionStats(ctx context.Context) ([]SectionStatsDTO, error) {
	rows, err := s.repo.GetSectionStats(ctx)
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
