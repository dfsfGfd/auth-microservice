package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Ошибки кэша
var (
	ErrTokenNotFound = errors.New("token not found")
)

// RedisCache реализация TokenCache на Redis.
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache создаёт новый Redis кэш для токенов.
func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
	if prefix == "" {
		prefix = "refresh:"
	}
	return &RedisCache{
		client: client,
		prefix: prefix,
	}
}

// key возвращает полный ключ для Redis.
func (c *RedisCache) key(token string) string {
	return c.prefix + token
}

// Set сохраняет refresh токен с привязкой к account ID.
func (c *RedisCache) Set(ctx context.Context, token string, accountID string, ttl time.Duration) error {
	key := c.key(token)
	return c.client.Set(ctx, key, accountID, ttl).Err()
}

// Get получает account ID по refresh токену.
func (c *RedisCache) Get(ctx context.Context, token string) (string, error) {
	key := c.key(token)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrTokenNotFound
		}
		return "", fmt.Errorf("redis get: %w", err)
	}
	return val, nil
}

// Delete удаляет refresh токен из кэша.
func (c *RedisCache) Delete(ctx context.Context, token string) error {
	key := c.key(token)
	_, err := c.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	return nil
}

// Exists проверяет существование refresh токена.
func (c *RedisCache) Exists(ctx context.Context, token string) (bool, error) {
	key := c.key(token)
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists: %w", err)
	}
	return result > 0, nil
}
