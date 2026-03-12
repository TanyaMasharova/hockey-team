package auth

import "errors"

var (
    ErrPhoneAlreadyRegistered = errors.New("phone already registered")
    ErrEmailAlreadyRegistered = errors.New("email already registered")
    ErrValidationFailed       = errors.New("validation failed")
    ErrInvalidCredentials     = errors.New("invalid credentials")
)