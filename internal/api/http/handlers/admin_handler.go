// internal/api/http/handlers/admin.go
package handlers

import (
	"net/http"

	"github.com/TanyaMasharova/hockey-team/internal/service/admin"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AdminHandler struct {
	adminService *admin.Service
	logger       *logrus.Logger
}

func NewAdminHandler(adminService *admin.Service, logger *logrus.Logger) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		logger:       logger,
	}
}

func (h *AdminHandler) GetStatsSummary(c *gin.Context) {
	stats, err := h.adminService.GetStatsSummary(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get stats summary")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) GetAllStats(c *gin.Context) {
	// Параллельно получаем все данные
	salesByMonth, err := h.adminService.GetSalesByMonth(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get sales by month")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}

	sectorPopularity, err := h.adminService.GetSectorPopularity(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get sector popularity")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}

	ticketStatus, err := h.adminService.GetTicketStatus(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get ticket status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}

	avgPriceBySector, err := h.adminService.GetAvgPriceBySector(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("Failed to get avg price by sector")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}

	topBuyers, err := h.adminService.GetTopBuyers(c.Request.Context(), 10)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get top buyers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sales_by_month":      salesByMonth,
		"sector_popularity":   sectorPopularity,
		"ticket_status":       ticketStatus,
		"avg_price_by_sector": avgPriceBySector,
		"top_buyers":          topBuyers,
	})
}
