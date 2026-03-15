// Package repository предоставляет интерфейсы и реализации репозиториев.
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"auth-microservice/internal/model"
)

// UserRepository интерфейс для работы с пользователями.
//
// Все методы принимают context.Context для управления временем жизни запроса.
// Ошибки возвращаются в соответствии с пакетом repository/errors.
type UserRepository interface {
	// ==================== Write Operations ====================

	// Create создаёт нового пользователя в хранилище.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrUserAlreadyExists — пользователь уже существует
	//   - repository.ErrEmailAlreadyExists — email уже зарегистрирован
	//   - repository.ErrUsernameAlreadyExists — username уже занят
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	Create(ctx context.Context, user *model.User) error

	// Update обновляет существующего пользователя.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrUserNotFound — пользователь не найден
	//   - repository.ErrEmailAlreadyExists — новый email уже занят
	//   - repository.ErrUsernameAlreadyExists — новый username уже занят
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено времени ожидания
	Update(ctx context.Context, user *model.User) error

	// Delete удаляет пользователя по идентификатору.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrUserNotFound — пользователь не найден
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteByEmail удаляет пользователя по email.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrUserNotFound — пользователь не найден
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	DeleteByEmail(ctx context.Context, email string) error

	// DeleteByUsername удаляет пользователя по username.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrUserNotFound — пользователь не найден
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	DeleteByUsername(ctx context.Context, username string) error

	// ==================== Read Operations ====================

	// GetByID получает пользователя по идентификатору.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrUserNotFound — пользователь не найден
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// GetByEmail получает пользователя по email.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrUserNotFound — пользователь не найден
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// GetByUsername получает пользователя по username.
	//
	// Возвращаемые ошибки:
	//   - repository.ErrUserNotFound — пользователь не найден
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	GetByUsername(ctx context.Context, username string) (*model.User, error)

	// Find находит пользователей по спецификации.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	Find(ctx context.Context, spec UserSpec) ([]*model.User, error)

	// List возвращает список пользователей с пагинацией.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	List(ctx context.Context, limit, offset int) ([]*model.User, error)

	// ==================== Check Operations ====================

	// Exists проверяет существование пользователя по идентификатору.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	Exists(ctx context.Context, id uuid.UUID) (bool, error)

	// ExistsByEmail проверяет существование пользователя по email.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	ExistsByEmail(ctx context.Context, email string) (bool, error)

	// ExistsByUsername проверяет существование пользователя по username.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	// ==================== Count Operations ====================

	// Count подсчитывает общее количество пользователей.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	Count(ctx context.Context) (int64, error)

	// CountByEmail подсчитывает количество пользователей с указанным email.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	CountByEmail(ctx context.Context, email string) (int64, error)

	// CountByUsername подсчитывает количество пользователей с указанным username.
	//
	// Возвращаемые ошибки:
	//   - context.Canceled — контекст отменён
	//   - context.DeadlineExceeded — превышено время ожидания
	CountByUsername(ctx context.Context, username string) (int64, error)
}

// UserSpec спецификация для поиска пользователей.
//
// Пример использования:
//
//	spec := repository.UserSpec{
//	    Email:    ptr.String("user@example.com"),
//	    Limit:    10,
//	    Offset:   0,
//	    OrderBy:  "created_at",
//	    OrderDir: "DESC",
//	}
//	users, err := repo.Find(ctx, spec)
type UserSpec struct {
	// Email для фильтрации по email (опционально)
	Email *string

	// Username для фильтрации по username (опционально)
	Username *string

	// IDs для фильтрации по списку ID (опционально)
	IDs []uuid.UUID

	// CreatedAfter фильтрация по дате создания (опционально)
	CreatedAfter *time.Time

	// CreatedBefore фильтрация по дате создания (опционально)
	CreatedBefore *time.Time

	// Limit максимальное количество результатов (по умолчанию 100)
	Limit int

	// Offset смещение результатов (по умолчанию 0)
	Offset int

	// OrderBy поле для сортировки (по умолчанию "created_at")
	OrderBy string

	// OrderDir направление сортировки: "ASC" или "DESC" (по умолчанию "DESC")
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

// Validate валидирует спецификацию.
func (s *UserSpec) Validate() error {
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
	return nil
}
