package handlers

import (
	"net/http"

	"github.com/TanyaMasharova/hockey-team/internal/service/seat"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SeatHandler struct {
	seatService *seat.Service
	logger      *logrus.Logger
}

func NewSeatHandler(seatService *seat.Service, logger *logrus.Logger) *SeatHandler {
	return &SeatHandler{
		seatService: seatService,
		logger:      logger,
	}
}

func (h *SeatHandler) GetSeatsBySector(c *gin.Context) {
	sectorID := c.Param("sectorId")
	matchID := c.Query("matchId") // получаем matchId из query параметра

	if sectorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sector ID is required"})
		return
	}

	if matchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "match ID is required"})
		return
	}

	seats, err := h.seatService.GetSeatsBySector(c.Request.Context(), sectorID, matchID)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get seats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get seats"})
		return
	}

	c.JSON(http.StatusOK, seats)
}
