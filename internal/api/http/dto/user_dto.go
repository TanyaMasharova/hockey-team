package dto

import "time"

// UserResponse - ответ с данными пользователя
type UserResponse struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	BirthDate *string   `json:"birth_date,omitempty"`
	Role      string    `db:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateProfileRequest - запрос на обновление профиля (несколько полей)
type UpdateProfileRequest struct {
	FullName  *string `json:"full_name,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Email     *string `json:"email,omitempty"`
	BirthDate *string `json:"birth_date,omitempty"`
}

// UpdateProfileResponse - ответ после обновления профиля
type UpdateProfileResponse struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	BirthDate *string   `json:"birth_date,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateFieldRequest - запрос на обновление конкретного поля
type UpdateFieldRequest struct {
	Field string `json:"field" binding:"required"` // full_name, phone, email, birth_date
	Value string `json:"value" binding:"required"`
}

// UpdateFieldResponse - ответ после обновления поля
type UpdateFieldResponse struct {
	Message string `json:"message"`
	Field   string `json:"field"`
	Value   string `json:"value"`
}
