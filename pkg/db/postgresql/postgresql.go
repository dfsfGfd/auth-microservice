// Package postgresql предоставляет подключение к PostgreSQL.
//
// Пример использования:
//
//	cfg := postgresql.Config{
//	    DSN: "postgres://user:pass@localhost:5432/db?sslmode=disable",
//	}
//	pool, err := postgresql.NewPool(ctx, cfg)
package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config конфигурация PostgreSQL
type Config struct {
	// DSN строка подключения
	DSN string
	// MaxConns максимальное количество подключений в пуле
	MaxConns int32
	// MinConns минимальное количество подключений в пуле
	MinConns int32
	// ConnTimeout таймаут подключения
	ConnTimeout time.Duration
	// MaxConnLifetime максимальное время жизни подключения
	MaxConnLifetime time.Duration
	// MaxConnIdleTime максимальное время простоя подключения
	MaxConnIdleTime time.Duration
}

// Validate валидирует и устанавливает значения по умолчанию
func (c *Config) Validate() {
	if c.MaxConns <= 0 {
		c.MaxConns = 25
	}
	if c.MinConns <= 0 {
		c.MinConns = 0
	}
	if c.ConnTimeout <= 0 {
		c.ConnTimeout = 10 * time.Second
	}
	if c.MaxConnLifetime <= 0 {
		c.MaxConnLifetime = time.Hour
	}
	if c.MaxConnIdleTime <= 0 {
		c.MaxConnIdleTime = 30 * time.Minute
	}
}

// NewPool создаёт пул подключений к PostgreSQL
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	cfg.Validate()

	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.ConnConfig.ConnectTimeout = cfg.ConnTimeout

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Проверка подключения
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// NewPoolNoPing создаёт пул подключений без проверки подключения
func NewPoolNoPing(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	cfg.Validate()

	poolConfig, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.ConnConfig.ConnectTimeout = cfg.ConnTimeout

	return pgxpool.NewWithConfig(ctx, poolConfig)
}
