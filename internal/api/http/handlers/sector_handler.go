package handlers

import (
	"net/http"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/dto"
	"github.com/TanyaMasharova/hockey-team/internal/service/sector"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SectorHandler struct {
    sectorService *sector.Service
    logger        *logrus.Logger
}

func NewSectorHandler(sectorService *sector.Service, logger *logrus.Logger) *SectorHandler {
    return &SectorHandler{
        sectorService: sectorService,
        logger:        logger,
    }
}

func (h *SectorHandler) GetAllSectors(c *gin.Context) {
    sectors, err := h.sectorService.GetAllSectors(c.Request.Context())
    if err != nil {
        h.logger.WithError(err).Error("Failed to get sectors")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get sectors"})
        return
    }
    
    response := make([]dto.SectorResponse, 0, len(sectors))
    for _, s := range sectors {
        response = append(response, dto.SectorResponse{
            ID:               s.ID,
            SectorNumber:     s.SectorNumber,
            Capacity:         s.Capacity,
            SectorType:       s.SectorType,
            PriceCoefficient: s.PriceCoefficient,
            ColorCode:        s.ColorCode,
        })
    }
    
    c.JSON(http.StatusOK, response)
}