package interfaces

import (
	"context"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
)

type SeatRepository interface {
	GetSeatsBySector(ctx context.Context, sectorID string, matchID string) ([]*domain.Seat, error)
	GetSeatByID(ctx context.Context, seatID string) (*domain.Seat, error)
	UpdateSeatStatus(ctx context.Context, seatID string, isTaken bool) error
}
