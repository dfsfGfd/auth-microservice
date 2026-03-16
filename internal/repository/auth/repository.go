// Package auth предоставляет PostgreSQL реализацию репозитория.
package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-microservice/internal/repository"
)

// AccountRepository реализация repository.AccountRepository.
type AccountRepository struct {
	pool *pgxpool.Pool
}

// NewAccountRepository создаёт новый AccountRepository.
func NewAccountRepository(pool *pgxpool.Pool) repository.AccountRepository {
	return &AccountRepository{
		pool: pool,
	}
}
