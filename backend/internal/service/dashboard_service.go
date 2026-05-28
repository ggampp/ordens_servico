package service

import (
	"context"

	"github.com/ggampp/ordens_servico/backend/internal/model"
	"github.com/ggampp/ordens_servico/backend/internal/repository"
)

// DashboardService exposes aggregate indicators.
type DashboardService struct {
	repo *repository.DashboardRepository
}

// NewDashboardService builds a DashboardService.
func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

// Summary returns the dashboard indicators.
func (s *DashboardService) Summary(ctx context.Context) (*model.DashboardSummary, error) {
	return s.repo.Summary(ctx)
}
