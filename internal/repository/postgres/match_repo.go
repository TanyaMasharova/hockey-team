package postgres

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/dto"
	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/jmoiron/sqlx"
)

type matchRepo struct {
    db *sqlx.DB
}

func NewMatchRepository(db *sqlx.DB) *matchRepo {
	return &matchRepo{db: db}
}

func (r *matchRepo) GetMatches(ctx context.Context) ([]dto.MatchResponse, error) {
	query := `
		SELECT m.id, 
            o.name as opponent,
						o.logo_url as logo_opponent,
            m.match_date, 
            m.home_away, 
            m.our_score, 
            m.opponent_score, 
            m.status, 
            m.is_derby
		FROM matches m
		JOIN opponents o ON m.opponent_id = o.id
		ORDER BY m.match_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query matches: %w", err)
	}
	defer rows.Close() //гарантирует, что после выполнения всей функции GetMatches ресурсы будут освобождены
	var matches []dto.MatchResponse
	for rows.Next() { //цикл для каждой строки (для БД!)
		var match dto.MatchResponse
		err := rows.Scan(
			&match.ID,
			&match.Opponent,
			&match.LogoOpponent,
			&match.MatchDate,
			&match.HomeAway,
			&match.OurScore,
			&match.OpponentScore,
			// &match.Season,
			&match.Status,
			&match.IsDerby,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan match: %w", err)
		}
		matches = append(matches, match)
	}
	if err = rows.Err(); err != nil { //для отслеживания ошибок во время итерации
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return  matches, nil

}




func (r *matchRepo) GetMatchByID(ctx context.Context, id string) (*domain.Match, error) {
	return &domain.Match{}, nil
}

func 	(r *matchRepo) GetMatchesBySeason(ctx context.Context, season string) ([]domain.Match, error) {
	return []domain.Match{}, nil
}