package dto

import "time"

type TicketResponse struct {
    ID           string    `json:"id"`
    MatchDate    time.Time `json:"match_date"`
    OpponentName string    `json:"opponent_name"`
    HomeAway     string    `json:"home_away"`
    SectorNumber string    `json:"sector_number"`
    SeatRow      string    `json:"seat_row"`
    SeatNumber   string    `json:"seat_number"`
    PurchaseDate time.Time `json:"purchase_date"`
    FinalPrice   float64   `json:"final_price"`
    Status       string    `json:"status"`
}

type UserTicketsResponse struct {
	UserID  string           `json:"user_id"`
	Tickets []TicketResponse `json:"tickets"`
	Total   int              `json:"total"`
}
type CreateTicketRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	MatchID    string `json:"match_id" binding:"required"`
	SeatID     string `json:"seat_id" binding:"required"`
	FinalPrice int    `json:"final_price" binding:"required"`
}

// CreateTicketResponse ответ после создания билета
type CreateTicketResponse struct {
    ID           string    `json:"id"`
    TicketNumber string    `json:"ticket_number"`
    QRCode       string    `json:"qr_code"`
    Status       string    `json:"status"`
    PurchaseDate time.Time `json:"purchase_date"`
}

// MatchResponse ответ с информацией о матче
// type MatchResponse struct {
//     ID            string `json:"id"`
//     OpponentName  string `json:"opponent_name"`
//     MatchDate     string `json:"match_date"`
//     HomeAway      string `json:"home_away"`
//     OurScore      int    `json:"our_score"`
//     OpponentScore int    `json:"opponent_score"`
//     Status        string `json:"status"`
//     IsDerby       bool   `json:"is_derby"`
// }

// SectorResponse ответ с информацией о секторе
type SectorResponse struct {
    ID               string  `json:"id"`
    SectorNumber     string  `json:"sector_number"`
    Capacity         int     `json:"capacity"`
    SectorType       string  `json:"sector_type"`
    PriceCoefficient float64 `json:"price_coefficient"`
    ColorCode        string  `json:"color_code"`
}

// SeatResponse ответ с информацией о месте
type SeatResponse struct {
    ID                   string `json:"id"`
    SeatRow              string `json:"seat_row"`
    SeatNumber           string `json:"seat_number"`
    IsHandicapAccessible bool   `json:"is_handicap_accessible"`
    IsTaken              bool   `json:"is_taken"`
}