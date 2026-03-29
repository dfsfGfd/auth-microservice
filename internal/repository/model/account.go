// Package model предоставляет модели данных для слоя хранения PostgreSQL.
package model

import (
	"time"
)

// Account модель аккаунта для хранения в БД.
//
// Соответствует таблице accounts:
//
//	CREATE TABLE accounts (
//	    id              BIGINT PRIMARY KEY,
//	    email           VARCHAR(254) NOT NULL UNIQUE,
//	    password        VARCHAR(72) NOT NULL,
//	    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
//	    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
//	);
type Account struct {
	// ID идентификатор аккаунта (PRIMARY KEY)
	ID int64 `db:"id"`

	// Email адрес электронной почты (UNIQUE, NOT NULL)
	Email string `db:"email"`

	// PasswordHash хеш пароля bcrypt (NOT NULL, max 72 символа)
	PasswordHash string `db:"password"`

	// CreatedAt время создания (NOT NULL, DEFAULT NOW())
	CreatedAt time.Time `db:"created_at"`

	// UpdatedAt время последнего обновления (NOT NULL, DEFAULT NOW())
	UpdatedAt time.Time `db:"updated_at"`
}
