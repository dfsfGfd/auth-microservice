// Package auth предоставляет PostgreSQL реализацию репозиториев.
package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-microservice/internal/model"
	"auth-microservice/internal/repository"
	repoerrors "auth-microservice/internal/repository/errors"
	"auth-microservice/internal/repository/converter"
	dbmodel "auth-microservice/internal/repository/model"
)

// ensure UserRepository implements interface
var _ repository.UserRepository = (*UserRepository)(nil)

// UserRepository PostgreSQL реализация repository.UserRepository
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository создаёт новый UserRepository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

// ============================================================================
// Write Operations
// ============================================================================

// Create создаёт нового пользователя в хранилище.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	const query = `
		INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	dbUser := converter.UserToDB(user)

	_, err := r.pool.Exec(ctx, query,
		dbUser.ID,
		dbUser.Email,
		dbUser.Username,
		dbUser.PasswordHash,
		dbUser.CreatedAt,
		dbUser.UpdatedAt,
	)

	if err != nil {
		return mapPostgresError(err)
	}

	return nil
}

// Update обновляет существующего пользователя.
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	const query = `
		UPDATE users
		SET email = $2, username = $3, password_hash = $4, updated_at = $5
		WHERE id = $1
	`

	dbUser := converter.UserToDB(user)

	result, err := r.pool.Exec(ctx, query,
		dbUser.ID,
		dbUser.Email,
		dbUser.Username,
		dbUser.PasswordHash,
		dbUser.UpdatedAt,
	)

	if err != nil {
		return mapPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return repoerrors.ErrUserNotFound
	}

	return nil
}

// Delete удаляет пользователя по идентификатору.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return mapPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return repoerrors.ErrUserNotFound
	}

	return nil
}

// DeleteByEmail удаляет пользователя по email.
func (r *UserRepository) DeleteByEmail(ctx context.Context, email string) error {
	const query = `DELETE FROM users WHERE email = $1`

	result, err := r.pool.Exec(ctx, query, email)
	if err != nil {
		return mapPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return repoerrors.ErrUserNotFound
	}

	return nil
}

// DeleteByUsername удаляет пользователя по username.
func (r *UserRepository) DeleteByUsername(ctx context.Context, username string) error {
	const query = `DELETE FROM users WHERE username = $1`

	result, err := r.pool.Exec(ctx, query, username)
	if err != nil {
		return mapPostgresError(err)
	}

	if result.RowsAffected() == 0 {
		return repoerrors.ErrUserNotFound
	}

	return nil
}

// ============================================================================
// Read Operations
// ============================================================================

// GetByID получает пользователя по идентификатору.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	const query = `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	dbUser, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrUserNotFound
		}
		return nil, err
	}

	return converter.UserToDomain(dbUser)
}

// GetByEmail получает пользователя по email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	const query = `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	row := r.pool.QueryRow(ctx, query, email)

	dbUser, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrUserNotFound
		}
		return nil, err
	}

	return converter.UserToDomain(dbUser)
}

// GetByUsername получает пользователя по username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	const query = `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	row := r.pool.QueryRow(ctx, query, username)

	dbUser, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrUserNotFound
		}
		return nil, err
	}

	return converter.UserToDomain(dbUser)
}

// Find находит пользователей по спецификации.
func (r *UserRepository) Find(ctx context.Context, spec repository.UserSpec) ([]*model.User, error) {
	query, args := buildFindQuery(spec)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, mapPostgresError(err)
	}
	defer rows.Close()

	dbUsers, err := scanUsers(rows)
	if err != nil {
		return nil, err
	}

	return converter.UserListToDomain(dbUsers)
}

// List возвращает список пользователей с пагинацией.
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	spec := repository.UserSpec{
		Limit:  limit,
		Offset: offset,
	}
	return r.Find(ctx, spec)
}

// ============================================================================
// Check Operations
// ============================================================================

