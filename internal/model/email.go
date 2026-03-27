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
	// Нормализуем: убираем пробелы по краям
	value = strings.TrimSpace(value)

	if value == "" {
		return nil, errs.ErrEmailInvalid
	}

	if len(value) > EmailMaxLength {
		return nil, errs.ErrEmailTooLong
	}

	// mail.ParseAddress парсит "Name <email>" и просто "email"
	if _, err := mail.ParseAddress(value); err != nil {
		return nil, errs.ErrEmailInvalid
	}

	email := Email(value)
	return &email, nil
}

// Value возвращает значение email
func (e Email) Value() string {
	return string(e)
}

// String реализует fmt.Stringer
func (e Email) String() string {
	return e.Value()
}

// NewEmailFromDB создаёт Email из БД без валидации.
// Используется при загрузке из БД, где данные уже прошли валидацию.
// НЕ используйте для новых данных — используйте NewEmail().
func NewEmailFromDB(value string) *Email {
	email := Email(value)
	return &email
}
