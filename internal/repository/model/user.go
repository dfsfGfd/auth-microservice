// Package model предоставляет модели данных для слоя хранения.
//
// Эти модели используются для маппинга между domain моделями и таблицами БД.
package model

import (
	"time"

	"github.com/google/uuid"
)

// User модель пользователя для хранения в БД
type User struct {
	// ID идентификатор пользователя
	ID uuid.UUID `db:"id"`

	// Email адрес электронной почты
	Email string `db:"email"`

	// Username имя пользователя
	Username string `db:"username"`

	// PasswordHash хеш пароля
	PasswordHash string `db:"password_hash"`

	// CreatedAt время создания
	CreatedAt time.Time `db:"created_at"`

	// UpdatedAt время последнего обновления
	UpdatedAt time.Time `db:"updated_at"`
}

// TableName возвращает имя таблицы для модели
func (User) TableName() string {
	return "users"
}

// UserEmailIndex индекс для поиска по email
type UserEmailIndex struct {
	UserID uuid.UUID `db:"user_id"`
	Email  string    `db:"email"`
}

// UserUsernameIndex индекс для поиска по username
type UserUsernameIndex struct {
	UserID   uuid.UUID `db:"user_id"`
	Username string    `db:"username"`
}
