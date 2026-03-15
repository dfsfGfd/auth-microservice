// Package repository предоставляет интерфейсы и реализации репозиториев.
package repository

import (
	"context"
	"time"
)

// TokenRepository интерфейс для работы с токенами.
//
// Все методы принимают context.Context для управления временем жизни запроса.
// Ошибки возвращаются в соответствии с пакетом repository/errors.
type TokenRepository interface {
	// ==================== Write Operations ====================

	// Store сохраняет токен в хранилище.
	//
	// Параметры:
	//   - ctx контекст выполнения
	//   - token строка токена
	//   - userID идентификатор пользователя
	//   - ttl время жизни токена
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	Store(ctx context.Context, token string, userID string, ttl time.Duration) error

	// Delete удаляет токен из хранилища (отзыв токена).
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	Delete(ctx context.Context, token string) error

	// DeleteByUserID удаляет все токены пользователя.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	DeleteByUserID(ctx context.Context, userID string) error

	// Extend продлевает время жизни токена.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrTokenNotFound — токен не найден
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	Extend(ctx context.Context, token string, ttl time.Duration) error

	// ==================== Read Operations ====================

	// GetUserID получает идентификатор пользователя по токену.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrTokenNotFound — токен не найден
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	GetUserID(ctx context.Context, token string) (string, error)

	// GetToken получает информацию о токене.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrTokenNotFound — токен не найден
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	GetToken(ctx context.Context, token string) (*TokenInfo, error)

	// ==================== Check Operations ====================

	// Exists проверяет существование токена в хранилище.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	Exists(ctx context.Context, token string) (bool, error)

	// IsExpired проверяет истёк ли токен.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrTokenNotFound — токен не найден
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	IsExpired(ctx context.Context, token string) (bool, error)

	// ==================== Count Operations ====================

	// CountByUserID подсчитывает количество токенов пользователя.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	CountByUserID(ctx context.Context, userID string) (int64, error)
}

// TokenInfo информация о токене.
type TokenInfo struct {
	// Token строка токена
	Token string

	// UserID идентификатор пользователя
	UserID string

	// ExpiresAt время истечения токена
	ExpiresAt time.Time

	// CreatedAt время создания токена
	CreatedAt time.Time

	// Metadata дополнительные метаданные
	Metadata map[string]string
}

// IsExpired проверяет истёк ли токен.
func (t *TokenInfo) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
