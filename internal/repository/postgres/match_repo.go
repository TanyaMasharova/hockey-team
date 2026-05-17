package postgres

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/dto"
	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type matchRepo struct {
    db *sqlx.DB
}

func NewMatchRepository(db *sqlx.DB) *matchRepo {
	return &matchRepo{db: db}
}

func (r *matchRepo) GetMatches(ctx context.Context, limit *int, futurePast *string) ([]dto.MatchResponse, error) {
	query := `
		SELECT m.id, 
            o.name as opponent,
						COALESCE(o.logo_url, '') as logo_opponent,
            m.match_date, 
            m.home_away, 
            m.our_score, 
            m.opponent_score, 
            m.status, 
            m.is_derby, 
						m.win_type
		FROM matches m
		JOIN opponents o ON m.opponent_id = o.id
		WHERE 1=1
	`
	

	if futurePast != nil && *futurePast != "" {
		if *futurePast == "future" {
			query += " AND DATE(m.match_date) >= CURRENT_DATE"
			
		} else if *futurePast == "past" {
			query += " AND DATE(m.match_date) < CURRENT_DATE"
		}
	}
logrus.Info(query)
	if limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *limit)
	}

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
			&match.WinType,
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




func (r *matchRepo) GetMatchByID(ctx context.Context, id string) (*dto.MatchResponse, error) {
	query := `
		SELECT m.id, 
			o.name as opponent,
			COALESCE(o.logo_url, '') as logo_opponent,
			m.match_date, 
			m.home_away, 
			m.our_score, 
			m.opponent_score, 
			m.status, 
			m.is_derby, 
			COALESCE(m.win_type, '') as win_type
		FROM matches m
		JOIN opponents o ON m.opponent_id = o.id
		WHERE m.id = $1
	`

	var match dto.MatchResponse
	err := r.db.GetContext(ctx, &match, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get match by id: %w", err)
	}
	
	return &match, nil
}

func 	(r *matchRepo) GetMatchesBySeason(ctx context.Context, season string) ([]dto.MatchResponse, error) {
	return []dto.MatchResponse{}, nil
}

func (r *matchRepo) GetMatchesStats(ctx context.Context) (*dto.MatchStatsResponse, error) {
	query := `
		SELECT 
    -- Победы
    COUNT(CASE WHEN our_score > opponent_score AND win_type = 'regular' THEN 1 END) as wins_regular,
    COUNT(CASE WHEN our_score > opponent_score AND win_type = 'overtime' THEN 1 END) as wins_overtime,
    COUNT(CASE WHEN our_score > opponent_score AND win_type = 'penalty' THEN 1 END) as wins_penalty,
    -- Поражения
    COUNT(CASE WHEN our_score < opponent_score AND win_type = 'regular' THEN 1 END) as losses_regular,
    COUNT(CASE WHEN our_score < opponent_score AND win_type = 'overtime' THEN 1 END) as losses_overtime,
    COUNT(CASE WHEN our_score < opponent_score AND win_type = 'penalty' THEN 1 END) as losses_penalty
		FROM matches 
		WHERE status = 'finished';
	`
	var stats dto.MatchStatsResponse

	err := r.db.QueryRowContext(ctx, query).Scan(
		 &stats.Wins.Regular,
        &stats.Wins.Overtime,
        &stats.Wins.Penalty,
        &stats.Losses.Regular,
        &stats.Losses.Overtime,
        &stats.Losses.Penalty,
	)
	if err != nil {
        return nil, fmt.Errorf("failed to get match stats: %w", err)
    }
		stats.Total = stats.Wins.Regular + stats.Wins.Overtime + stats.Wins.Penalty + 
                  stats.Losses.Regular + stats.Losses.Overtime + stats.Losses.Penalty
		return &stats, nil

}
func (r *matchRepo) GetMatchWithOpponent(ctx context.Context, matchID string) (*domain.MatchWithOpponent, error) {
    var match domain.MatchWithOpponent
    
    query := `
        SELECT 
            m.id as match_id,
            o.name as opponent_name,
            o.logo as opponent_logo,
            m.match_date,
            m.home_away,
            m.our_score,
            m.opponent_score,
            m.status,
            m.is_derby,
            COALESCE(m.win_type, '') as win_type
        FROM matches m
        JOIN opponents o ON m.opponent_id = o.id
        WHERE m.id = $1 AND m.deleted_at IS NULL
    `
    
    err := r.db.GetContext(ctx, &match, query, matchID)
    if err != nil {
        return nil, fmt.Errorf("failed to get match with opponent: %w", err)
    }
    
    return &match, nil
}
