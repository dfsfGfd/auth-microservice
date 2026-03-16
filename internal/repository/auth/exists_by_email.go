package auth

import (
	"context"
)

// ExistsByEmail проверяет существование аккаунта по email.
func (r *AccountRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM accounts WHERE email = $1)
	`, email).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
