package matches

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/dto"
	"github.com/TanyaMasharova/hockey-team/internal/repository/interfaces"
)

type Service struct {
	matchRepo interfaces.MatchRepository
}

func NewService(matchRepo interfaces.MatchRepository) *Service {
	if matchRepo == nil {
		panic("userRepo is required") // защита от ошибок
	}

	return &Service{
		matchRepo: matchRepo,
	}
}

func (s *Service) GetMatches(ctx context.Context, limit *int, futurePast *string) ([]dto.MatchResponse, error) {

	//создать переменную куда загружать ответ от репозитория
	var matches []dto.MatchResponse

	var limitValue *int

	// Проверяем, передан ли параметр limit
	if limit != nil && *limit > 0 {
		limitValue = limit
	}

	var futurePastValue *string
	if futurePast != nil && *futurePast != "" {
		futurePastValue = futurePast
	}

	//вызвать метод интерфейса репозитория
	result, err := s.matchRepo.GetMatches(ctx, limitValue, futurePastValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches: %w", err)
	}
	matches = result //преобразование вот тут будет
	return matches, nil
	//если ошибка - прерывание и вывод ошибки
}

func (s *Service) GetStatsMatches(ctx context.Context) (*dto.MatchStatsResponse, error) {

	stats, err := s.matchRepo.GetMatchesStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches stats: %w", err)
	}

	return stats, nil

}

// GetMatchByID возвращает матч по ID
func (s *Service) GetMatchByID(ctx context.Context, matchID string) (*dto.MatchResponse, error) {
	if matchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}

	match, err := s.matchRepo.GetMatchByID(ctx, matchID)
	if err != nil {
		return nil, fmt.Errorf("failed to get match: %w", err)
	}

	return match, nil
}
