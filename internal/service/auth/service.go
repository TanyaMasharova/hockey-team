package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TanyaMasharova/hockey-team/internal/domain"
	"github.com/TanyaMasharova/hockey-team/internal/repository/interfaces"
	"github.com/sirupsen/logrus"
)

type Service struct {
	userRepo interfaces.UserRepository //интерфейсы - это тоже указатели (под капотом)
}

// конструктор
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
		return nil, fmt.Errorf("phone already registered") //бизнес-ошибка
	}
	// 3. Хеширование
	// hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to hash password: %w", err)
	// }

	user := &domain.User{
		Phone: data.Phone,
		//PasswordHash: string(hash), это будет когда реализую хэш

		PasswordHash: data.Password,
		Email:        data.Email,
		FullName:     data.FullName,
		Role:         "user",
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

type LoginResponseData struct {
	Token string
	User  *domain.User
	Role  string
}

func (s *Service) Login(ctx context.Context, email, password string) (*LoginResponseData, error) {
	// 1. Валидация входных данных
	if email == "" {
		return nil, fmt.Errorf("%w: email is required", ErrValidationFailed)
	}
	if password == "" {
		return nil, fmt.Errorf("%w: password is required", ErrValidationFailed)
	}

	// 2. Поиск пользователя по email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, interfaces.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// 3. Проверка пароля
	if user.PasswordHash != password {
		return nil, ErrInvalidCredentials
	}

	// 4. Генерация токена
	token, err := s.generateJWTToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	logrus.Info("user role: ", user.Role)

	return &LoginResponseData{
		Token: token,
		User:  user,
		Role:  user.Role,
	}, nil
}

// GetUserByID получает пользователя по ID (как строке)
func (s *Service) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: invalid user id", ErrValidationFailed)
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, interfaces.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Очищаем sensitive данные
	user.PasswordHash = ""

	return user, nil
}

// Вспомогательный метод для генерации JWT токена (пока заглушка)
func (s *Service) generateJWTToken(user *domain.User) (string, error) {
	// TODO: Реализовать генерацию JWT токена
	// Пока возвращаем временный токен для тестирования
	// В реальном коде нужно использовать библиотеку github.com/golang-jwt/jwt

	// Это временное решение для тестирования
	// tempToken := fmt.Sprintf("temp_token_for_user_%d", user.ID)
	return user.ID, nil
}

// Добавьте эти методы в ваш auth сервис

// UpdateProfile обновляет профиль пользователя
func (s *Service) UpdateProfile(ctx context.Context, userID string, updateData domain.UpdateProfileData) (*domain.User, error) {
	// 1. Получаем текущего пользователя
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, interfaces.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	// 2. Обновляем только те поля, которые переданы
	if updateData.FullName != nil {
		user.FullName = *updateData.FullName
	}
	if updateData.Phone != nil {
		user.Phone = *updateData.Phone
	}
	if updateData.Email != nil {
		user.Email = *updateData.Email
	}
	if updateData.BirthDate != nil {
		user.BirthDate = updateData.BirthDate
	}

	// 3. Сохраняем изменения
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		if errors.Is(err, interfaces.ErrDuplicatePhone) {
			return nil, ErrPhoneAlreadyRegistered
		}
		if errors.Is(err, interfaces.ErrDuplicateEmail) {
			return nil, ErrEmailAlreadyRegistered
		}
		return nil, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

// UpdateField обновляет конкретное поле пользователя
func (s *Service) UpdateField(ctx context.Context, userID string, field string, value string) error {
	// Валидация поля
	switch field {
	case "phone":
		if value == "" {
			return ErrValidationFailed
		}
		// Дополнительная валидация телефона
		if len(value) < 10 {
			return ErrValidationFailed
		}
	case "email":
		if value == "" {
			return ErrValidationFailed
		}
		// Простая валидация email
		if !strings.Contains(value, "@") {
			return ErrValidationFailed
		}
	case "full_name":
		if value == "" {
			return ErrValidationFailed
		}
	case "birth_date":
		// Валидация даты рождения (опционально)
		if value != "" {
			if _, err := time.Parse("2006-01-02", value); err != nil {
				return ErrValidationFailed
			}
		}
	default:
		return fmt.Errorf("unknown field: %s", field)
	}

	err := s.userRepo.UpdateField(ctx, userID, field, value)
	if err != nil {
		if errors.Is(err, interfaces.ErrDuplicatePhone) {
			return ErrPhoneAlreadyRegistered
		}
		if errors.Is(err, interfaces.ErrDuplicateEmail) {
			return ErrEmailAlreadyRegistered
		}
		if errors.Is(err, interfaces.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("update field: %w", err)
	}

	return nil
}
