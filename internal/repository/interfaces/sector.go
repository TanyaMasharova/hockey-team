package interfaces

import (
	"context"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
)

type SectorRepository interface {
	GetAllSectors(ctx context.Context) ([]*domain.StadiumSector, error)
	GetSectorByID(ctx context.Context, sectorID string) (*domain.StadiumSector, error)
}
