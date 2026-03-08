package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/TanyaMasharova/hockey-team/internal/repository/interfaces"
)


type Service struct {
    userRepo interfaces.UserRepository
}

//конструктор
func NewService(userRepo interfaces.UserRepository) *Service {
    if userRepo == nil {
        panic("userRepo is required") // защита от ошибок
    }
    
    return &Service{
        userRepo: userRepo,
    }
}

func (s *Service) Register(ctx context.Context, data domain.RegistrationData) (*domain.User, error) {
	if err := s.validatePhone(data.Phone); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	if err := s.validatePassword(data.Password); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	//так же можно сделать валидацию на email

	if data.FullName == "" {
        return nil, fmt.Errorf("validation failed: full name required")
    }

		exists, err := s.userRepo.Exists(ctx, data.Phone)
		if err != nil {
			return nil, fmt.Errorf("failed to check phone: %w", err)

		}
		if exists {
			return nil, fmt.Errorf("phone already registered")	//бизнес-ошибка
		}
		// 3. Хеширование
    // hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
    // if err != nil {
    //     return nil, fmt.Errorf("failed to hash password: %w", err)
    // }

		user := &domain.User{
        Phone:        data.Phone,
				//PasswordHash: string(hash), это будет когда реализую хэш

        PasswordHash: data.Password,
        Email:        data.Email,
        FullName:     data.FullName,
    }

		if err := s.userRepo.Create(ctx, user); err != nil {
			//преобразование технических ошибок в бизнес
			if errors.Is(err, interfaces.ErrDuplicatePhone) {
            return nil, fmt.Errorf("phone already registered")
        }
        if errors.Is(err, interfaces.ErrDuplicateEmail) {
            return nil, fmt.Errorf("email already registered")
        }
        return nil, fmt.Errorf("database error: %w", err)
		}

		// user.PasswordHash = "" не поняла зачем это тут 
    return user, nil
}

func (s *Service) validatePhone(phone string) error {
	 //заглушка. тут должна быть проверка на коректность номера телефона
    // if phone == "" {
    //     return domain.ErrPhoneRequired
    // }
    // // Простая проверка: только цифры, может быть с +
    // phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
    // if !phoneRegex.MatchString(phone) {
    //     return domain.ErrInvalidPhone
    // }
    return nil
}

func (s *Service) validatePassword(password string) error {
	//заглушка. проверка на кол-во символов в пароле, наличие спец. символов
	return nil
}