package model

import (
	"net/mail"
	"strings"

	errs "auth-microservice/internal/errors"
)

// EmailMaxLength максимальная длина email по RFC 5321
const EmailMaxLength = 254

// Email представляет email-адрес в домене.
//
// Инварианты:
//   - Непустое значение
//   - Длина не более EmailMaxLength символов
//   - Валидный формат email по RFC 5321
type Email string

// NewEmail создаёт новый Email с валидацией инвариантов домена
func NewEmail(value string) (*Email, error) {
	if value == "" {
		return nil, errs.ErrEmailInvalid
	}

	if len(value) > EmailMaxLength {
		return nil, errs.ErrEmailTooLong
	}

	// Нормализуем: убираем пробелы по краям
	value = strings.TrimSpace(value)

	// mail.ParseAddress парсит "Name <email>" и просто "email"
	if _, err := mail.ParseAddress(value); err != nil {
		return nil, errs.ErrEmailInvalid
	}

	email := Email(value)
	return &email, nil
}

// String возвращает строковое представление email
func (e Email) String() string {
	return string(e)
}

// Value возвращает значение email
func (e Email) Value() string {
	return string(e)
}

// Equal сравнивает два Email
func (e Email) Equal(other *Email) bool {
	if other == nil {
		return false
	}
	return e == *other
}
