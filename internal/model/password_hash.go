package model

import (
	"strings"

	errs "auth-microservice/internal/errors"
)

// PasswordHash представляет хеш пароля (bcrypt)
type PasswordHash string

// NewPasswordHash создаёт новый PasswordHash с валидацией формата bcrypt
// Ожидает уже готовый хеш, не хеширует пароль
func NewPasswordHash(hash string) (*PasswordHash, error) {
	if hash == "" {
		return nil, errs.ErrPasswordInvalid
	}

	// bcrypt хеши начинаются с $2a$, $2b$ или $2y$
	if !strings.HasPrefix(hash, "$2a$") &&
		!strings.HasPrefix(hash, "$2b$") &&
		!strings.HasPrefix(hash, "$2y$") {
		return nil, errs.ErrPasswordInvalid
	}

	// bcrypt хеш имеет длину около 60 символов (минимум 53, максимум 72)
	if len(hash) < 53 || len(hash) > 72 {
		return nil, errs.ErrPasswordInvalid
	}

	passwordHash := PasswordHash(hash)
	return &passwordHash, nil
}

// String возвращает строковое представление (для логгирования)
func (p PasswordHash) String() string {
	return "[REDACTED]"
}

// Value возвращает значение хеша (для сравнения и сохранения)
func (p PasswordHash) Value() string {
	return string(p)
}

// Equal сравнивает два хеша
func (p PasswordHash) Equal(other *PasswordHash) bool {
	if other == nil {
		return false
	}
	return p == *other
}

// NewPasswordHashFromString создаёт PasswordHash из строки без валидации.
// Используется при загрузке из БД (валидация уже выполнена при сохранении).
func NewPasswordHashFromString(hash string) *PasswordHash {
	passwordHash := PasswordHash(hash)
	return &passwordHash
}
