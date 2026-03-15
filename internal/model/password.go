package model

import (
	"unicode"

	errs "auth-microservice/internal/errors"
)

// PlainPassword представляет обычный (нехешированный) пароль
type PlainPassword string

// NewPlainPassword создаёт новый PlainPassword с валидацией правил
func NewPlainPassword(value string) (*PlainPassword, error) {
	if value == "" {
		return nil, errs.ErrPasswordInvalid
	}

	if len(value) < 8 {
		return nil, errs.ErrPasswordTooShort
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, r := range value {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return nil, errs.ErrPasswordInvalid
	}

	password := PlainPassword(value)
	return &password, nil
}

// String возвращает замаскированное представление (для логгирования)
func (p PlainPassword) String() string {
	return "****"
}

// Value возвращает оригинальное значение пароля (для хеширования)
func (p PlainPassword) Value() string {
	return string(p)
}

// Equal сравнивает два пароля
func (p PlainPassword) Equal(other *PlainPassword) bool {
	if other == nil {
		return false
	}
	return p == *other
}
