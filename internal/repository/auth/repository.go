// Package auth предоставляет PostgreSQL реализацию репозитория.
package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-microservice/internal/repository"
)

// UserRepository реализация repository.UserRepository.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository создаёт новый UserRepository.
func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &UserRepository{
		pool: pool,
	}
}
