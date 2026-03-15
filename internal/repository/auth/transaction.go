// Package auth предоставляет реализацию транзакций для PostgreSQL.
package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	repoerrors "auth-microservice/internal/repository/errors"
)

// ensure UnitOfWork implements interface
var _ UnitOfWork = (*UnitOfWorkImpl)(nil)

// UnitOfWork интерфейс для группировки операций в транзакцию.
type UnitOfWork interface {
	// Begin начинает транзакцию
	Begin(ctx context.Context) (context.Context, error)

	// Commit фиксирует транзакцию
	Commit(ctx context.Context) error

	// Rollback откатывает транзакцию
	Rollback(ctx context.Context) error

	// Users возвращает UserRepository в контексте транзакции
	Users() UserRepository

	// Tokens возвращает TokenRepository в контексте транзакции
	Tokens() TokenRepository
}

// UnitOfWorkImpl реализация UnitOfWork с поддержкой транзакций.
type UnitOfWorkImpl struct {
	pool       *pgxpool.Pool
	tx         pgx.Tx
	userRepo   *UserRepository
	tokenRepo  *TokenRepository
}

// NewUnitOfWork создаёт новый UnitOfWork
func NewUnitOfWork(pool *pgxpool.Pool, userRepo *UserRepository, tokenRepo *TokenRepository) *UnitOfWorkImpl {
	return &UnitOfWorkImpl{
		pool:      pool,
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

// Begin начинает транзакцию.
func (u *UnitOfWorkImpl) Begin(ctx context.Context) (context.Context, error) {
	if u.tx != nil {
		return nil, repoerrors.ErrTransactionActive
	}

	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	u.tx = tx

	// Создаём репозитории с транзакцией
	u.userRepo = &UserRepository{pool: tx}

	return ctx, nil
}

// Commit фиксирует транзакцию.
func (u *UnitOfWorkImpl) Commit(ctx context.Context) error {
	if u.tx == nil {
		return repoerrors.ErrNoTransaction
	}

	err := u.tx.Commit(ctx)
	u.tx = nil

	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Rollback откатывает транзакцию.
func (u *UnitOfWorkImpl) Rollback(ctx context.Context) error {
	if u.tx == nil {
		return repoerrors.ErrNoTransaction
	}

	err := u.tx.Rollback(ctx)
	u.tx = nil

	if err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return nil
}

// Users возвращает UserRepository в контексте транзакции.
func (u *UnitOfWorkImpl) Users() UserRepository {
	return u.userRepo
}

// Tokens возвращает TokenRepository в контексте транзакции.
func (u *UnitOfWorkImpl) Tokens() TokenRepository {
	return u.tokenRepo
}

// txPool обёртка для использования транзакции как пула.
type txPool struct {
	tx pgx.Tx
}

func (t *txPool) Exec(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return t.tx.Exec(ctx, sql, args...)
}

func (t *txPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return t.tx.Query(ctx, sql, args...)
}

func (t *txPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return t.tx.QueryRow(ctx, sql, args...)
}
