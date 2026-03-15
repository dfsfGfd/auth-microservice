package tests

import (
	"testing"

	"auth-microservice/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPasswordHash(t *testing.T) {
	t.Run("valid bcrypt hashes", func(t *testing.T) {
		validHashes := []struct {
			name string
			hash string
		}{
			{"$2a$", "$2a$12$abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuvwx"},
			{"$2b$", "$2b$12$abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuvwx"},
			{"$2y$", "$2y$12$abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuvwx"},
		}

		for _, tt := range validHashes {
			t.Run(tt.name, func(t *testing.T) {
				passwordHash, err := model.NewPasswordHash(tt.hash)

				assert.NoError(t, err)
				assert.NotNil(t, passwordHash)
				assert.Equal(t, "[REDACTED]", passwordHash.String())
			})
		}
	})

	t.Run("invalid hashes", func(t *testing.T) {
		tests := []struct {
			name  string
			value string
			err   string
		}{
			{"empty string", "", "invalid password"},
			{"wrong prefix", "$1$12$abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuvwx", "invalid password"},
			{"too short", "$2a$12$abc", "invalid password"},
			{"too long", "$2a$12$abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz", "invalid password"},
			{"no prefix", "abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuvwx", "invalid password"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				passwordHash, err := model.NewPasswordHash(tt.value)

				assert.Error(t, err)
				assert.Nil(t, passwordHash)
				assert.Contains(t, err.Error(), tt.err)
			})
		}
	})
}

func TestPasswordHash_Equal(t *testing.T) {
	t.Run("equal hashes", func(t *testing.T) {
		hash1, err := model.NewPasswordHash("$2a$12$abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuvwx")
		require.NoError(t, err)

		hash2, err := model.NewPasswordHash("$2a$12$abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuvwx")
		require.NoError(t, err)

		assert.True(t, hash1.Equal(hash2))
	})

	t.Run("different hashes", func(t *testing.T) {
		hash1, err := model.NewPasswordHash("$2a$12$abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuvwx")
		require.NoError(t, err)

		hash2, err := model.NewPasswordHash("$2a$12$xyz123abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuv")
		require.NoError(t, err)

		assert.False(t, hash1.Equal(hash2))
	})

	t.Run("nil comparison", func(t *testing.T) {
		hash, err := model.NewPasswordHash("$2a$12$abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuvwx")
		require.NoError(t, err)

		assert.False(t, hash.Equal(nil))
	})
}
