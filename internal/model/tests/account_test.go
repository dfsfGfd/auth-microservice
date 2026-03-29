package tests

import (
	"testing"

	"auth-microservice/internal/errors"
	"auth-microservice/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	t.Run("valid account creation", func(t *testing.T) {
		email, _ := model.NewEmail("user@example.com")
		passwordHash := model.NewPasswordHashFromString("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")

		account, err := model.NewAccount(1, email, passwordHash)

		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, int64(1), account.ID())
		assert.Equal(t, "user@example.com", account.Email().String())
		assert.NotNil(t, account.PasswordHash())
		assert.NotNil(t, account.CreatedAt())
		assert.NotNil(t, account.UpdatedAt())
	})

	t.Run("invalid ID - zero", func(t *testing.T) {
		email, _ := model.NewEmail("user@example.com")
		passwordHash := model.NewPasswordHashFromString("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")

		account, err := model.NewAccount(0, email, passwordHash)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrAccountInvalidID)
		assert.Nil(t, account)
	})

	t.Run("invalid ID - negative", func(t *testing.T) {
		email, _ := model.NewEmail("user@example.com")
		passwordHash := model.NewPasswordHashFromString("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")

		account, err := model.NewAccount(-1, email, passwordHash)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrAccountInvalidID)
		assert.Nil(t, account)
	})

	t.Run("large valid ID", func(t *testing.T) {
		email, _ := model.NewEmail("user@example.com")
		passwordHash := model.NewPasswordHashFromString("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")

		// Snowflake ID из реального примера
		snowflakeID := int64(296494707175849985)
		account, err := model.NewAccount(snowflakeID, email, passwordHash)

		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, snowflakeID, account.ID())
	})
}
