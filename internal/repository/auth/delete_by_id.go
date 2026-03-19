package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteByID удаляет аккаунт по ID.
func (r *AccountRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM accounts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete account by id [id=%s]: %w", id.String(), err)
	}
	return nil
}
