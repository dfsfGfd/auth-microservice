package errors

import "errors"

var (
	// Общие ошибки
	ErrInternal        = errors.New("internal error")
	ErrNotFound        = errors.New("not found")
	ErrAlreadyExists   = errors.New("already exists")
	ErrInvalidArgument = errors.New("invalid argument")

	// Ошибки аутентификации
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrTokenInvalid      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrSessionNotFound   = errors.New("session not found")
	ErrSessionRevoked    = errors.New("session revoked")

	// Ошибки пользователя
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Ошибки пароля
	ErrPasswordInvalid     = errors.New("invalid password")
	ErrPasswordTooShort    = errors.New("password too short")

	// Ошибки email
	ErrEmailInvalid = errors.New("invalid email")
	ErrEmailTooLong = errors.New("email too long")

	// Ошибки username
	ErrUsernameInvalid   = errors.New("invalid username")
	ErrUsernameTooShort  = errors.New("username too short")
	ErrUsernameTooLong   = errors.New("username too long")

	// Ошибки репозитория
	ErrRepository        = errors.New("repository error")
	ErrDBConnection      = errors.New("database connection error")
	ErrDBQuery           = errors.New("database query error")
)
