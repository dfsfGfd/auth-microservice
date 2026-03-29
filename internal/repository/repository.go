// Package repository предоставляет интерфейсы для работы с хранилищем данных.
package repository

import (
	"context"

	"auth-microservice/internal/model"
)

// AccountRepository интерфейс для работы с аккаунтами в PostgreSQL.
type AccountRepository interface {
	// Save сохраняет аккаунт (создаёт или обновляет).
	Save(ctx context.Context, account *model.Account) error

	// GetByEmail получает аккаунт по email.
	GetByEmail(ctx context.Context, email string) (*model.Account, error)

	// ExistsByID проверяет существование аккаунта по ID.
	ExistsByID(ctx context.Context, id int64) (bool, error)

	// ExistsByEmail проверяет существование аккаунта по email.
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
