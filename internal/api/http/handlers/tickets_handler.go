package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/dto"
	"github.com/TanyaMasharova/hockey-team/internal/service/ticket"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TicketHandler struct {
	ticketService *ticket.Service
	logger        *logrus.Logger
}

func NewTicketHandler(ticketService *ticket.Service, logger *logrus.Logger) *TicketHandler {
	if ticketService == nil {
		panic("ticketService is required")
	}
	return &TicketHandler{
		ticketService: ticketService,
		logger:        logger,
	}
}

func (h *TicketHandler) GetUserTickets(c *gin.Context) {
    userID := c.Param("user_id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
        return
    }
    
    tickets, err := h.ticketService.GetUserTickets(c.Request.Context(), userID)
    if err != nil {
        h.logger.WithError(err).Error("Failed to get user tickets")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tickets"})
        return
    }
    
    response := dto.UserTicketsResponse{
        UserID:  userID,
        Total:   len(tickets),
        Tickets: make([]dto.TicketResponse, 0, len(tickets)),
    }
    
    for _, t := range tickets {
        homeAwayText := "В гостях"
        if t.HomeAway == "home" {
            homeAwayText = "Дома"
        }
        
        response.Tickets = append(response.Tickets, dto.TicketResponse{
            ID:           t.ID,
            MatchDate:    t.MatchDate,
            OpponentName: t.OpponentName,
            HomeAway:     homeAwayText,
            SectorNumber: t.SectorNumber,
            SeatRow:      t.SeatRow,
            SeatNumber:   t.SeatNumber,
            PurchaseDate: t.PurchaseDate,
            FinalPrice:   t.FinalPrice,
            Status:       t.Status,
        })
    }
    
    c.JSON(http.StatusOK, response)
}
func (h *TicketHandler) CreateTicket(c *gin.Context) {
    var req dto.CreateTicketRequest

    body, _ := io.ReadAll(c.Request.Body)
    fmt.Println("Raw body:", string(body))
    
    // Восстанавливаем тело для ShouldBindJSON
    c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
    
    if err := c.ShouldBindJSON(&req); err != nil {
        // Выводим конкретную ошибку
        fmt.Printf("Bind error: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})  // ← возвращаем реальную ошибку
        return
    }
    
    fmt.Printf("Parsed req: %+v\n", req)
    
    serviceReq := &ticket.CreateTicketRequest{
        UserID:     req.UserID,
        MatchID:    req.MatchID,
        SeatID:     req.SeatID,
        FinalPrice: req.FinalPrice,
    }
    
    createdTicket, err := h.ticketService.CreateTicket(c.Request.Context(), serviceReq)
    if err != nil {
        h.logger.WithError(err).Error("Failed to create ticket")
        
        if err.Error() == "seat is already taken" {
            c.JSON(http.StatusConflict, gin.H{"error": "Seat is already taken"})
            return
        }
        
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "id":            createdTicket.ID,
        "status":        createdTicket.Status,
        "purchase_date": createdTicket.PurchaseDate,
    })
}