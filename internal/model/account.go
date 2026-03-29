package model

import (
	"auth-microservice/internal/errors"
	"time"
)

// Account — агрегат аккаунта.
//
// Инварианты:
//   - ID всегда присутствует и > 0
//   - Email, PasswordHash не могут быть nil после создания
//   - UpdatedAt >= CreatedAt
type Account struct {
	id           int64
	email        *Email
	passwordHash *PasswordHash
	createdAt    time.Time
	updatedAt    time.Time
}

// NewAccount создаёт новый аккаунт.
// Принимает готовые валидированные Value Objects и ID.
// Хеширование пароля выполняется в сервисном слое.
func NewAccount(id int64, email *Email, passwordHash *PasswordHash) (*Account, error) {
	if id <= 0 {
		return nil, errors.ErrAccountInvalidID
	}

	now := time.Now()
	return &Account{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ID возвращает идентификатор аккаунта.
func (a *Account) ID() int64 {
	return a.id
}

// Email возвращает email.
func (a *Account) Email() *Email {
	return a.email
}

// PasswordHash возвращает хеш пароля.
func (a *Account) PasswordHash() *PasswordHash {
	return a.passwordHash
}

// CreatedAt возвращает время создания.
func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

// UpdatedAt возвращает время последнего обновления.
func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

// NewAccountFromDB создаёт Account из данных БД (internal API для repository).
// Не валидирует ID — предполагается, что данные из БД корректны.
func NewAccountFromDB(id int64, email *Email, passwordHash *PasswordHash, createdAt, updatedAt time.Time) *Account {
	return &Account{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}
