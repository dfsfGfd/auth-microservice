package auth

import (
	"context"

	"auth-microservice/internal/model"
	"auth-microservice/internal/repository/converter"
)

// Save сохраняет пользователя (создаёт или обновляет).
func (r *UserRepository) Save(ctx context.Context, user *model.User) error {
	dbUser := converter.UserToDB(user)

	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (email) DO UPDATE SET
			username = EXCLUDED.username,
			password_hash = EXCLUDED.password_hash,
			updated_at = EXCLUDED.updated_at
	`,
		dbUser.ID,
		dbUser.Email,
		dbUser.Username,
		dbUser.PasswordHash,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	)

	return err
}
