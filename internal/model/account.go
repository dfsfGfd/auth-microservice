package model

import (
	"time"
)

// Account — агрегат аккаунта.
//
// Инварианты:
//   - ID всегда присутствует
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
// Принимает готовые валидированные Value Objects.
// Хеширование пароля выполняется в сервисном слое.
func NewAccount(email *Email, passwordHash *PasswordHash) (*Account, error) {
	now := time.Now()
	return &Account{
		id:           0, // ID устанавливается через SetID после сохранения в БД
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

// SetID устанавливает идентификатор аккаунта (для конвертеров из БД).
func (a *Account) SetID(id int64) {
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
