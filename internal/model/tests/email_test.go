package tests

import (
	"strings"
	"testing"

	"auth-microservice/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEmail(t *testing.T) {
	t.Run("valid emails", func(t *testing.T) {
		validEmails := []string{
			"user@example.com",
			"user.name@example.com",
			"user+tag@example.com",
			"user@mail.example.com",
			"user_name@example.co.uk",
		}

		for _, emailStr := range validEmails {
			t.Run(emailStr, func(t *testing.T) {
				email, err := model.NewEmail(emailStr)

				assert.NoError(t, err)
				assert.NotNil(t, email)
				assert.Equal(t, emailStr, email.String())
			})
		}
	})

	t.Run("valid emails with spaces", func(t *testing.T) {
		// Email с пробелами по краям должен нормализоваться
		email, err := model.NewEmail("  user@example.com  ")

		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", email.String())
	})

	t.Run("invalid emails", func(t *testing.T) {
		tests := []struct {
			name  string
			value string
			err   string
		}{
			{"empty string", "", "invalid email"},
			{"no @ symbol", "userexample.com", "invalid email"},
			{"no local part", "@example.com", "invalid email"},
			{"no domain", "user@", "invalid email"},
			{"spaces in email", "user@exam ple.com", "invalid email"},
			{"double @", "user@@example.com", "invalid email"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				email, err := model.NewEmail(tt.value)

				assert.Error(t, err)
				assert.Nil(t, email)
				assert.Contains(t, err.Error(), tt.err)
			})
		}
	})

	t.Run("email too long", func(t *testing.T) {
		// Email длиной ровно EmailMaxLength символов — валиден
		// len("@example.com") = 12, поэтому 254-12 = 242 символа local part
		maxEmail := strings.Repeat("a", model.EmailMaxLength-12) + "@example.com"
		email, err := model.NewEmail(maxEmail)
		assert.NoError(t, err)
		assert.NotNil(t, email)
		assert.Equal(t, model.EmailMaxLength, len(maxEmail))

		// Email длиной больше EmailMaxLength — невалиден
		tooLongEmail := strings.Repeat("a", model.EmailMaxLength-11) + "@example.com"
		_, err = model.NewEmail(tooLongEmail)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too long")
	})
}

func TestEmail_Value(t *testing.T) {
	email, err := model.NewEmail("user@example.com")
	require.NoError(t, err)

	assert.Equal(t, "user@example.com", email.Value())
}

func TestEmail_String(t *testing.T) {
	email, err := model.NewEmail("user@example.com")
	require.NoError(t, err)

	assert.Equal(t, "user@example.com", email.String())
}

func TestNewEmailFromDB(t *testing.T) {
	t.Run("creates email without validation", func(t *testing.T) {
		// NewEmailFromDB не валидирует, принимает любой email
		email := model.NewEmailFromDB("user@example.com")

		assert.NotNil(t, email)
		assert.Equal(t, "user@example.com", email.String())
	})

	t.Run("accepts emails that would fail NewEmail", func(t *testing.T) {
		// Эти email не прошли бы валидацию в NewEmail,
		// но NewEmailFromDB принимает их (доверенные данные из БД)
		email := model.NewEmailFromDB("invalid-without-at")

		assert.NotNil(t, email)
		assert.Equal(t, "invalid-without-at", email.String())
	})

	t.Run("accepts empty string without error", func(t *testing.T) {
		// NewEmailFromDB не возвращает ошибку, даже если email пустой
		email := model.NewEmailFromDB("")

		assert.NotNil(t, email)
		assert.Equal(t, "", email.String())
	})

	t.Run("Value method works correctly", func(t *testing.T) {
		email := model.NewEmailFromDB("test@example.com")

		assert.Equal(t, "test@example.com", email.Value())
	})
}
