package domain

import "time"

// User — просто хранилище данных, никакой логики
type User struct {
    ID           string     `db:"id" json:"id"`
    Phone        string    `db:"phone" json:"phone"`
    PasswordHash string    `db:"password_hash" json:"-"`
    Email        string    `db:"email" json:"email"`
    FullName     string    `db:"full_name" json:"full_name"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// RegistrationData — данные для регистрации (отдельная структура!)
type RegistrationData struct {
    Phone    string 
    Password string 
    Email    string 
    FullName string 
}

// LoginData — данные для входа
type LoginData struct {
    Phone    string
    Password string
}

// UpdateData — данные для обновления
// type UpdateData struct {
//     Email    *string // указатели, чтобы отличать пустую строку от "не меняем"
//     FullName *string
// }