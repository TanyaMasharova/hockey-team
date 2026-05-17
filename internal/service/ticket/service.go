package ticket

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/TanyaMasharova/hockey-team/internal/repository/interfaces"
)

type CreateTicketRequest struct {
	UserID     string
	MatchID    string
	SeatID     string
	FinalPrice int
}

type Service struct {
	ticketRepo interfaces.TicketRepository
}

func NewService(ticketRepo interfaces.TicketRepository) *Service {
	if ticketRepo == nil {
		panic("ticketRepo is required")
	}
	return &Service{
		ticketRepo: ticketRepo,
	}
}

func (s *Service) GetUserTickets(ctx context.Context, userID string) ([]*domain.TicketWithDetails, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	
	tickets, err := s.ticketRepo.GetUserTickets(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tickets: %w", err)
	}
	
	return tickets, nil
}

func (s *Service) CreateTicket(ctx context.Context, req *CreateTicketRequest) (*domain.Ticket, error) {
	// Валидация
	if req.UserID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	if req.MatchID == "" {
		return nil, fmt.Errorf("match ID is required")
	}
	if req.SeatID == "" {
		return nil, fmt.Errorf("seat ID is required")
	}
	
	// Проверяем доступность места
	available, err := s.ticketRepo.CheckSeatAvailability(ctx, req.SeatID, req.MatchID)
	if err != nil {
		return nil, fmt.Errorf("failed to check seat availability: %w", err)
	}
	
	if !available {
		return nil, fmt.Errorf("seat is already taken")
	}
	
	// Генерируем хэш для QR кода
	qrHash := generateQRHash(req.UserID, req.MatchID, req.SeatID)
	
	ticket := &domain.Ticket{
		ID:           generateID(),
		UserID:       req.UserID,
		MatchID:      req.MatchID,
		SeatID:       req.SeatID,
		FinalPrice:   req.FinalPrice,
		PurchaseDate: time.Now(),
		QRCodeHash:   qrHash,
		Status:       "active",
	}
	
	err = s.ticketRepo.CreateTicket(ctx, ticket)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}
	
	return ticket, nil
}

func generateQRHash(userID, matchID, seatID string) string {
	data := fmt.Sprintf("%s:%s:%s:%d", userID, matchID, seatID, time.Now().UnixNano())
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func generateID() string {
	bytes := make([]byte, 16)
	// Используйте uuid, а не rand
	// Лучше использовать github.com/google/uuid
	return fmt.Sprintf("%x", bytes)
}