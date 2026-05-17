package interfaces

import (
	"context"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/dto"
	"github.com/TanyaMasharova/hockey-team/internal/domain"
)

type MatchRepository interface {
	GetMatches(ctx context.Context, limit *int, futurePast *string) ([]dto.MatchResponse, error)
	// GetFutureMatches(ctx context.Context) ([]domain.Match, error)
	GetMatchByID(ctx context.Context, id string) (*dto.MatchResponse, error)
	GetMatchesBySeason(ctx context.Context, season string) ([]dto.MatchResponse, error)
	GetMatchesStats(ctx context.Context) (*dto.MatchStatsResponse, error)

	 GetMatchWithOpponent(ctx context.Context, matchID string) (*domain.MatchWithOpponent, error)
}