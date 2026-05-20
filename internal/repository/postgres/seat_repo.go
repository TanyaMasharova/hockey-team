package postgres

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/jmoiron/sqlx"
)

type seatRepo struct {
	db *sqlx.DB
}

func NewSeatRepository(db *sqlx.DB) *seatRepo {
	return &seatRepo{db: db}
}

// GetSeatsBySector возвращает места с учётом занятости на конкретный матч
func (r *seatRepo) GetSeatsBySector(ctx context.Context, sectorID string, matchID string) ([]*domain.Seat, error) {
	var seats []*domain.Seat

	query := `
        SELECT 
            s.id,
            s.sector_id,
            s.seat_row,
            s.seat_number,
            s.is_handicap_accessible,
            CASE 
                WHEN t.id IS NOT NULL AND t.status IN ('active', 'used') THEN true 
                ELSE false 
            END as is_taken
        FROM seats s
        LEFT JOIN tickets t ON t.seat_id = s.id AND t.match_id = $2
        WHERE s.sector_id = $1
        ORDER BY s.seat_row, s.seat_number::int
    `

	err := r.db.SelectContext(ctx, &seats, query, sectorID, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seats by sector: %w", err)
	}

	return seats, nil
}

func (r *seatRepo) GetSeatByID(ctx context.Context, seatID string) (*domain.Seat, error) {
	var seat domain.Seat

	query := `
        SELECT 
            id,
            sector_id,
            seat_row,
            seat_number,
            is_handicap_accessible,
            is_taken
        FROM seats
        WHERE id = $1
    `

	err := r.db.GetContext(ctx, &seat, query, seatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get seat by id: %w", err)
	}

	return &seat, nil
}

// UpdateSeatStatus - этот метод лучше не использовать для билетов,
// так как занятость зависит от матча, а не от самого места
func (r *seatRepo) UpdateSeatStatus(ctx context.Context, seatID string, isTaken bool) error {
	query := `
        UPDATE seats 
        SET is_taken = $1, updated_at = NOW()
        WHERE id = $2
    `

	_, err := r.db.ExecContext(ctx, query, isTaken, seatID)
	if err != nil {
		return fmt.Errorf("failed to update seat status: %w", err)
	}

	return nil
}
