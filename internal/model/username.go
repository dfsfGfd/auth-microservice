package model

import (
	"strings"
	"unicode"

	errs "auth-microservice/internal/errors"
)

type Username string

func NewUsername(value string) (*Username, error) {
	if value == "" {
		return nil, errs.ErrUsernameInvalid
	}

	if len(value) < 3 {
		return nil, errs.ErrUsernameTooShort
	}

	if len(value) > 30 {
		return nil, errs.ErrUsernameTooLong
	}

	for _, r := range value {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return nil, errs.ErrUsernameInvalid
		}
	}

	if strings.HasPrefix(value, "_") || strings.HasSuffix(value, "_") {
		return nil, errs.ErrUsernameInvalid
	}

	username := Username(value)
	return &username, nil
}

func (u Username) String() string {
	return string(u)
}

// Value возвращает значение username.
func (u Username) Value() string {
	return string(u)
}

func (u Username) Equal(other *Username) bool {
	if other == nil {
		return false
	}
	return u == *other
}
