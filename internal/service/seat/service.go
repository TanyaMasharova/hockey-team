package seat

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/TanyaMasharova/hockey-team/internal/repository/interfaces"
)

type Service struct {
    seatRepo interfaces.SeatRepository
}

func NewService(seatRepo interfaces.SeatRepository) *Service {
    if seatRepo == nil {
        panic("seatRepo is required")
    }
    return &Service{seatRepo: seatRepo}
}

func (s *Service) GetSeatsBySector(ctx context.Context, sectorID string, matchID string) ([]*domain.Seat, error) {
    if sectorID == "" {
        return nil, fmt.Errorf("sector ID is required")
    }
    if matchID == "" {
        return nil, fmt.Errorf("match ID is required")
    }
    
    seats, err := s.seatRepo.GetSeatsBySector(ctx, sectorID, matchID)
    if err != nil {
        return nil, fmt.Errorf("failed to get seats: %w", err)
    }
    
    return seats, nil
}