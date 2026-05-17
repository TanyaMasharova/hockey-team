package postgres

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/jmoiron/sqlx"
)

type sectorRepo struct {
    db *sqlx.DB
}

func NewSectorRepository(db *sqlx.DB) *sectorRepo {
    return &sectorRepo{db: db}
}

func (r *sectorRepo) GetAllSectors(ctx context.Context) ([]*domain.StadiumSector, error) {
    var sectors []*domain.StadiumSector
    
    query := `
        SELECT 
            id,
            sector_number,
            capacity,
            sector_type,
            price_coefficient,
            color_code
        FROM stadium_sectors
        ORDER BY sector_number
    `
    
    err := r.db.SelectContext(ctx, &sectors, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get all sectors: %w", err)
    }
    
    return sectors, nil
}

func (r *sectorRepo) GetSectorByID(ctx context.Context, sectorID string) (*domain.StadiumSector, error) {
    var sector domain.StadiumSector
    
    query := `
        SELECT 
            id,
            sector_number,
            capacity,
            sector_type,
            price_coefficient,
            color_code
        FROM stadium_sectors
        WHERE id = $1
    `
    
    err := r.db.GetContext(ctx, &sector, query, sectorID)
    if err != nil {
        return nil, fmt.Errorf("failed to get sector by id: %w", err)
    }
    
    return &sector, nil
}