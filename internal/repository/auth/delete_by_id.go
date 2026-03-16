package auth

import (
	"context"

	"github.com/google/uuid"
)

// DeleteByID удаляет аккаунт по ID.
func (r *AccountRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM accounts WHERE id = $1`, id)
	return err
}
