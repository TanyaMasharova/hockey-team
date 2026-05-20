// internal/repository/postgres/admin_repo.go
package postgres

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type adminRepo struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) GetSalesByMonth(ctx context.Context) ([]domain.SalesByMonth, error) {
	query := `
        SELECT 
            TO_CHAR(m.match_date, 'YYYY-MM') as month,
            COUNT(t.id) as tickets,
            COALESCE(SUM(t.final_price), 0) as revenue
        FROM tickets t
        JOIN matches m ON t.match_id = m.id
        WHERE t.status != 'cancelled'
        GROUP BY TO_CHAR(m.match_date, 'YYYY-MM')
        ORDER BY month DESC
        LIMIT 12
    `

	var result []domain.SalesByMonth
	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales by month: %w", err)
	}

	// Добавьте лог для проверки
	logrus.Infof("Sales by month result: %+v", result)

	return result, nil
}
func (r *adminRepo) GetSectorPopularity(ctx context.Context) ([]domain.SectorPopularity, error) {
	query := `
        SELECT 
            ss.sector_number as sector,
            COUNT(t.id) as sold,
            ss.capacity,
            ROUND(COUNT(t.id) * 100.0 / ss.capacity, 2) as occupancy_percent
        FROM stadium_sectors ss
        LEFT JOIN seats s ON s.sector_id = ss.id
        LEFT JOIN tickets t ON t.seat_id = s.id AND t.status != 'cancelled'
        GROUP BY ss.sector_number, ss.capacity
        ORDER BY sold DESC
    `

	var result []domain.SectorPopularity
	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sector popularity: %w", err)
	}
	return result, nil
}

func (r *adminRepo) GetTicketStatus(ctx context.Context) ([]domain.TicketStatus, error) {
	query := `
        SELECT 
            status,
            COUNT(*) as count
        FROM tickets
        GROUP BY status
    `

	var result []domain.TicketStatus
	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket status: %w", err)
	}
	return result, nil
}

func (r *adminRepo) GetAvgPriceBySector(ctx context.Context) ([]domain.AvgPriceBySector, error) {
	query := `
        SELECT 
            ss.sector_number as sector,
            ss.sector_type,
            ROUND(AVG(t.final_price), 0) as avg_price,
            COUNT(t.id) as tickets_sold
        FROM tickets t
        JOIN seats s ON t.seat_id = s.id
        JOIN stadium_sectors ss ON s.sector_id = ss.id
        WHERE t.status != 'cancelled'
        GROUP BY ss.sector_number, ss.sector_type
        ORDER BY avg_price DESC
    `

	var result []domain.AvgPriceBySector
	err := r.db.SelectContext(ctx, &result, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get avg price by sector: %w", err)
	}
	return result, nil
}

func (r *adminRepo) GetTopBuyers(ctx context.Context, limit int) ([]domain.TopBuyer, error) {
	query := `
        SELECT 
            u.full_name,
            u.email,
            COUNT(t.id) as tickets_count,
            COALESCE(SUM(t.final_price), 0) as total_spent
        FROM users u
        JOIN tickets t ON t.user_id = u.id
        WHERE t.status != 'cancelled'
        GROUP BY u.id, u.full_name, u.email
        ORDER BY tickets_count DESC
        LIMIT $1
    `

	var result []domain.TopBuyer
	err := r.db.SelectContext(ctx, &result, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top buyers: %w", err)
	}
	return result, nil
}

func (r *adminRepo) GetStatsSummary(ctx context.Context) (*domain.AdminStats, error) {
	query := `
        SELECT 
            (SELECT COUNT(*) FROM users) as total_users,
            (SELECT COUNT(*) FROM tickets WHERE status != 'cancelled') as total_tickets,
            (SELECT COALESCE(SUM(final_price), 0) FROM tickets WHERE status != 'cancelled') as total_revenue,
            (SELECT COUNT(*) FROM tickets WHERE status = 'active') as active_tickets,
            (SELECT COUNT(*) FROM tickets WHERE status = 'used') as used_tickets
    `

	var result domain.AdminStats
	err := r.db.GetContext(ctx, &result, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats summary: %w", err)
	}

	logrus.Infof("Stats summary: total_users=%d, total_tickets=%d, total_revenue=%f, active=%d, used=%d",
		result.TotalUsers, result.TotalTickets, result.TotalRevenue, result.ActiveTickets, result.UsedTickets)

	return &result, nil
}
