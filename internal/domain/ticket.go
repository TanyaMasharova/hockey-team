package domain

import (
	"time"
)

type TicketWithDetails struct {
	ID           string    `db:"ticket_id" json:"id"`
	FinalPrice   float64   `db:"final_price" json:"final_price"`
	PurchaseDate time.Time `db:"purchase_date" json:"purchase_date"`
	Status       string    `db:"ticket_status" json:"status"`
	
	// Информация о матче
	MatchDate    time.Time `db:"match_date" json:"match_date"`
	HomeAway     string    `db:"home_away" json:"home_away"`
	OpponentName string    `db:"opponent_name" json:"opponent_name"`
	OpponentCity string    `db:"opponent_city" json:"opponent_city,omitempty"`
	
	// Информация о месте
	SectorNumber string `db:"sector_number" json:"sector_number"`
	SectorType   string `db:"sector_type" json:"sector_type"`
	SeatRow      string `db:"seat_row" json:"seat_row"`
	SeatNumber   string `db:"seat_number" json:"seat_number"`
}