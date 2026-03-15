// Package model предоставляет модели данных для слоя хранения PostgreSQL.
package model

import (
	"time"

	"github.com/google/uuid"
)

// User модель пользователя для хранения в БД.
//
// Соответствует таблице users:
//
//	CREATE TABLE users (
//	    id              UUID PRIMARY KEY,
//	    email           VARCHAR(254) NOT NULL UNIQUE,
//	    username        VARCHAR(30) NOT NULL UNIQUE,
//	    password_hash   VARCHAR(72) NOT NULL,
//	    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
//	    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
//	);
type User struct {
	// ID идентификатор пользователя (PRIMARY KEY)
	ID uuid.UUID `db:"id"`

	// Email адрес электронной почты (UNIQUE, NOT NULL)
	Email string `db:"email"`

	// Username имя пользователя (UNIQUE, NOT NULL)
	Username string `db:"username"`

	// PasswordHash хеш пароля bcrypt (NOT NULL, max 72 символа)
	PasswordHash string `db:"password_hash"`

	// CreatedAt время создания (NOT NULL, DEFAULT NOW())
	CreatedAt time.Time `db:"created_at"`

	// UpdatedAt время последнего обновления (NOT NULL, DEFAULT NOW())
	UpdatedAt time.Time `db:"updated_at"`
}

// TableName возвращает имя таблицы для модели.
func (User) TableName() string {
	return "users"
}
