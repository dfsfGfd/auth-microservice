package auth

import (
	"context"
	stderrors "errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"auth-microservice/internal/errors"
	"auth-microservice/internal/model"
	"auth-microservice/internal/repository/converter"
	dbmodel "auth-microservice/internal/repository/model"
)

// GetByID получает аккаунт по ID.
func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	var dbAccount dbmodel.Account

	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`, id).Scan(
		&dbAccount.ID,
		&dbAccount.Email,
		&dbAccount.PasswordHash,
		&dbAccount.CreatedAt,
		&dbAccount.UpdatedAt,
	)

	if stderrors.Is(err, pgx.ErrNoRows) {
		return nil, errors.ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}

	return converter.AccountToDomain(&dbAccount)
}
