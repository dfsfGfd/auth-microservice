package auth

import (
	"context"

	"auth-microservice/internal/model"
	"auth-microservice/internal/repository/converter"
	dbmodel "auth-microservice/internal/repository/model"
)

// GetAll возвращает всех пользователей.
func (r *UserRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbUsers []*dbmodel.User
	for rows.Next() {
		var dbUser dbmodel.User
		err := rows.Scan(
			&dbUser.ID,
			&dbUser.Email,
			&dbUser.Username,
			&dbUser.PasswordHash,
			&dbUser.CreatedAt,
			&dbUser.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		dbUsers = append(dbUsers, &dbUser)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return converter.UserListToDomain(dbUsers)
}
