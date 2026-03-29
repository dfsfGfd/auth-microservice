package auth

import (
	"context"
	stderrors "errors"
	"fmt"

	"auth-microservice/internal/errors"
	"auth-microservice/internal/model"
	"auth-microservice/internal/model/converter"
	dbmodel "auth-microservice/internal/repository/model"
	"github.com/jackc/pgx/v5"
)

// GetByEmail получает аккаунт по email.
func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*model.Account, error) {
	var dbAccount dbmodel.Account

	err := r.pool.QueryRow(ctx, `
		SELECT id, email, password, created_at, updated_at
		FROM accounts
		WHERE email = $1
	`, email).Scan(
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
		return nil, fmt.Errorf("get account by email [email=%s]: %w", email, err)
	}

	return converter.AccountToDomain(&dbAccount)
}
