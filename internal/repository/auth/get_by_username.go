package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"auth-microservice/internal/model"
	"auth-microservice/internal/repository/converter"
	dbmodel "auth-microservice/internal/repository/model"
)

// GetByUsername получает пользователя по username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var dbUser dbmodel.User

	err := r.pool.QueryRow(ctx, `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE username = $1
	`, username).Scan(
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
