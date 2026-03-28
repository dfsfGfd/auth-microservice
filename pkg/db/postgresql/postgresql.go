package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config конфигурация PostgreSQL
type Config struct {
	DSN         string
	MaxConns    int32
	ConnTimeout time.Duration
}

// NewPool создаёт пул подключений к PostgreSQL
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	// парсинг конфигурации
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)

	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn: %w", err)
	}
	// применение настроек
	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.ConnConfig.ConnectTimeout = cfg.ConnTimeout

	// создание нового пулла
	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init pool: %w", err)
	}

	// Проверка подключения
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
