package auth

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrPhoneAlreadyRegistered = errors.New("phone already registered")
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrValidationFailed       = errors.New("validation failed")
	ErrInvalidCredentials     = errors.New("invalid credentials")
)
