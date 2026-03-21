// Package token предоставляет интерфейс и реализацию кэша для токенов.
package token

import (
	"context"
	"time"
)

// TokenCache интерфейс для кэширования токенов.
type TokenCache interface {
	// Set сохраняет refresh токен с привязкой к account ID.
	Set(ctx context.Context, token, accountID string, ttl time.Duration) error

	// Get получает account ID по refresh токену.
	Get(ctx context.Context, token string) (string, error)

	// Delete удаляет refresh токен из кэша.
	Delete(ctx context.Context, token string) error

	// Exists проверяет существование refresh токена.
	Exists(ctx context.Context, token string) (bool, error)
}
