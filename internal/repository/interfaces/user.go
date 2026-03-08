package interfaces

import (
	"context"

	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
)

type UserRepository interface {

	Create(ctx context.Context, user *domain.User) error

	//  GetByID(ctx context.Context, id int64) (*domain.User, error)

	//  GetByPhone(ctx context.Context, phone string) (*domain.User, error)

	 // GetByEmail получает пользователя по email (для проверки дубликатов)
    // GetByEmail(ctx context.Context, email string) (*domain.User, error)

		// UpdatePassword обновляет только пароль
    // Ошибка: ErrNotFound
    // UpdatePassword(ctx context.Context, id int64, passwordHash string) error

	//  Delete(ctx context.Context, id int64) error

	 // Exists проверяет существование пользователя по телефону
    Exists(ctx context.Context, phone string) (bool, error)

		
    // GetAll возвращает список пользователей с пагинацией
    // GetAll(ctx context.Context, limit, offset int) ([]*domain.User, int64, error)
}
// Ошибки репозитория
var (
    ErrUserNotFound      = fmt.Errorf("user not found")
    ErrDuplicatePhone    = fmt.Errorf("phone number already exists")
    ErrDuplicateEmail    = fmt.Errorf("email already exists")
)