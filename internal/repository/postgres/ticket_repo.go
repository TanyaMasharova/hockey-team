package postgres

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ticketRepo struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) *ticketRepo {
	return &ticketRepo{db: db}
}

// GetUserTickets возвращает все билеты пользователя с детальной информацией
func (r *ticketRepo) GetUserTickets(ctx context.Context, userID string) ([]*domain.TicketWithDetails, error) {
    var tickets []*domain.TicketWithDetails
    
    query := `
        SELECT 
            t.id as ticket_id,
            t.final_price,
            t.purchase_date,
            t.status as ticket_status,
            m.match_date,
            m.home_away,
            o.name as opponent_name,
            COALESCE(o.city, '') as opponent_city,
            ss.sector_number,
            ss.sector_type,
            s.seat_row,
            s.seat_number
        FROM tickets t
        JOIN matches m ON t.match_id = m.id
        JOIN opponents o ON m.opponent_id = o.id
        JOIN seats s ON t.seat_id = s.id
        JOIN stadium_sectors ss ON s.sector_id = ss.id
        WHERE t.user_id = $1
        ORDER BY m.match_date DESC
    `
    
    err := r.db.SelectContext(ctx, &tickets, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user tickets: %w", err)
    }
    
    return tickets, nil
}

// CreateTicket создает новый билет (только поля из БД)
func (r *ticketRepo) CreateTicket(ctx context.Context, ticket *domain.Ticket) error {
    query := `
        INSERT INTO tickets (
            id, user_id, match_id, seat_id, final_price, 
            purchase_date, qr_code_hash, status
        ) VALUES (
            gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7
        )
        RETURNING id
    `
    
    var id string
    err := r.db.QueryRowContext(ctx, query,
        ticket.UserID,
        ticket.MatchID,
        ticket.SeatID,
        ticket.FinalPrice,
        ticket.PurchaseDate,
        ticket.QRCodeHash,
        ticket.Status,
    ).Scan(&id)
    
    if err != nil {
        return fmt.Errorf("failed to create ticket: %w", err)
    }
    
    ticket.ID = id
    return nil
}

// CheckSeatAvailability проверяет, свободно ли место на конкретный матч
func (r *ticketRepo) CheckSeatAvailability(ctx context.Context, seatID string, matchID string) (bool, error) {
	var count int
	
	query := `
		SELECT COUNT(*)
		FROM tickets t
		WHERE t.seat_id = $1 AND t.match_id = $2 AND t.status != 'cancelled'
	`
	
	err := r.db.GetContext(ctx, &count, query, seatID, matchID)
	if err != nil {
		return false, fmt.Errorf("failed to check seat availability: %w", err)
	}
	
	return count == 0, nil
}

func (r *ticketRepo) GetTicketByID(ctx context.Context, ticketID string) (*domain.Ticket, error) {
	var ticket domain.Ticket
	
	query := `
		SELECT 
			id, user_id, match_id, seat_id, final_price,
			purchase_date, qr_code_hash, status
		FROM tickets
		WHERE id = $1
	`
	
	err := r.db.GetContext(ctx, &ticket, query, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket by id: %w", err)
	}
	
	return &ticket, nil
}