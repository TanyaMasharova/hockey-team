package handlers

import (
	"net/http"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/dto"
	"github.com/TanyaMasharova/hockey-team/internal/service/matches"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


type MatchesHandler struct {
	// userService *auth.Service
	matchService *matches.Service
	logger *logrus.Logger
}

func NewMatchesHandler(matchService *matches.Service, logger *logrus.Logger) *MatchesHandler{
	if matchService == nil {
		panic("matchService is required")
	}
	return &MatchesHandler {
		matchService: matchService,
		logger: logger,
	}
}
//gin.Context помогает отслеживатьс= состояние запросов и ответов (http)
func (h *MatchesHandler) GetMatches(c *gin.Context){
	matches, err := h.matchService.GetMatches(c.Request.Context())
	if err != nil {
				h.logger.WithError(err).Error("GetMatches failed")
	}

	var resp []dto.MatchResponse
	for _, match := range matches {
		resp = append(resp, dto.MatchResponse{
			ID:             match.ID,
			Opponent:     match.Opponent,
			LogoOpponent: match.LogoOpponent,
			MatchDate:      match.MatchDate,
			HomeAway:       match.HomeAway,
			OurScore:       match.OurScore,
			OpponentScore:  match.OpponentScore,
			// Season:         match.Season,
			Status:         match.Status,
			IsDerby:        match.IsDerby,
		})
	}
	c.JSON(http.StatusOK, resp)
}