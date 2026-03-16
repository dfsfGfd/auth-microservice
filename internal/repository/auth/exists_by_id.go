package auth

import (
	"context"

	"github.com/google/uuid"
)

// ExistsByID проверяет существование аккаунта по ID.
func (r *AccountRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)
	`, id).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
