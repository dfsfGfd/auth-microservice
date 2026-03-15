package auth

import (
	"context"

	"github.com/google/uuid"
)

// DeleteByID удаляет пользователя по ID.
func (r *UserRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}
