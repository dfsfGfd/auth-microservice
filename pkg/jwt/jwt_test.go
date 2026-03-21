package jwt_test

import (
	"testing"
	"time"

	"auth-microservice/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	t.Run("успешное создание сервиса", func(t *testing.T) {
		config := jwt.Config{
			SecretKey:       "test-secret-key",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 14 * 24 * time.Hour,
			Issuer:          "auth-service",
		}

		service, err := jwt.NewService(config)

		require.NoError(t, err)
		require.NotNil(t, service)
	})

	t.Run("ошибка при пустом SecretKey", func(t *testing.T) {
		config := jwt.Config{
			SecretKey:       "",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 14 * 24 * time.Hour,
			Issuer:          "auth-service",
		}

		service, err := jwt.NewService(config)

		assert.Error(t, err)
		assert.Nil(t, service)
	})

	t.Run("ошибка при нулевом AccessTokenTTL", func(t *testing.T) {
		config := jwt.Config{
			SecretKey:       "test-secret-key",
			AccessTokenTTL:  0,
			RefreshTokenTTL: 14 * 24 * time.Hour,
			Issuer:          "auth-service",
		}

		service, err := jwt.NewService(config)

		assert.Error(t, err)
		assert.Nil(t, service)
	})

	t.Run("ошибка при нулевом RefreshTokenTTL", func(t *testing.T) {
		config := jwt.Config{
			SecretKey:       "test-secret-key",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 0,
			Issuer:          "auth-service",
		}

		service, err := jwt.NewService(config)

		assert.Error(t, err)
		assert.Nil(t, service)
	})

	t.Run("ошибка при пустом Issuer", func(t *testing.T) {
		config := jwt.Config{
			SecretKey:       "test-secret-key",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 14 * 24 * time.Hour,
			Issuer:          "",
		}

		service, err := jwt.NewService(config)

		assert.Error(t, err)
		assert.Nil(t, service)
	})
}

func TestService_GenerateTokens(t *testing.T) {
	service := createTestService(t)

	t.Run("успешная генерация токенов", func(t *testing.T) {
		accountID := "550e8400-e29b-41d4-a716-446655440000"
		email := "user@example.com"

		tokens, err := service.GenerateTokens(accountID, email)

		require.NoError(t, err)
		require.NotNil(t, tokens)
		assert.NotEmpty(t, tokens.AccessToken)
		assert.NotEmpty(t, tokens.RefreshToken)
		assert.Equal(t, int64(900), tokens.ExpiresIn)
		assert.Equal(t, "Bearer", tokens.TokenType)
	})

	t.Run("токены имеют разные значения", func(t *testing.T) {
		accountID := "550e8400-e29b-41d4-a716-446655440000"
		email := "user@example.com"

		tokens, err := service.GenerateTokens(accountID, email)

		require.NoError(t, err)
		assert.NotEqual(t, tokens.AccessToken, tokens.RefreshToken)
	})
}

func TestService_ValidateToken(t *testing.T) {
	service := createTestService(t)

	t.Run("успешная валидация access токена", func(t *testing.T) {
		accountID := "550e8400-e29b-41d4-a716-446655440000"
		email := "user@example.com"

		tokens, err := service.GenerateTokens(accountID, email)
		require.NoError(t, err)

		claims, err := service.ValidateToken(tokens.AccessToken)

		require.NoError(t, err)
		require.NotNil(t, claims)
		assert.Equal(t, accountID, claims.AccountID)
		assert.Equal(t, email, claims.Email)
		assert.Equal(t, jwt.AccessToken, claims.Type)
		assert.Equal(t, "auth-service", claims.Issuer)
	})

	t.Run("успешная валидация refresh токена", func(t *testing.T) {
		accountID := "550e8400-e29b-41d4-a716-446655440000"
		email := "user@example.com"

		tokens, err := service.GenerateTokens(accountID, email)
		require.NoError(t, err)

		claims, err := service.ValidateToken(tokens.RefreshToken)

		require.NoError(t, err)
		require.NotNil(t, claims)
		assert.Equal(t, jwt.RefreshToken, claims.Type)
	})

	t.Run("ошибка при невалидном токене", func(t *testing.T) {
		claims, err := service.ValidateToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("ошибка при токене с другим issuer", func(t *testing.T) {
		// Создаём сервис с другим issuer
		wrongService, err := jwt.NewService(jwt.Config{
			SecretKey:       "test-secret-key",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 14 * 24 * time.Hour,
			Issuer:          "wrong-service",
		})
		require.NoError(t, err)

		accountID := "550e8400-e29b-41d4-a716-446655440000"
		tokens, err := service.GenerateTokens(accountID, "user@example.com")
		require.NoError(t, err)

		claims, err := wrongService.ValidateToken(tokens.AccessToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestService_ValidateAccessToken(t *testing.T) {
	service := createTestService(t)

	t.Run("успешная валидация access токена", func(t *testing.T) {
		tokens, err := service.GenerateTokens("account-id", "user@example.com")
		require.NoError(t, err)

		claims, err := service.ValidateAccessToken(tokens.AccessToken)

		require.NoError(t, err)
		assert.Equal(t, jwt.AccessToken, claims.Type)
	})

	t.Run("ошибка при refresh токене", func(t *testing.T) {
		tokens, err := service.GenerateTokens("account-id", "user@example.com")
		require.NoError(t, err)

		claims, err := service.ValidateAccessToken(tokens.RefreshToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestService_ValidateRefreshToken(t *testing.T) {
	service := createTestService(t)

	t.Run("успешная валидация refresh токена", func(t *testing.T) {
		tokens, err := service.GenerateTokens("account-id", "user@example.com")
		require.NoError(t, err)

		claims, err := service.ValidateRefreshToken(tokens.RefreshToken)

		require.NoError(t, err)
		assert.Equal(t, jwt.RefreshToken, claims.Type)
	})

	t.Run("ошибка при access токене", func(t *testing.T) {
		tokens, err := service.GenerateTokens("account-id", "user@example.com")
		require.NoError(t, err)

		claims, err := service.ValidateRefreshToken(tokens.AccessToken)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestService_GetConfig(t *testing.T) {
	config := jwt.Config{
		SecretKey:       "test-secret-key",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 14 * 24 * time.Hour,
		Issuer:          "auth-service",
	}

	service, err := jwt.NewService(config)
	require.NoError(t, err)

	returnedConfig := service.GetConfig()

	assert.Equal(t, config.SecretKey, returnedConfig.SecretKey)
	assert.Equal(t, config.AccessTokenTTL, returnedConfig.AccessTokenTTL)
	assert.Equal(t, config.RefreshTokenTTL, returnedConfig.RefreshTokenTTL)
	assert.Equal(t, config.Issuer, returnedConfig.Issuer)
}

func createTestService(t *testing.T) *jwt.Service {
	t.Helper()

	service, err := jwt.NewService(jwt.Config{
		SecretKey:       "test-secret-key",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 14 * 24 * time.Hour,
		Issuer:          "auth-service",
	})
	require.NoError(t, err)
	return service
}
