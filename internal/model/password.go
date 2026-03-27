package model

import (
	errs "auth-microservice/internal/errors"
)

// PlainPassword представляет обычный (нехешированный) пароль
type PlainPassword string

// NewPlainPassword создаёт новый PlainPassword с валидацией правил
// Требования:
//   - Минимум 8 символов
//   - Без ограничений на регистр и наличие цифр
func NewPlainPassword(value string) (*PlainPassword, error) {
	if value == "" {
		return nil, errs.ErrPasswordInvalid
	}

	if len(value) < 8 {
		return nil, errs.ErrPasswordTooShort
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
