package domain

import "time"

// Match структура матча
type Match struct {
	ID            string    `db:"id"`
	OpponentID    string    `db:"opponent_id"`
	MatchDate     time.Time `db:"match_date"`
	HomeAway      string    `db:"home_away"`
	OurScore      int16     `db:"our_score"`
	OpponentScore int16     `db:"opponent_score"`
	Status        string    `db:"status"`
	IsDerby       bool      `db:"is_derby"`
	WinType       string    `db:"win_type"`
}

// MatchWithOpponent матч с информацией о сопернике
type MatchWithOpponent struct {
	ID            string    `db:"match_id"`
	Opponent      string    `db:"opponent_name"`
	LogoOpponent  string    `db:"opponent_logo"`
	MatchDate     time.Time `db:"match_date"`
	HomeAway      string    `db:"home_away"`
	OurScore      int16     `db:"our_score"`
	OpponentScore int16     `db:"opponent_score"`
	Status        string    `db:"status"`
	IsDerby       bool      `db:"is_derby"`
	WinType       string    `db:"win_type"`
}

// StadiumSector структура сектора стадиона
type StadiumSector struct {
	ID               string  `db:"id"`
	SectorNumber     string  `db:"sector_number"`
	Capacity         int     `db:"capacity"`
	SectorType       string  `db:"sector_type"`
	PriceCoefficient float64 `db:"price_coefficient"`
	ColorCode        string  `db:"color_code"`
}

// Seat структура места
type Seat struct {
	ID                   string `db:"id"`
	SectorID             string `db:"sector_id"`
	SeatRow              string `db:"seat_row"`
	SeatNumber           string `db:"seat_number"`
	IsHandicapAccessible bool   `db:"is_handicap_accessible"`
	IsTaken              bool   `db:"is_taken"`
}

// Ticket структура билета
type Ticket struct {
	ID           string    `db:"id"`
	TicketNumber string    `db:"ticket_number"`
	UserID       string    `db:"user_id"`
	MatchID      string    `db:"match_id"`
	SeatID       string    `db:"seat_id"`
	FinalPrice   int       `db:"final_price"`
	FullName     string    `db:"full_name"`
	Phone        string    `db:"phone"`
	Email        string    `db:"email"`
	PurchaseDate time.Time `db:"purchase_date"`
	Status       string    `db:"status"`
	QRCodeHash   string    `db:"qr_code"`
}
