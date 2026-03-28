package model

import (
	"time"

	"github.com/google/uuid"
)

// Account — агрегат аккаунта.
//
// Инварианты:
//   - ID всегда присутствует
//   - Email, PasswordHash не могут быть nil после создания
//   - UpdatedAt >= CreatedAt
type Account struct {
	id           uuid.UUID
	email        *Email
	passwordHash *PasswordHash
	createdAt    time.Time
	updatedAt    time.Time
}

// NewAccount создаёт новый аккаунт.
// Принимает готовые валидированные Value Objects.
// Хеширование пароля выполняется в сервисном слое.
func NewAccount(email *Email, passwordHash *PasswordHash) (*Account, error) {
	now := time.Now()
	return &Account{
		id:           uuid.New(),
		email:        email,
		passwordHash: passwordHash,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ID возвращает идентификатор аккаунта.
func (a *Account) ID() uuid.UUID {
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

// SetID устанавливает идентификатор аккаунта (для конвертеров из БД).
func (a *Account) SetID(id uuid.UUID) {
	a.id = id
}

// SetCreatedAt устанавливает время создания (для конвертеров из БД).
func (a *Account) SetCreatedAt(t time.Time) {
	a.createdAt = t
}

// SetUpdatedAt устанавливает время обновления (для конвертеров из БД).
func (a *Account) SetUpdatedAt(t time.Time) {
	a.updatedAt = t
}
