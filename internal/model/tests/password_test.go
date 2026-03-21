package tests

import (
	"testing"

	"auth-microservice/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPlainPassword(t *testing.T) {
	t.Run("valid passwords", func(t *testing.T) {
		validPasswords := []string{
			"Password1",
			"Secure123",
			"MyPass123",
			"Test1234",
			"Abcdefg1",
		}

		for _, passwordStr := range validPasswords {
			t.Run(passwordStr, func(t *testing.T) {
				password, err := model.NewPlainPassword(passwordStr)

				assert.NoError(t, err)
				assert.NotNil(t, password)
				assert.Equal(t, "****", password.String())
			})
		}
	})

	t.Run("invalid passwords", func(t *testing.T) {
		tests := []struct {
			name  string
			value string
			err   string
		}{
			{"empty string", "", "invalid password"},
			{"too short (4 chars)", "Pass1", "password too short"},
			{"too short (7 chars)", "Pass123", "password too short"},
			{"too short (1 char)", "a", "password too short"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				password, err := model.NewPlainPassword(tt.value)

				assert.Error(t, err)
				assert.Nil(t, password)
				assert.Contains(t, err.Error(), tt.err)
			})
		}
	})

	t.Run("valid passwords - no uppercase required", func(t *testing.T) {
		validPasswords := []string{
			"password1",
			"qwerty123",
			"abcdefgh",
		}

		for _, passwordStr := range validPasswords {
			t.Run(passwordStr, func(t *testing.T) {
				password, err := model.NewPlainPassword(passwordStr)

				assert.NoError(t, err)
				assert.NotNil(t, password)
			})
		}
	})

	t.Run("valid passwords - no lowercase required", func(t *testing.T) {
		validPasswords := []string{
			"PASSWORD1",
			"QWERTY123",
			"ABCDEFGH",
		}

		for _, passwordStr := range validPasswords {
			t.Run(passwordStr, func(t *testing.T) {
				password, err := model.NewPlainPassword(passwordStr)

				assert.NoError(t, err)
				assert.NotNil(t, password)
			})
		}
	})

	t.Run("valid passwords - no digit required", func(t *testing.T) {
		validPasswords := []string{
			"Password",
			"qwertyui",
			"abcdefgh",
		}

		for _, passwordStr := range validPasswords {
			t.Run(passwordStr, func(t *testing.T) {
				password, err := model.NewPlainPassword(passwordStr)

				assert.NoError(t, err)
				assert.NotNil(t, password)
			})
		}
	})
}

func TestPlainPassword_Equal(t *testing.T) {
	t.Run("equal passwords", func(t *testing.T) {
		password1, err := model.NewPlainPassword("Password1")
		require.NoError(t, err)

		password2, err := model.NewPlainPassword("Password1")
		require.NoError(t, err)

		assert.True(t, password1.Equal(password2))
	})

	t.Run("different passwords", func(t *testing.T) {
		password1, err := model.NewPlainPassword("Password1")
		require.NoError(t, err)

		password2, err := model.NewPlainPassword("Secure123")
		require.NoError(t, err)

		assert.False(t, password1.Equal(password2))
	})

	t.Run("nil comparison", func(t *testing.T) {
		password, err := model.NewPlainPassword("Password1")
		require.NoError(t, err)

		assert.False(t, password.Equal(nil))
	})
}

func TestPlainPassword_Value(t *testing.T) {
	password, err := model.NewPlainPassword("SuperSecret123")
	require.NoError(t, err)

	assert.Equal(t, "SuperSecret123", password.Value())
}

func TestPlainPassword_String_MasksValue(t *testing.T) {
	password, err := model.NewPlainPassword("SuperSecret123")
	require.NoError(t, err)

	assert.Equal(t, "****", password.String())
	assert.NotContains(t, password.String(), "SuperSecret123")
}
