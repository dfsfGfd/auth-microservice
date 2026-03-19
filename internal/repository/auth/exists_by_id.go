package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// ExistsByID проверяет существование аккаунта по ID.
func (r *AccountRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool

	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM accounts WHERE id = $1)
	`, id).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check account exists by id [id=%s]: %w", id.String(), err)
	}

	return exists, nil
}
