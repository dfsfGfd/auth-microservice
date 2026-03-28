// Package middleware предоставляет HTTP и gRPC middleware для приложения.
package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiterConfig конфигурация rate limiter
type RateLimiterConfig struct {
	// Window размер окна времени (например, 1*time.Minute)
	Window time.Duration
	// Limit максимальное количество запросов в окно
	Limit int
	// Prefix префикс для ключей в Redis
	Prefix string
}

// Validate валидирует конфигурацию
func (c *RateLimiterConfig) Validate() error {
	if c.Window <= 0 {
		return fmt.Errorf("window must be positive")
	}
	if c.Limit <= 0 {
		return fmt.Errorf("limit must be positive")
	}
	if c.Prefix == "" {
		c.Prefix = "ratelimit:"
	}
	return nil
}

// RateLimiter Redis-based rate limiter с использованием sliding window
type RateLimiter struct {
	client  *redis.Client
	configs map[string]RateLimiterConfig
}

// NewRateLimiter создаёт новый rate limiter с конфигурациями для разных endpoint'ов
func NewRateLimiter(client *redis.Client, configs map[string]RateLimiterConfig) (*RateLimiter, error) {
	// Валидируем все конфигурации
	for name, cfg := range configs {
		if err := cfg.Validate(); err != nil {
			return nil, fmt.Errorf("invalid config for %s: %w", name, err)
		}
	}

	return &RateLimiter{
		client:  client,
		configs: configs,
	}, nil
}

// Allow проверяет, разрешён ли запрос для данного ключа и endpoint
// Возвращает (allowed, remaining, resetTime, error)
func (rl *RateLimiter) Allow(ctx context.Context, endpoint, key string) (bool, int, time.Time, error) {
	config, ok := rl.configs[endpoint]
	if !ok {
		// Если конфигурации нет, пропускаем запрос
		return true, 0, time.Now(), nil
	}

	// Добавляем таймаут на Redis операцию для защиты от зависаний
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	now := time.Now()
	windowStart := now.Add(-config.Window)
	windowKey := config.Prefix + key

	// Используем Redis pipeline для атомарности
	pipe := rl.client.Pipeline()

	// Удаляем старые записи за пределами окна
	pipe.ZRemRangeByScore(ctx, windowKey, "0", strconv.FormatInt(windowStart.UnixNano(), 10))

	// Добавляем текущий запрос
	pipe.ZAdd(ctx, windowKey, redis.Z{
		Score:  float64(now.UnixNano()),
		Member: strconv.FormatInt(now.UnixNano(), 10) + ":" + strconv.FormatInt(time.Now().UnixNano(), 10),
	})

	// Устанавливаем TTL для ключа (автоочистка)
	pipe.Expire(ctx, windowKey, config.Window*2)

	// Считаем количество запросов в окне
	countCmd := pipe.ZCard(ctx, windowKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		// Для критичных endpoint'ов (login, register) используем fail-close
		// Это предотвращает brute-force атаки при недоступности Redis
		if endpoint == "login" || endpoint == "register" {
			return false, 0, now, fmt.Errorf("rate limiter unavailable: %w", err)
		}
		// Для остальных endpoint'ов — fail-open (пропускаем запрос)
		return true, 0, now, nil
	}

	count := int(countCmd.Val())
	remaining := config.Limit - count
	if remaining < 0 {
		remaining = 0
	}

	// Время сброса лимита
	resetTime := now.Add(config.Window)

	if count > config.Limit {
		return false, remaining, resetTime, nil
	}

	return true, remaining, resetTime, nil
}
