package sector

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/TanyaMasharova/hockey-team/internal/repository/interfaces"
)

type Service struct {
	sectorRepo interfaces.SectorRepository
}

func NewService(sectorRepo interfaces.SectorRepository) *Service {
	if sectorRepo == nil {
		panic("sectorRepo is required")
	}
	return &Service{
		sectorRepo: sectorRepo,
	}
}

// GetAllSectors возвращает все секторы стадиона
func (s *Service) GetAllSectors(ctx context.Context) ([]*domain.StadiumSector, error) {
	sectors, err := s.sectorRepo.GetAllSectors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all sectors: %w", err)
	}
	return sectors, nil
}

// GetSectorByID возвращает сектор по ID
func (s *Service) GetSectorByID(ctx context.Context, sectorID string) (*domain.StadiumSector, error) {
	if sectorID == "" {
		return nil, fmt.Errorf("sector ID is required")
	}

	sector, err := s.sectorRepo.GetSectorByID(ctx, sectorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sector: %w", err)
	}

	return sector, nil
}
