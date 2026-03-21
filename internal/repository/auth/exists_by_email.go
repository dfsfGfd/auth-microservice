package auth

import (
	"context"
	"fmt"
)

// ExistsByEmail проверяет существование аккаунта по email.
func (r *AccountRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM accounts WHERE email = $1)
	`, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check account exists by email [email=%s]: %w", email, err)
	}

	return exists, nil
}
