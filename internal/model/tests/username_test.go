package tests

import (
	"strings"
	"testing"

	"auth-microservice/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUsername(t *testing.T) {
	t.Run("valid usernames", func(t *testing.T) {
		validUsernames := []string{
			"user",
			"john_doe",
			"user123",
			"abc",
			"user_name_123",
		}

		for _, usernameStr := range validUsernames {
			t.Run(usernameStr, func(t *testing.T) {
				username, err := model.NewUsername(usernameStr)

				assert.NoError(t, err)
				assert.NotNil(t, username)
				assert.Equal(t, usernameStr, username.String())
			})
		}
	})

	t.Run("invalid usernames", func(t *testing.T) {
		tests := []struct {
			name  string
			value string
			err   string
		}{
			{"empty string", "", "invalid username"},
			{"too short (1 char)", "a", "username too short"},
			{"too short (2 chars)", "ab", "username too short"},
			{"starts with underscore", "_user", "invalid username"},
			{"ends with underscore", "user_", "invalid username"},
			{"special character", "user!", "invalid username"},
			{"space in username", "user name", "invalid username"},
			{"hyphen", "user-name", "invalid username"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				username, err := model.NewUsername(tt.value)

				assert.Error(t, err)
				assert.Nil(t, username)
				assert.Contains(t, err.Error(), tt.err)
			})
		}
	})

	t.Run("username too long", func(t *testing.T) {
		longUsername := strings.Repeat("a", 31)

		username, err := model.NewUsername(longUsername)

		assert.Error(t, err)
		assert.Nil(t, username)
		assert.Contains(t, err.Error(), "too long")
	})
}

func TestUsername_Equal(t *testing.T) {
	t.Run("equal usernames", func(t *testing.T) {
		username1, err := model.NewUsername("john_doe")
		require.NoError(t, err)

		username2, err := model.NewUsername("john_doe")
		require.NoError(t, err)

		assert.True(t, username1.Equal(username2))
	})

	t.Run("different usernames", func(t *testing.T) {
		username1, err := model.NewUsername("john_doe")
		require.NoError(t, err)

		username2, err := model.NewUsername("jane_doe")
		require.NoError(t, err)

		assert.False(t, username1.Equal(username2))
	})

	t.Run("nil comparison", func(t *testing.T) {
		username, err := model.NewUsername("john_doe")
		require.NoError(t, err)

		assert.False(t, username.Equal(nil))
	})
}
