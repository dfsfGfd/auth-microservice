// Package errors предоставляет ошибки репозитория.
package errors

import "errors"

// Ошибки репозитория
var (
	// ErrNotFound агрегат не найден
	ErrNotFound = errors.New("repository: not found")

	// ErrAlreadyExists агрегат уже существует
	ErrAlreadyExists = errors.New("repository: already exists")

	// ErrUniqueViolation нарушение уникальности
	ErrUniqueViolation = errors.New("repository: unique violation")

	// ErrForeignKeyViolation нарушение внешнего ключа
	ErrForeignKeyViolation = errors.New("repository: foreign key violation")

	// ErrTransactionActive активная транзакция
	ErrTransactionActive = errors.New("repository: transaction active")

	// ErrNoTransaction нет активной транзакции
	ErrNoTransaction = errors.New("repository: no transaction")
)

// Ошибки пользователя
var (
	// ErrUserNotFound пользователь не найден
	ErrUserNotFound = errors.New("repository: user not found")

	// ErrUserAlreadyExists пользователь уже существует
	ErrUserAlreadyExists = errors.New("repository: user already exists")

	// ErrEmailAlreadyExists email уже зарегистрирован
	ErrEmailAlreadyExists = errors.New("repository: email already exists")

	// ErrUsernameAlreadyExists username уже занят
	ErrUsernameAlreadyExists = errors.New("repository: username already exists")
)

// Ошибки токена
var (
	// ErrTokenNotFound токен не найден
	ErrTokenNotFound = errors.New("repository: token not found")

	// ErrTokenExpired токен истёк
	ErrTokenExpired = errors.New("repository: token expired")
)
