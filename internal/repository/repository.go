// Package repository предоставляет интерфейсы для работы с хранилищем данных.
package repository

import (
	"context"

	"github.com/google/uuid"

	"auth-microservice/internal/model"
)

// AccountRepository интерфейс для работы с аккаунтами в PostgreSQL.
type AccountRepository interface {
	// Save сохраняет аккаунт (создаёт или обновляет).
	Save(ctx context.Context, account *model.Account) error

	// DeleteByID удаляет аккаунт по ID.
	DeleteByID(ctx context.Context, id uuid.UUID) error

	// GetByID получает аккаунт по ID.
	GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error)

	// GetByEmail получает аккаунт по email.
	GetByEmail(ctx context.Context, email string) (*model.Account, error)
}