// Exists проверяет существование пользователя по идентификатору.
func (r *UserRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, mapPostgresError(err)
	}

	return exists, nil
}

// ExistsByEmail проверяет существование пользователя по email.
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, mapPostgresError(err)
	}

	return exists, nil
}

// ExistsByUsername проверяет существование пользователя по username.
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, mapPostgresError(err)
	}

	return exists, nil
}

// ============================================================================
// Count Operations
// ============================================================================

// Count подсчитывает общее количество пользователей.
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	const query = `SELECT COUNT(*) FROM users`

	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, mapPostgresError(err)
	}

	return count, nil
}

// CountByEmail подсчитывает количество пользователей с указанным email.
func (r *UserRepository) CountByEmail(ctx context.Context, email string) (int64, error) {
	const query = `SELECT COUNT(*) FROM users WHERE email = $1`

	var count int64
	err := r.pool.QueryRow(ctx, query, email).Scan(&count)
	if err != nil {
		return 0, mapPostgresError(err)
	}

	return count, nil
}

// CountByUsername подсчитывает количество пользователей с указанным username.
func (r *UserRepository) CountByUsername(ctx context.Context, username string) (int64, error) {
	const query = `SELECT COUNT(*) FROM users WHERE username = $1`

	var count int64
	err := r.pool.QueryRow(ctx, query, username).Scan(&count)
	if err != nil {
		return 0, mapPostgresError(err)
	}

	return count, nil
}

// ============================================================================
// Helper Functions
// ============================================================================

// scanUser сканирует пользователя из одной строки результата.
func scanUser(row pgx.Row) (*dbmodel.User, error) {
	var user dbmodel.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, mapPostgresError(err)
	}
	return &user, nil
}

// scanUsers сканирует список пользователей из строк результата.
func scanUsers(rows pgx.Rows) ([]*dbmodel.User, error) {
	var users []*dbmodel.User

	for rows.Next() {
		var user dbmodel.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, mapPostgresError(err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, mapPostgresError(err)
	}

	return users, nil
}

// buildFindQuery строит SQL запрос для поиска по спецификации.
func buildFindQuery(spec repository.UserSpec) (string, []interface{}) {
	spec.Validate()

	query := `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE 1=1
	`

	args := make([]interface{}, 0)
	argIndex := 1

	if spec.Email != nil {
		query += fmt.Sprintf(" AND email = $%d", argIndex)
		args = append(args, *spec.Email)
		argIndex++
	}

	if spec.Username != nil {
		query += fmt.Sprintf(" AND username = $%d", argIndex)
		args = append(args, *spec.Username)
		argIndex++
	}

	if len(spec.IDs) > 0 {
		query += fmt.Sprintf(" AND id = ANY($%d)", argIndex)
		args = append(args, spec.IDs)
		argIndex++
	}

	if spec.CreatedAfter != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *spec.CreatedAfter)
		argIndex++
	}

	if spec.CreatedBefore != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *spec.CreatedBefore)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY %s %s", spec.OrderBy, spec.OrderDir)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, spec.Limit, spec.Offset)

	return query, args
}

// mapPostgresError преобразует ошибки PostgreSQL в доменные ошибки.
func mapPostgresError(err error) error {
	if err == nil {
		return nil
	}

	// Unique violation (23505)
	if isUniqueViolation(err) {
		return repoerrors.ErrUniqueViolation
	}

	// Foreign key violation (23503)
	if isForeignKeyViolation(err) {
		return repoerrors.ErrForeignKeyViolation
	}

	return err
}

// isUniqueViolation проверяет нарушение уникальности.
func isUniqueViolation(err error) bool {
	// TODO: реализовать проверку через pgconn
	return false
}

// isForeignKeyViolation проверяет нарушение внешнего ключа.
func isForeignKeyViolation(err error) bool {
	// TODO: реализовать проверку через pgconn
	return false
}
