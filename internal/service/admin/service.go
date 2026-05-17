// internal/service/admin/service.go
package admin

import (
	"context"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/TanyaMasharova/hockey-team/internal/repository/interfaces"
)

type Service struct {
    adminRepo interfaces.AdminRepository
}

func NewService(adminRepo interfaces.AdminRepository) *Service {
    return &Service{
        adminRepo: adminRepo,
    }
}

func (s *Service) GetSalesByMonth(ctx context.Context) ([]domain.SalesByMonth, error) {
    return s.adminRepo.GetSalesByMonth(ctx)
}

func (s *Service) GetSectorPopularity(ctx context.Context) ([]domain.SectorPopularity, error) {
    return s.adminRepo.GetSectorPopularity(ctx)
}

func (s *Service) GetTicketStatus(ctx context.Context) ([]domain.TicketStatus, error) {
    return s.adminRepo.GetTicketStatus(ctx)
}

func (s *Service) GetAvgPriceBySector(ctx context.Context) ([]domain.AvgPriceBySector, error) {
    return s.adminRepo.GetAvgPriceBySector(ctx)
}

func (s *Service) GetTopBuyers(ctx context.Context, limit int) ([]domain.TopBuyer, error) {
    return s.adminRepo.GetTopBuyers(ctx, limit)
}

func (s *Service) GetStatsSummary(ctx context.Context) (*domain.AdminStats, error) {
    return s.adminRepo.GetStatsSummary(ctx)
}