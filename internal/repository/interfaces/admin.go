// internal/repository/interfaces/admin.go
package interfaces

import (
	"context"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
)

type AdminRepository interface {
    GetSalesByMonth(ctx context.Context) ([]domain.SalesByMonth, error)
    GetSectorPopularity(ctx context.Context) ([]domain.SectorPopularity, error)
    GetTicketStatus(ctx context.Context) ([]domain.TicketStatus, error)
    GetAvgPriceBySector(ctx context.Context) ([]domain.AvgPriceBySector, error)
    GetTopBuyers(ctx context.Context, limit int) ([]domain.TopBuyer, error)
    GetStatsSummary(ctx context.Context) (*domain.AdminStats, error)
}