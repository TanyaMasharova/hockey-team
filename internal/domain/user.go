package domain

import "time"

type User struct {
    ID           string     `db:"id" json:"id"`
    Phone        string     `db:"phone" json:"phone"`
    PasswordHash string     `db:"password_hash" json:"-"`
    Email        string     `db:"email" json:"email"`
    FullName     string     `db:"full_name" json:"full_name"`
    BirthDate    *string    `db:"birth_date" json:"birth_date"` // добавляем дату рождения
    Role      string    `db:"role" json:"role"`
    CreatedAt    time.Time  `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}

type RegistrationData struct {
    Phone    string 
    Password string 
    Email    string 
    FullName string 
}

type LoginData struct {
    Phone    string
    Password string
}

// НОВАЯ СТРУКТУРА: данные для обновления профиля
type UpdateProfileData struct {
    FullName  *string `json:"full_name,omitempty"`
    Phone     *string `json:"phone,omitempty"`
    Email     *string `json:"email,omitempty"`
    BirthDate *string `json:"birth_date,omitempty"`
}