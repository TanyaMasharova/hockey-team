package postgres

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/TanyaMasharova/hockey-team/internal/repository/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, phone, email, password_hash, full_name, birth_date, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, interfaces.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (phone, password_hash, email, full_name, birth_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		user.Phone,
		user.PasswordHash,
		nullString(user.Email),
		user.FullName,
		nullStringPtr(user.BirthDate),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
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

func (r *userRepo) Exists(ctx context.Context, phone string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)`
	err := r.db.GetContext(ctx, &exists, query, phone)
	if err != nil {
		return false, fmt.Errorf("check phone existence: %w", err)
	}
	return exists, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	query := `
		SELECT id, phone, email, password_hash, full_name, birth_date, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, interfaces.ErrUserNotFound
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &user, nil
}

// НОВЫЙ МЕТОД: полное обновление пользователя
func (r *userRepo) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET full_name = $1, phone = $2, email = $3, birth_date = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		user.FullName,
		user.Phone,
		nullString(user.Email),
		nullStringPtr(user.BirthDate),
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return interfaces.ErrUserNotFound
		}

		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				if pgErr.Constraint == "users_phone_key" {
					return interfaces.ErrDuplicatePhone
				}
				if pgErr.Constraint == "users_email_key" {
					return interfaces.ErrDuplicateEmail
				}
			}
		}
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

// НОВЫЙ МЕТОД: обновление конкретного поля
func (r *userRepo) UpdateField(ctx context.Context, id string, field string, value interface{}) error {
	var query string

	switch field {
	case "full_name":
		query = `UPDATE users SET full_name = $1, updated_at = NOW() WHERE id = $2`
	case "phone":
		query = `UPDATE users SET phone = $1, updated_at = NOW() WHERE id = $2`
	case "email":
		query = `UPDATE users SET email = $1, updated_at = NOW() WHERE id = $2`
	case "birth_date":
		query = `UPDATE users SET birth_date = $1, updated_at = NOW() WHERE id = $2`
	default:
		return fmt.Errorf("unknown field: %s", field)
	}

	result, err := r.db.ExecContext(ctx, query, value, id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				if field == "phone" && pgErr.Constraint == "users_phone_key" {
					return interfaces.ErrDuplicatePhone
				}
				if field == "email" && pgErr.Constraint == "users_email_key" {
					return interfaces.ErrDuplicateEmail
				}
			}
		}
		return fmt.Errorf("update field %s: %w", field, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return interfaces.ErrUserNotFound
	}

	return nil
}

func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullStringPtr(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}
