package domain

import (
	"time"
)

type Match struct {
	ID             string `db:"id" json:"id"`
	OpponentID     string `db:"opponent_id" json:"opponent_id"`
	MatchDate      time.Time `db:"match_date" json:"match_date"`
	HomeAway       string    `db:"home_away" json:"home_away"`
	OurScore       int16     `db:"our_score" json:"our_score"`
	OpponentScore  int16     `db:"opponent_score" json:"opponent_score"`
	Season         string    `db:"season" json:"season"`
	Status         string    `db:"status" json:"status"`
	IsDerby        bool      `db:"is_derby" json:"is_derby"`
}

