package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"auth-microservice/internal/model"
	"auth-microservice/internal/repository/converter"
	dbmodel "auth-microservice/internal/repository/model"
)

// GetByID получает пользователя по ID.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var dbUser dbmodel.User

	err := r.pool.QueryRow(ctx, `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&dbUser.ID,
		&dbUser.Email,
		&dbUser.Username,
		&dbUser.PasswordHash,
		&dbUser.CreatedAt,
		&dbUser.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return converter.UserToDomain(&dbUser)
}
