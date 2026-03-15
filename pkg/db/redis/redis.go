// Package redis предоставляет подключение к Redis.
//
// Пример использования:
//
//	cfg := redis.Config{
//	    Addr: "localhost:6379",
//	}
//	client, err := redis.NewClient(cfg)
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config конфигурация Redis клиента
type Config struct {
	// Addr адрес Redis сервера (host:port)
	Addr string
	// Password пароль для аутентификации (опционально)
	Password string
	// DB номер базы данных (0-15)
	DB int
	// PoolSize размер пула подключений
	PoolSize int
	// MinIdleConns минимальное количество idle подключений
	MinIdleConns int
	// ConnTimeout таймаут подключения
	ConnTimeout time.Duration
	// ReadTimeout таймаут чтения
	ReadTimeout time.Duration
	// WriteTimeout таймаут записи
	WriteTimeout time.Duration
}

// Validate валидирует и устанавливает значения по умолчанию
func (c *Config) Validate() {
	if c.PoolSize <= 0 {
		c.PoolSize = 10
	}
	if c.MinIdleConns <= 0 {
		c.MinIdleConns = 0
	}
	if c.ConnTimeout <= 0 {
		c.ConnTimeout = 5 * time.Second
	}
	if c.ReadTimeout <= 0 {
		c.ReadTimeout = 3 * time.Second
	}
	if c.WriteTimeout <= 0 {
		c.WriteTimeout = 3 * time.Second
	}
}

// NewClient создаёт новый Redis клиент
func NewClient(cfg Config) (*redis.Client, error) {
	cfg.Validate()

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.ConnTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// Проверка подключения
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

// NewClientNoPing создаёт Redis клиент без проверки подключения
func NewClientNoPing(cfg Config) *redis.Client {
	cfg.Validate()

	return redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.ConnTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
}
