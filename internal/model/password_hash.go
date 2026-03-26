package model

import errs "auth-microservice/internal/errors"

// PasswordHash представляет хеш пароля (bcrypt)
type PasswordHash string

// NewPasswordHash создаёт новый PasswordHash с базовой валидацией длины
func NewPasswordHash(hash string) (*PasswordHash, error) {
	if hash == "" {
		return nil, errs.ErrPasswordInvalid
	}
	// bcrypt хеш: минимум 53, максимум 72 символа
	if len(hash) < 53 || len(hash) > 72 {
		return nil, errs.ErrPasswordInvalid
	}
	h := PasswordHash(hash)
	return &h, nil
}

// NewPasswordHashFromString создаёт PasswordHash без валидации (для чтения из БД)
func NewPasswordHashFromString(hash string) *PasswordHash {
	h := PasswordHash(hash)
	return &h
}

// Value возвращает хеш для сравнения и сохранения
func (p PasswordHash) Value() string {
	return string(p)
}

// String реализует fmt.Stringer (безопасное логирование)
func (p PasswordHash) String() string {
	return "[REDACTED]"
}

// Equal сравнивает два хеша
func (p PasswordHash) Equal(other *PasswordHash) bool {
	if other == nil {
		return false
	}
	return p == *other
}
