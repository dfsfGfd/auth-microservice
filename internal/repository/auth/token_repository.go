// Package auth предоставляет Redis реализацию репозитория для токенов.
package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"auth-microservice/internal/repository"
	repoerrors "auth-microservice/internal/repository/errors"
)

// ensure TokenRepository implements interface
var _ repository.TokenRepository = (*TokenRepository)(nil)

// TokenRepository Redis реализация repository.TokenRepository
type TokenRepository struct {
	client *redis.Client
	prefix string
}

// NewTokenRepository создаёт новый TokenRepository
func NewTokenRepository(client *redis.Client, prefix string) *TokenRepository {
	return &TokenRepository{
		client: client,
		prefix: prefix,
	}
}

// key формирует ключ с префиксом
func (r *TokenRepository) key(token string) string {
	return fmt.Sprintf("%s:refresh:%s", r.prefix, token)
}

// userKey формирует ключ для списка токенов пользователя
func (r *TokenRepository) userKey(userID string) string {
	return fmt.Sprintf("%s:user:%s:tokens", r.prefix, userID)
}

// ============================================================================
// Write Operations
// ============================================================================

// Store сохраняет токен в хранилище.
func (r *TokenRepository) Store(ctx context.Context, token string, userID string, ttl time.Duration) error {
	key := r.key(token)

	// Сохраняем токен
	if err := r.client.Set(ctx, key, userID, ttl).Err(); err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	// Добавляем токен в список токенов пользователя
	if err := r.client.SAdd(ctx, r.userKey(userID), token).Err(); err != nil {
		return fmt.Errorf("failed to add token to user set: %w", err)
	}

	// Устанавливаем TTL для списка токенов
	if err := r.client.Expire(ctx, r.userKey(userID), ttl).Err(); err != nil {
		return fmt.Errorf("failed to set TTL for user token set: %w", err)
	}

	return nil
}

// Delete удаляет токен из хранилища (отзыв токена).
func (r *TokenRepository) Delete(ctx context.Context, token string) error {
	key := r.key(token)

	// Получаем userID для удаления из списка
	userID, _ := r.client.Get(ctx, key).Result()

	// Удаляем токен
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	// Удаляем токен из списка токенов пользователя
	if userID != "" {
		if err := r.client.SRem(ctx, r.userKey(userID), token).Err(); err != nil {
			return fmt.Errorf("failed to remove token from user set: %w", err)
		}
	}

	return nil
}

// DeleteByUserID удаляет все токены пользователя.
func (r *TokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	// Получаем все токены пользователя
	tokens, err := r.client.SMembers(ctx, r.userKey(userID)).Result()
	if err != nil {
		return fmt.Errorf("failed to get user tokens: %w", err)
	}

	// Удаляем каждый токен
	for _, token := range tokens {
		if err := r.client.Del(ctx, r.key(token)).Err(); err != nil {
			return fmt.Errorf("failed to delete token: %w", err)
		}
	}

	// Удаляем список токенов пользователя
	if err := r.client.Del(ctx, r.userKey(userID)).Err(); err != nil {
		return fmt.Errorf("failed to delete user token set: %w", err)
	}

	return nil
}

// Extend продлевает время жизни токена.
func (r *TokenRepository) Extend(ctx context.Context, token string, ttl time.Duration) error {
	key := r.key(token)

	// Проверяем существование токена
	exists, err := r.Exists(ctx, token)
	if err != nil {
		return err
	}
	if !exists {
		return repoerrors.ErrTokenNotFound
	}

	// Продлеваем TTL
	if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
		return fmt.Errorf("failed to extend token TTL: %w", err)
	}

	return nil
}

// ============================================================================
// Read Operations
// ============================================================================

// GetUserID получает идентификатор пользователя по токену.
func (r *TokenRepository) GetUserID(ctx context.Context, token string) (string, error) {
	key := r.key(token)

	userID, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", repoerrors.ErrTokenNotFound
		}
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	return userID, nil
}

// GetToken получает информацию о токене.
func (r *TokenRepository) GetToken(ctx context.Context, token string) (*repository.TokenInfo, error) {
	key := r.key(token)

	// Получаем userID
	userID, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, repoerrors.ErrTokenNotFound
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	// Получаем TTL
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get token TTL: %w", err)
	}

	return &repository.TokenInfo{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

// ============================================================================
// Check Operations
// ============================================================================

// Exists проверяет существование токена в хранилище.
func (r *TokenRepository) Exists(ctx context.Context, token string) (bool, error) {
	key := r.key(token)

	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token: %w", err)
	}

	return result > 0, nil
}

// IsExpired проверяет истёк ли токен.
func (r *TokenRepository) IsExpired(ctx context.Context, token string) (bool, error) {
	key := r.key(token)

	// Получаем оставшееся время жизни
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return true, repoerrors.ErrTokenNotFound
		}
		return false, fmt.Errorf("failed to get token TTL: %w", err)
	}

	// Если TTL <= 0, токен истёк
	return ttl <= 0, nil
}

// ============================================================================
// Count Operations
// ============================================================================

// CountByUserID подсчитывает количество токенов пользователя.
func (r *TokenRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	count, err := r.client.SCard(ctx, r.userKey(userID)).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to count user tokens: %w", err)
	}

	return count, nil
}
