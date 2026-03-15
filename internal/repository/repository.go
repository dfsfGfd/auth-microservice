// Package repository предоставляет интерфейсы для работы с хранилищем данных.
package repository

import (
	"context"

	"github.com/google/uuid"

	"auth-microservice/internal/model"
)

// UserRepository интерфейс для работы с пользователями в PostgreSQL.
type UserRepository interface {
	// Save сохраняет пользователя (создаёт или обновляет).
	Save(ctx context.Context, user *model.User) error

	// DeleteByID удаляет пользователя по ID.
	DeleteByID(ctx context.Context, id uuid.UUID) error

	// GetByID получает пользователя по ID.
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)

	// GetByEmail получает пользователя по email.
	GetByEmail(ctx context.Context, email string) (*model.User, error)

	// GetByUsername получает пользователя по username.
	GetByUsername(ctx context.Context, username string) (*model.User, error)

	// GetAll возвращает всех пользователей.
	GetAll(ctx context.Context) ([]*model.User, error)
}
