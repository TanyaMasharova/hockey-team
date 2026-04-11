package dto

//приходящие данные при регистрации
type RegisterRequest struct {
    Phone    string `json:"phone" validate:"required"`
    Password string `json:"password" validate:"required,min=6"`
    Email    string `json:"email" validate:"required,email"`
    FullName string `json:"full_name" validate:"required"`
}

//ответ на фронт
type RegisterResponse struct {
    ID        string  `json:"id"`
    Phone     string `json:"phone"`
    Email     string `json:"email"`
    FullName  string `json:"full_name"`
    CreatedAt string `json:"created_at"`
}

//приходящие данные при входе
type LoginRequest struct {
    Phone    string `json:"phone" validate:"required"`
    Password string `json:"password" validate:"required"`
}