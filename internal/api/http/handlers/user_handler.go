package handlers

import (
	"errors"
	"net/http"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/dto"
	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/TanyaMasharova/hockey-team/internal/service/auth"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	userService *auth.Service
	logger *logrus.Logger
}

func NewUserHandler(authService *auth.Service, logger *logrus.Logger) *UserHandler{
	if authService == nil {
		panic("userService is required")
	}
	return &UserHandler {
		userService: authService,
		logger: logger,
	}
}

//r - указатель на запрос, содержит в себе заголовки, тело 
//w - интерфейс, содержит в себе структуру для отправки ответа
func (h *UserHandler) Register(c *gin.Context) {

	//1. Декодирование JSON
	//явно создаём переменную, т.к. данные приходят извне
	var req dto.RegisterRequest 
	//NewDecoder() - конструктор, создает новый декодер
	//Decode() - метод
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warn("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H {
			"error" : "Invalud request body: " + err.Error(),
				})
				return
	}

		// Логируем входящий запрос
	h.logger.WithFields(logrus.Fields{
		"phone":    req.Phone,
		"email":    req.Email,
		"fullName": req.FullName,
	}).Info("Processing registration request")

	regData := domain.RegistrationData{
		Phone: req.Phone,
		Password: req.Password,
		Email: req.Email,
		FullName: req.FullName,
	}

	user, err := h.userService.Register(c.Request.Context(), regData)
	if err != nil {

		h.logger.WithError(err).Error("Registration failed")
		
		h.logger.Errorf("Error type: %T", err)
		 switch {
        case errors.Is(err, auth.ErrPhoneAlreadyRegistered):
           c.JSON(http.StatusConflict, gin.H{
						"error" : "Phone number already registered",
					 })
        case errors.Is(err, auth.ErrEmailAlreadyRegistered):
            c.JSON(http.StatusConflict, gin.H{
							"error" : "Email is already registered",
						})
        case errors.Is(err, auth.ErrValidationFailed):
            c.JSON(http.StatusBadRequest, gin.H{
							"error" : err.Error(),
						})
        default:
            c.JSON(http.StatusInternalServerError, gin.H{
							"error" : "Internal server error",
						})
        }
				return
	}

	resp := dto.RegisterResponse {
		ID: user.ID,
		Phone: user.Phone,
		Email: user.Email,
		FullName: user.FullName,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, resp)

}

