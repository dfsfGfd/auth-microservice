package auth

import (
	"context"

	"auth-microservice/internal/model"
	"auth-microservice/internal/repository/converter"
	dbmodel "auth-microservice/internal/repository/model"
)

// GetAll возвращает все аккаунты.
func (r *AccountRepository) GetAll(ctx context.Context) ([]*model.Account, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, email, password, created_at, updated_at
		FROM accounts
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbAccounts []*dbmodel.Account
	for rows.Next() {
		var dbAccount dbmodel.Account
		err := rows.Scan(
			&dbAccount.ID,
			&dbAccount.Email,
			&dbAccount.PasswordHash,
			&dbAccount.CreatedAt,
			&dbAccount.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		dbAccounts = append(dbAccounts, &dbAccount)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return converter.AccountListToDomain(dbAccounts)
}
