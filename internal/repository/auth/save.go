package auth

import (
	"context"
	"fmt"

	"auth-microservice/internal/model"
	"auth-microservice/internal/repository/converter"
)

// Save сохраняет аккаунт (создаёт или обновляет).
func (r *AccountRepository) Save(ctx context.Context, account *model.Account) error {
	dbAccount := converter.AccountToDB(account)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO accounts (id, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (email) DO UPDATE SET
			password = EXCLUDED.password,
			updated_at = EXCLUDED.updated_at
	`,
		dbAccount.ID,
		dbAccount.Email,
		dbAccount.PasswordHash,
		dbAccount.CreatedAt,
		dbAccount.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("save account [email=%s]: %w", account.Email().Value(), err)
	}
	return nil
}
