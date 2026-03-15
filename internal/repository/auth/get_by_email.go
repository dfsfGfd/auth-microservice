package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"auth-microservice/internal/model"
	"auth-microservice/internal/repository/converter"
	dbmodel "auth-microservice/internal/repository/model"
)

// GetByEmail получает пользователя по email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var dbUser dbmodel.User

	err := r.pool.QueryRow(ctx, `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
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
