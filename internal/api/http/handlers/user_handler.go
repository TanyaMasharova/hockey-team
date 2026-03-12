package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TanyaMasharova/hockey-team/internal/api/http/dto"
	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/TanyaMasharova/hockey-team/internal/service/auth"
)

type UserHandler struct {
	userService *auth.Service
}

func NewUserHandler(authService *auth.Service) *UserHandler{
	if authService == nil {
		panic("userService is required")
	}
	return &UserHandler {
		userService: authService,
	}
}

//r - указатель на запрос, содержит в себе заголовки, тело 
//w - интерфейс, содержит в себе структуру для отправки ответа
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	//1. Проверка метода
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//2. Декодирование JSON
	//явно создаём переменную, т.к. данные приходят извне
	var req dto.RegisterRequest 
	//NewDecoder() - конструктор, создает новый декодер
	//Decode() - метод
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() //закрываем тело запроса явно

	//3. Преобразование в доменную структуру
	/*
	var regData domain.RegistrationData        // 1. Только объявление
regData := domain.RegistrationData{}       // 2. Объявление + пустая инициализация
regData := domain.RegistrationData{        // 3. Объявление + заполнение
    Phone: req.Phone,
}
		*/
	regData := domain.RegistrationData{
		Phone: req.Phone,
		Password: req.Password,
		Email: req.Email,
		FullName: req.FullName,
	}

	user, err := h.userService.Register(r.Context(), regData)
	if err != nil {
		 switch {
        case errors.Is(err, auth.ErrPhoneAlreadyRegistered):
            http.Error(w, "Phone number already registered", http.StatusConflict)
        case errors.Is(err, auth.ErrEmailAlreadyRegistered):
            http.Error(w, "Email already registered", http.StatusConflict)
        case errors.Is(err, auth.ErrValidationFailed):
            http.Error(w, err.Error(), http.StatusBadRequest)
        default:
            // Логируем внутреннюю ошибку (в реальном проекте используйте логгер)
            // log.Printf("Internal error: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
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

	//Header() устанавливает заголовок (в даннос случае говорящий про json)
	//WriteHeader() устанавливает http статус-код 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)


}

