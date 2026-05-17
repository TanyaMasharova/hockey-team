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

func (h *UserHandler) Login(c *gin.Context) {
    // 1. Декодирование JSON
    var req dto.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.WithError(err).Warn("Invalid login request body")
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body: " + err.Error(),
        })
        return
    }

    // Логируем входящий запрос (без пароля!)
    h.logger.WithFields(logrus.Fields{
        "email": req.Email,
    }).Info("Processing login request")

    // 2. Вызов сервиса аутентификации
    loginData, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        h.logger.WithError(err).Error("Login failed")
        
        // Обработка ошибок
        switch {
        case errors.Is(err, auth.ErrInvalidCredentials):
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid email or password",
            })
        case errors.Is(err, auth.ErrUserNotFound):
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "User not found",
            })
        default:
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal server error",
            })
        }
        
        return
    }

    // 3. Успешный ответ с токеном и id пользователя
    resp := dto.LoginResponse{
        ID:        loginData.User.ID,
        Token:     loginData.Token,
        TokenType: "Bearer",
        Role:        loginData.User.Role,

    }
    h.logger.WithFields(logrus.Fields{
        "user_id": loginData.User.ID,
        "role":    loginData.User.Role,
    }).Info("Login successful, sending response")

    c.JSON(http.StatusOK, resp)
}


// В user_handler.go обновите GetUserByID:
func (h *UserHandler) GetUserByID(c *gin.Context) {
    id := c.Param("id")
    
    if id == "" {
        h.logger.Warn("Empty user ID")
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "User ID is required",
        })
        return
    }
    
    h.logger.WithFields(logrus.Fields{
        "user_id": id,
    }).Info("Processing get user by ID request")
    
    user, err := h.userService.GetUserByID(c.Request.Context(), id)
    if err != nil {
        h.logger.WithError(err).Error("Failed to get user")
        
        switch {
        case errors.Is(err, auth.ErrUserNotFound):
            c.JSON(http.StatusNotFound, gin.H{
                "error": "User not found",
            })
        default:
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal server error",
            })
        }
        return
    }
    
    // Используем обновленный UserResponse с BirthDate
    resp := dto.UserResponse{
        ID:        user.ID,
        Phone:     user.Phone,
        Email:     user.Email,
        FullName:  user.FullName,
        BirthDate: user.BirthDate,
        Role:      user.Role,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
    
    c.JSON(http.StatusOK, resp)
}


func (h *UserHandler) UpdateField(c *gin.Context) {
    // Берем user_id из параметра URL
    userID := c.Param("user_id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "User ID is required",
        })
        return
    }
    
    var req dto.UpdateFieldRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.WithError(err).Warn("Invalid update field request")
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body: " + err.Error(),
        })
        return
    }
    
    h.logger.WithFields(logrus.Fields{
        "user_id": userID,
        "field":   req.Field,
        "value":   req.Value,
    }).Info("Processing field update request")
    
    err := h.userService.UpdateField(c.Request.Context(), userID, req.Field, req.Value)
    if err != nil {
        h.logger.WithError(err).Error("Field update failed")
        
        switch {
        case errors.Is(err, auth.ErrUserNotFound):
            c.JSON(http.StatusNotFound, gin.H{
                "error": "User not found",
            })
        case errors.Is(err, auth.ErrPhoneAlreadyRegistered):
            c.JSON(http.StatusConflict, gin.H{
                "error": "Phone number already registered",
            })
        case errors.Is(err, auth.ErrEmailAlreadyRegistered):
            c.JSON(http.StatusConflict, gin.H{
                "error": "Email already registered",
            })
        case errors.Is(err, auth.ErrValidationFailed):
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
        default:
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Internal server error",
            })
        }
        return
    }
    
    c.JSON(http.StatusOK, dto.UpdateFieldResponse{
        Message: "Field updated successfully",
        Field:   req.Field,
        Value:   req.Value,
    })
}