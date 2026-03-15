package model

import (
	"time"

	"github.com/google/uuid"

	errs "auth-microservice/internal/errors"
)

// User — агрегат пользователя.
//
// Инварианты:
//   - ID всегда присутствует
//   - Email, Username, PasswordHash не могут быть nil после создания
//   - UpdatedAt >= CreatedAt
type User struct {
	id           uuid.UUID
	email        *Email
	username     *Username
	passwordHash *PasswordHash
	createdAt    time.Time
	updatedAt    time.Time
}

// NewUser создаёт нового пользователя.
// Принимает готовые валидированные Value Objects.
// Хеширование пароля выполняется в сервисном слое.
func NewUser(email *Email, username *Username, passwordHash *PasswordHash) (*User, error) {
	if email == nil {
		return nil, errs.ErrEmailInvalid
	}
	if username == nil {
		return nil, errs.ErrUsernameInvalid
	}
	if passwordHash == nil {
		return nil, errs.ErrPasswordInvalid
	}

	now := time.Now()
	return &User{
		id:           uuid.New(),
		email:        email,
		username:     username,
		passwordHash: passwordHash,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// ID возвращает идентификатор пользователя.
func (u *User) ID() uuid.UUID {
	return u.id
}

// Email возвращает email.
func (u *User) Email() *Email {
	return u.email
}

// Username возвращает username.
func (u *User) Username() *Username {
	return u.username
}

// PasswordHash возвращает хеш пароля.
func (u *User) PasswordHash() *PasswordHash {
	return u.passwordHash
}

// CreatedAt возвращает время создания.
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt возвращает время последнего обновления.
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// UpdateEmail обновляет email.
func (u *User) UpdateEmail(email *Email) error {
	if email == nil {
		return errs.ErrEmailInvalid
	}
	u.email = email
	u.updatedAt = time.Now()
	return nil
}

// UpdateUsername обновляет username.
func (u *User) UpdateUsername(username *Username) error {
	if username == nil {
		return errs.ErrUsernameInvalid
	}
	u.username = username
	u.updatedAt = time.Now()
	return nil
}

// UpdatePasswordHash обновляет хеш пароля.
func (u *User) UpdatePasswordHash(passwordHash *PasswordHash) error {
	if passwordHash == nil {
		return errs.ErrPasswordInvalid
	}
	u.passwordHash = passwordHash
	u.updatedAt = time.Now()
	return nil
}
