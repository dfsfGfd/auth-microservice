package errors

import "errors"

var (
	// Ошибки аутентификации
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrTokenInvalid         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")

	// Ошибки аккаунта
	ErrAccountNotFound    = errors.New("account not found")
	ErrAccountExists      = errors.New("account already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")

	// Ошибки пароля
	ErrPasswordInvalid  = errors.New("invalid password")
	ErrPasswordTooShort = errors.New("password too short")

	// Ошибки email
	ErrEmailInvalid = errors.New("invalid email")
	ErrEmailTooLong = errors.New("email too long")
)
