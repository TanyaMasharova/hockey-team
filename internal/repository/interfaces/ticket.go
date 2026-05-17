package interfaces

import (
	"context"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
)

type TicketRepository interface {
	// GetUserTickets возвращает все билеты пользователя с детальной информацией
	GetUserTickets(ctx context.Context, userID string) ([]*domain.TicketWithDetails, error)
	
	// CreateTicket создает новый билет
	CreateTicket(ctx context.Context, ticket *domain.Ticket) error
	
	// CheckSeatAvailability проверяет, свободно ли место на конкретный матч
	CheckSeatAvailability(ctx context.Context, seatID string, matchID string) (bool, error)
	
	// GetTicketByID возвращает билет по ID (опционально)
	GetTicketByID(ctx context.Context, ticketID string) (*domain.Ticket, error)
}