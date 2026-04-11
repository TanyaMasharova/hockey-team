package postgres

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/TanyaMasharova/hockey-team/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//указатель на подключение к БД. приватная структыра (т.к. начинается с маленькой буквы)
type userRepo struct {
    db *sqlx.DB
}

//Пуюбличный конструктор (т.к. с большой буквы), который инициализирующий db внутри структуры тем, что передано в параметрах
func NewUserRepository(db *sqlx.DB) *userRepo {
    return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (phone, password_hash, email, full_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	//QueryRowContext - запрос который ыернет не более 1 строки и использует context для отмены операций и установки таймаутов
	//Ошибки, в том числе ошибки, когда строки не найдены, откладываются до вызова Scan-метода возвращаемого значения
	err := r.db.QueryRowContext(ctx, query,
	user.Phone,
    user.PasswordHash,
    nullString(user.Email), 
	user.FullName,
).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
//формирует запрос и результаты кладет в &user.ID, &user.CreatedAt, &user.UpdatedAt

if err != nil {
        // Преобразуем ошибки БД в ошибки репозитория
				//проверка, является ли ошибка специфичной 
				//ok = true если это ошибка Postgres
        if pgErr, ok := err.(*pq.Error); ok {
            if pgErr.Code == "23505" { // unique_violation. код нарушения уникальности 
							//pgErr содержит имя ограничения. на основе его возвращается определённое значение
                if pgErr.Constraint == "users_phone_key" {
                    return interfaces.ErrDuplicatePhone
                }
                if pgErr.Constraint == "users_email_key" {
                    return interfaces.ErrDuplicateEmail
                }
            }
        }
        return fmt.Errorf("create user: %w", err)
    }
return nil
}
//проверка сущетвования пользователя по номеру телефона. (r *userRepo) - получатель метода
func (r *userRepo)  Exists(ctx context.Context, phone string) (bool, error) {
	var exists bool
	//запрос, возвращающий наличие хотя бы 1 строки где телефон = переданному аргументу (т.е. введённому пользователем на фронте номеру)
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`

		//сканирует результат в переменную, переданную по ссылке
    err := r.db.GetContext(ctx, &exists, query, phone)
		//ошибки (например): пользователь закрыл браузер, запрос отменён, истёк таймаут
    if err != nil {
        return false, fmt.Errorf("check phone existence: %w", err)
    }

    return exists, nil
}

//вспомогательная функция. если передали пустое значение - возвращаем nil, если нет - возвращаем переданную строку
func nullString(s string) interface{} {
    if s == "" {
        return nil
    }
    return s
}