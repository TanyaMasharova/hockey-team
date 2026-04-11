package dto

import "time"

// MatchResponse — ответ с данными матча
type MatchResponse struct {
    ID             string    `json:"id"`
    Opponent     string    `json:"opponent"`
		LogoOpponent	string	`json:"logo_opponent"`
    MatchDate      time.Time `json:"match_date"`
    HomeAway       string    `json:"home_away"`
    OurScore       int16     `json:"our_score"`
    OpponentScore  int16     `json:"opponent_score"`
    // Season         string    `json:"season"`
    Status         string    `json:"status"`
    IsDerby        bool      `json:"is_derby"`
}

// ListMatchesResponse — ответ со списком матчей
type ListMatchesResponse struct {
    Matches []MatchResponse `json:"matches"`
    Total   int             `json:"total"`
}

// фильтр по сезону 
type GetMatchesBySeasonRequest struct {
    Season string `query:"season" validate:"required"`
}

// фильтр по статусу
type GetMatchesByStatusRequest struct {
    Status string `query:"status" validate:"required,oneof=scheduled live finished cancelled"`
}