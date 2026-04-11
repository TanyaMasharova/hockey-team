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

func (s *Service) GetMatches(ctx context.Context) ([]dto.MatchResponse, error) {
		
	//создать переменную куда загружать ответ от репозитория
	var matches []dto.MatchResponse

	//вызвать метод интерфейса репозитория
	result, err := s.matchRepo.GetMatches(ctx)
	if err != nil {
				return nil, fmt.Errorf("failed to get matches: %w", err)
	}
	matches = result //преобразование вот тут будет
	return  matches, nil
	//если ошибка - прерывание и вывод ошибки
}