// Package auth предоставляет PostgreSQL реализацию репозитория.
package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// AccountRepository реализация repository.AccountRepository.
type AccountRepository struct {
	pool *pgxpool.Pool
}

// NewAccountRepository создаёт новый AccountRepository.
func NewAccountRepository(pool *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{
		pool: pool,
	}
}
