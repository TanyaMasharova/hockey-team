package interfaces

import (
	"context"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/dto"
	"github.com/TanyaMasharova/hockey-team/internal/domain"
)

type MatchRepository interface {
	GetMatches(ctx context.Context) ([]dto.MatchResponse, error)
	// GetFutureMatches(ctx context.Context) ([]domain.Match, error)
	GetMatchByID(ctx context.Context, id string) (*domain.Match, error)
	GetMatchesBySeason(ctx context.Context, season string) ([]domain.Match, error)
}