package dto

import "time"

// MatchResponse — ответ с данными матча
type MatchResponse struct {
    ID            string    `json:"id" db:"id"`
    Opponent      string    `json:"opponent" db:"opponent"`
    LogoOpponent  string    `json:"logo_opponent" db:"logo_opponent"`
    MatchDate     time.Time `json:"match_date" db:"match_date"`        // ← db:"match_date"
    HomeAway      string    `json:"home_away" db:"home_away"`          // ← db:"home_away"
    OurScore      int16     `json:"our_score" db:"our_score"`          // ← db:"our_score"
    OpponentScore int16     `json:"opponent_score" db:"opponent_score"` // ← db:"opponent_score"
    Status        string    `json:"status" db:"status"`                // ← db:"status"
    IsDerby       bool      `json:"is_derby" db:"is_derby"`            // ← db:"is_derby"
    WinType       string    `json:"win_type,omitempty" db:"win_type"`  // ← db:"win_type"
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


//статистика по прошедшим матчам
type MatchStatsResponse struct {
    Wins struct {
        Regular  int `json:"regular"`  // победы в основное время
        Overtime int `json:"overtime"` // победы в овертайме
        Penalty  int `json:"penalty"`  // победы по буллитам
    } `json:"wins"`
    Losses struct {
        Regular  int `json:"regular"`  // поражения в основное время
        Overtime int `json:"overtime"` // поражения в овертайме
        Penalty  int `json:"penalty"`  // поражения по буллитам
    } `json:"losses"`
    Total int `json:"total"`
}
