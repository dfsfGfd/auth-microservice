// Package repository предоставляет интерфейсы для работы с хранилищем данных.
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"auth-microservice/internal/model"
)

// UserRepository интерфейс для работы с пользователями в PostgreSQL.
type UserRepository interface {
	// ==================== Write Methods ====================

	// Create создаёт нового пользователя в БД.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	//   - postgres unique violation (email/username уже существует)
	Create(ctx context.Context, user *model.User) error

	// Update обновляет существующего пользователя в БД.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	//   - ErrNotFound (пользователь не найден)
	Update(ctx context.Context, user *model.User) error

	// Delete удаляет пользователя по ID.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByEmail удаляет пользователя по email.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	DeleteByEmail(ctx context.Context, email string) error

	// DeleteByUsername удаляет пользователя по username.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	DeleteByUsername(ctx context.Context, username string) error

	// ==================== Read Methods ====================

	// GetByID получает пользователя по ID.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	//   - ErrNotFound (пользователь не найден)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// GetByEmail получает пользователя по email.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	//   - ErrNotFound (пользователь не найден)
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// GetByUsername получает пользователя по username.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	//   - ErrNotFound (пользователь не найден)
	GetByUsername(ctx context.Context, username string) (*model.User, error)

	// Find находит пользователей по спецификации.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	Find(ctx context.Context, spec UserSpec) ([]*model.User, error)

	// List возвращает список пользователей с пагинацией.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	List(ctx context.Context, limit, offset int) ([]*model.User, error)

	// ==================== Check Methods ====================

	// Exists проверяет существование пользователя по ID.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	Exists(ctx context.Context, id uuid.UUID) (bool, error)

	// ExistsByEmail проверяет существование пользователя по email.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername проверяет существование пользователя по username.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ==================== Count Methods ====================

	// Count подсчитывает общее количество пользователей.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	Count(ctx context.Context) (int64, error)

	// CountByEmail подсчитывает количество пользователей с данным email.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	CountByEmail(ctx context.Context, email string) (int64, error)

	// CountByUsername подсчитывает количество пользователей с данным username.
	//
	// Ошибки:
	//   - context.Canceled
	//   - context.DeadlineExceeded
	CountByUsername(ctx context.Context, username string) (int64, error)
}

// UserSpec спецификация для поиска пользователей.
type UserSpec struct {
	// Email для фильтрации (опционально)
	Email *string

	// Username для фильтрации (опционально)
	Username *string

	// IDs для фильтрации по списку ID (опционально)
	IDs []uuid.UUID

	// CreatedAfter фильтрация по дате создания после (опционально)
	CreatedAfter *time.Time

	// CreatedBefore фильтрация по дате создания до (опционально)
	CreatedBefore *time.Time

	// Limit максимальное количество результатов (default: 100)
	Limit int

	// Offset смещение результатов (default: 0)
	Offset int

	// OrderBy поле для сортировки (default: "created_at")
	OrderBy string

	// OrderDir направление сортировки: "ASC" или "DESC" (default: "DESC")
	OrderDir string
}

// DefaultUserSpec возвращает спецификацию по умолчанию.
func DefaultUserSpec() UserSpec {
	return UserSpec{
		Limit:    100,
		Offset:   0,
		OrderBy:  "created_at",
		OrderDir: "DESC",
	}
}

// Validate валидирует и устанавливает значения по умолчанию.
func (s *UserSpec) Validate() {
	if s.Limit <= 0 {
		s.Limit = 100
	}
	if s.Offset < 0 {
		s.Offset = 0
	}
	if s.OrderBy == "" {
		s.OrderBy = "created_at"
	}
	if s.OrderDir != "ASC" && s.OrderDir != "DESC" {
		s.OrderDir = "DESC"
	}
}
