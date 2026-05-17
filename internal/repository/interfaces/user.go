package interfaces

import (
	"context"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Exists(ctx context.Context, phone string) (bool, error)
	
	// НОВЫЕ МЕТОДЫ ДЛЯ РЕДАКТИРОВАНИЯ
	Update(ctx context.Context, user *domain.User) error
	UpdateField(ctx context.Context, id string, field string, value interface{}) error
}

// Ошибки репозитория
var (
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrDuplicatePhone    = fmt.Errorf("phone number already exists")
	ErrDuplicateEmail    = fmt.Errorf("email already exists")
)