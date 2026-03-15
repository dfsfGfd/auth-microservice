# Repository Layer — DDD Pattern

Документ описывает требования к слою репозиториев в соответствии с паттернами Domain-Driven Design (DDD).

---

## 📋 Оглавление

- [Обзор](#обзор)
- [Структура пакета](#структура-пакета)
- [Неминг конвенции](#нейминг-конвенции)
- [Интерфейсы репозиториев](#интерфейсы-репозиториев)
- [Реализация](#реализация)
- [Примеры использования](#примеры-использования)
- [Best Practices](#best-practices)

---

## Обзор

Репозиторий в DDD — это абстракция над слоем хранения данных. Он предоставляет коллекционно-ориентированный интерфейс для доступа к агрегатам.

**Принципы:**

| Принцип | Описание |
|---------|----------|
| **Агрегат-ориентированность** | Репозитории работают с агрегатами, а не с отдельными таблицами |
| **Инкапсуляция** | Логика доступа к данным скрыта внутри репозитория |
| **Интерфейсы в domain** | Интерфейсы определяются в domain слое, реализация — в infrastructure |
| **Единая ответственность** | Один репозиторий = один тип агрегата |

---

## Структура пакета

```
internal/
├── model/                      # Domain модели (агрегаты, VO)
│   ├── user.go
│   └── ...
│
├── repository/                 # Repository layer
│   ├── repository.go           # Общие интерфейсы и типы
│   ├── user_repository.go      # UserRepository интерфейс
│   └── postgres/               # PostgreSQL реализация
│       ├── user_repository.go  # Реализация UserRepository
│       └── ...
│
└── service/                    # Service layer (использует репозитории)
```

---

## Нейминг конвенции

### Имена репозиториев

| Паттерн | Пример | Описание |
|---------|--------|----------|
| `{Aggregate}Repository` | `UserRepository` | Интерфейс репозитория |
| `{aggregate}Repository` | `userRepository` | Переменная репозитория |
| `{Vendor}{Aggregate}Repository` | `PostgresUserRepository` | Реализация |

### Имена методов

| Операция | Метод | Описание |
|----------|-------|----------|
| Создание | `Create(ctx, aggregate)` | Сохранение нового агрегата |
| Чтение | `GetByID(ctx, id)` | Получение по ID |
| Чтение | `Find(ctx, spec)` | Поиск по спецификации |
| Обновление | `Update(ctx, aggregate)` | Обновление существующего |
| Удаление | `Delete(ctx, id)` | Удаление по ID |
| Существование | `Exists(ctx, id)` | Проверка существования |
| Количество | `Count(ctx)` | Подсчёт количества |

### Спецификации (Search Criteria)

```go
// Спецификация для поиска пользователей
type UserSpec struct {
    Email    *string
    Username *string
    Limit    int
    Offset   int
}
```

---

## Интерфейсы репозиториев

### Общий интерфейс репозитория

```go
// repository.go

package repository

import (
    "context"
    
    "auth-microservice/internal/model"
)

// TransactionManager управляет транзакциями
type TransactionManager interface {
    // Begin начинает транзакцию
    Begin(ctx context.Context) (context.Context, error)
    
    // Commit фиксирует транзакцию
    Commit(ctx context.Context) error
    
    // Rollback откатывает транзакцию
    Rollback(ctx context.Context) error
}

// UnitOfWork паттерн для группировки операций
type UnitOfWork interface {
    TransactionManager
    Users() UserRepository
}
```

### UserRepository

```go
// user_repository.go

package repository

import (
    "context"
    
    "auth-microservice/internal/model"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
    // Create создаёт нового пользователя
    // Возвращает ошибку если пользователь уже существует
    Create(ctx context.Context, user *model.User) error
    
    // GetByID получает пользователя по ID
    // Возвращает model.ErrUserNotFound если пользователь не найден
    GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
    
    // GetByEmail получает пользователя по email
    // Возвращает model.ErrUserNotFound если пользователь не найден
    GetByEmail(ctx context.Context, email string) (*model.User, error)
    
    // GetByUsername получает пользователя по username
    // Возвращает model.ErrUserNotFound если пользователь не найден
    GetByUsername(ctx context.Context, username string) (*model.User, error)
    
    // Find находит пользователей по спецификации
    Find(ctx context.Context, spec UserSpec) ([]*model.User, error)
    
    // Update обновляет существующего пользователя
    // Возвращает ошибку если пользователь не найден
    Update(ctx context.Context, user *model.User) error
    
    // Delete удаляет пользователя по ID
    // Возвращает ошибку если пользователь не найден
    Delete(ctx context.Context, id uuid.UUID) error
    
    // Exists проверяет существование пользователя по ID
    Exists(ctx context.Context, id uuid.UUID) (bool, error)
    
    // ExistsByEmail проверяет существование пользователя по email
    ExistsByEmail(ctx context.Context, email string) (bool, error)
    
    // ExistsByUsername проверяет существование пользователя по username
    ExistsByUsername(ctx context.Context, username string) (bool, error)
    
    // Count подсчитывает количество пользователей
    Count(ctx context.Context) (int64, error)
}

// UserSpec спецификация для поиска пользователей
type UserSpec struct {
    // Email для фильтрации по email
    Email *string
    
    // Username для фильтрации по username
    Username *string
    
    // Limit максимальное количество результатов
    Limit int
    
    // Offset смещение результатов
    Offset int
}

// DefaultUserSpec возвращает спецификацию по умолчанию
func DefaultUserSpec() UserSpec {
    return UserSpec{
        Limit:  100,
        Offset: 0,
    }
}
```

### TokenRepository (для refresh токенов)

```go
// token_repository.go

package repository

import (
    "context"
    "time"
)

// TokenRepository интерфейс для работы с токенами
type TokenRepository interface {
    // Store сохраняет refresh токен
    // ttl — время жизни токена
    Store(ctx context.Context, token string, userID string, ttl time.Duration) error
    
    // GetUserID получает user_id по токену
    // Возвращает пустую строку если токен не найден
    GetUserID(ctx context.Context, token string) (string, error)
    
    // Delete удаляет токен (отзыв)
    Delete(ctx context.Context, token string) error
    
    // Exists проверяет существование токена
    Exists(ctx context.Context, token string) (bool, error)
    
    // TTL получает оставшееся время жизни токена
    TTL(ctx context.Context, token string) (time.Duration, error)
}
```

---

## Реализация

### PostgreSQL UserRepository

```go
// postgres/user_repository.go

package postgres

import (
    "context"
    "errors"
    "time"
    
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    
    "auth-microservice/internal/model"
    "auth-microservice/internal/repository"
)

// UserRepository реализация repository.UserRepository
type UserRepository struct {
    pool *pgxpool.Pool
}

// NewUserRepository создаёт новый UserRepository
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
    return &UserRepository{
        pool: pool,
    }
}

// Create создаёт нового пользователя
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    const query = `
        INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    
    _, err := r.pool.Exec(ctx, query,
        user.ID(),
        user.Email().String(),
        user.Username().String(),
        user.PasswordHash().String(),
        user.CreatedAt(),
        user.UpdatedAt(),
    )
    
    if err != nil {
        return handleError(err)
    }
    
    return nil
}

// GetByID получает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
    const query = `
        SELECT id, email, username, password_hash, created_at, updated_at
        FROM users
        WHERE id = $1
    `
    
    var user *model.User
    err := r.pool.QueryRow(ctx, query, id).Scan(
        &user.id,
        &user.email,
        &user.username,
        &user.passwordHash,
        &user.createdAt,
        &user.updatedAt,
    )
    
    if errors.Is(err, pgx.ErrNoRows) {
        return nil, model.ErrUserNotFound
    }
    if err != nil {
        return nil, handleError(err)
    }
    
    return user, nil
}

// GetByEmail получает пользователя по email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
    const query = `
        SELECT id, email, username, password_hash, created_at, updated_at
        FROM users
        WHERE email = $1
    `
    
    // ... аналогично GetByID
}

// Find находит пользователей по спецификации
func (r *UserRepository) Find(ctx context.Context, spec repository.UserSpec) ([]*model.User, error) {
    query := `
        SELECT id, email, username, password_hash, created_at, updated_at
        FROM users
        WHERE 1=1
    `
    
    args := []interface{}{}
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
    
    query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
    args = append(args, spec.Limit, spec.Offset)
    
    rows, err := r.pool.Query(ctx, query, args...)
    if err != nil {
        return nil, handleError(err)
    }
    defer rows.Close()
    
    var users []*model.User
    for rows.Next() {
        user, err := scanUser(rows)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, nil
}

// Update обновляет существующего пользователя
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
    const query = `
        UPDATE users
        SET email = $2, username = $3, password_hash = $4, updated_at = $5
        WHERE id = $1
    `
    
    result, err := r.pool.Exec(ctx, query,
        user.ID(),
        user.Email().String(),
        user.Username().String(),
        user.PasswordHash().String(),
        user.UpdatedAt(),
    )
    
    if err != nil {
        return handleError(err)
    }
    
    if result.RowsAffected() == 0 {
        return model.ErrUserNotFound
    }
    
    return nil
}

// Delete удаляет пользователя по ID
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
    const query = `DELETE FROM users WHERE id = $1`
    
    result, err := r.pool.Exec(ctx, query, id)
    if err != nil {
        return handleError(err)
    }
    
    if result.RowsAffected() == 0 {
        return model.ErrUserNotFound
    }
    
    return nil
}

// Exists проверяет существование пользователя по ID
func (r *UserRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
    const query = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
    
    var exists bool
    err := r.pool.QueryRow(ctx, query, id).Scan(&exists)
    if err != nil {
        return false, handleError(err)
    }
    
    return exists, nil
}

// ExistsByEmail проверяет существование пользователя по email
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
    const query = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
    
    var exists bool
    err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
    if err != nil {
        return false, handleError(err)
    }
    
    return exists, nil
}

// ExistsByUsername проверяет существование пользователя по username
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
    const query = `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`
    
    var exists bool
    err := r.pool.QueryRow(ctx, query, username).Scan(&exists)
    if err != nil {
        return false, handleError(err)
    }
    
    return exists, nil
}

// Count подсчитывает количество пользователей
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
    const query = `SELECT COUNT(*) FROM users`
    
    var count int64
    err := r.pool.QueryRow(ctx, query).Scan(&count)
    if err != nil {
        return 0, handleError(err)
    }
    
    return count, nil
}

// handleError преобразует ошибки БД в доменные ошибки
func handleError(err error) error {
    if err == nil {
        return nil
    }
    
    // Проверка на уникальность
    if isUniqueViolation(err) {
        return model.ErrUserAlreadyExists
    }
    
    // Проверка на foreign key
    if isForeignKeyViolation(err) {
        return model.ErrUserNotFound
    }
    
    return err
}

// scanUser сканирует пользователя из строки результата
func scanUser(row pgx.Row) (*model.User, error) {
    var user model.User
    err := row.Scan(
        &user.id,
        &user.email,
        &user.username,
        &user.passwordHash,
        &user.createdAt,
        &user.updatedAt,
    )
    if err != nil {
        return nil, handleError(err)
    }
    return &user, nil
}
```

### Redis TokenRepository

```go
// postgres/token_repository.go

package postgres

import (
    "context"
    "fmt"
    "time"
    
    "github.com/redis/go-redis/v9"
    
    "auth-microservice/internal/repository"
)

// TokenRepository реализация repository.TokenRepository
type TokenRepository struct {
    client *redis.Client
    prefix string
}

// NewTokenRepository создаёт новый TokenRepository
func NewTokenRepository(client *redis.Client, prefix string) *TokenRepository {
    return &TokenRepository{
        client: client,
        prefix: prefix,
    }
}

// Store сохраняет refresh токен
func (r *TokenRepository) Store(ctx context.Context, token string, userID string, ttl time.Duration) error {
    key := r.key(token)
    
    err := r.client.Set(ctx, key, userID, ttl).Err()
    if err != nil {
        return fmt.Errorf("failed to store token: %w", err)
    }
    
    return nil
}

// GetUserID получает user_id по токену
func (r *TokenRepository) GetUserID(ctx context.Context, token string) (string, error) {
    key := r.key(token)
    
    userID, err := r.client.Get(ctx, key).Result()
    if err != nil {
        if errors.Is(err, redis.Nil) {
            return "", nil
        }
        return "", fmt.Errorf("failed to get token: %w", err)
    }
    
    return userID, nil
}

// Delete удаляет токен (отзыв)
func (r *TokenRepository) Delete(ctx context.Context, token string) error {
    key := r.key(token)
    
    err := r.client.Del(ctx, key).Err()
    if err != nil {
        return fmt.Errorf("failed to delete token: %w", err)
    }
    
    return nil
}

// Exists проверяет существование токена
func (r *TokenRepository) Exists(ctx context.Context, token string) (bool, error) {
    key := r.key(token)
    
    result, err := r.client.Exists(ctx, key).Result()
    if err != nil {
        return false, fmt.Errorf("failed to check token: %w", err)
    }
    
    return result > 0, nil
}

// TTL получает оставшееся время жизни токена
func (r *TokenRepository) TTL(ctx context.Context, token string) (time.Duration, error) {
    key := r.key(token)
    
    ttl, err := r.client.TTL(ctx, key).Result()
    if err != nil {
        return 0, fmt.Errorf("failed to get token TTL: %w", err)
    }
    
    return ttl, nil
}

// key формирует ключ с префиксом
func (r *TokenRepository) key(token string) string {
    return fmt.Sprintf("%s:refresh:%s", r.prefix, token)
}
```

---

## Примеры использования

### Базовое использование

```go
type AuthService struct {
    users repository.UserRepository
    tokens repository.TokenRepository
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (*model.User, error) {
    // Проверка на существование
    exists, err := s.users.ExistsByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, model.ErrEmailAlreadyExists
    }
    
    // Создание VO
    emailVO, err := model.NewEmail(email)
    if err != nil {
        return nil, err
    }
    
    usernameVO, err := model.NewUsername(username)
    if err != nil {
        return nil, err
    }
    
    passwordVO, err := model.NewPassword(password)
    if err != nil {
        return nil, err
    }
    
    // Создание агрегата
    user, err := model.NewUser(emailVO, usernameVO, passwordVO)
    if err != nil {
        return nil, err
    }
    
    // Сохранение через репозиторий
    if err := s.users.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

### Транзакции

```go
func (s *AuthService) TransferUser(ctx context.Context, fromID, toID uuid.UUID) error {
    // Начало транзакции
    txCtx, err := s.uow.Begin(ctx)
    if err != nil {
        return err
    }
    
    defer func() {
        if r := recover(); r != nil {
            s.uow.Rollback(txCtx)
        }
    }()
    
    // Получение пользователя
    user, err := s.uow.Users().GetByID(txCtx, fromID)
    if err != nil {
        s.uow.Rollback(txCtx)
        return err
    }
    
    // Обновление
    // ...
    
    if err := s.uow.Users().Update(txCtx, user); err != nil {
        s.uow.Rollback(txCtx)
        return err
    }
    
    // Фиксация
    if err := s.uow.Commit(txCtx); err != nil {
        return err
    }
    
    return nil
}
```

---

## Best Practices

### ✅ Делайте

- Используйте интерфейсы в service layer
- Возвращайте доменные ошибки, а не ошибки БД
- Используйте спецификации для сложных запросов
- Инкапсулируйте логику в репозитории
- Используйте транзакции для группировки операций
- Закрывайте подключения и итераторы

### ❌ Не делайте

- Не возвращайте `sql.Rows` или `pgx.Rows` из репозитория
- Не используйте репозитории для бизнес-логики
- Не создавайте репозитории для отдельных таблиц
- Не смешивайте разные источники данных в одном репозитории
- Не используйте `interface{}` без необходимости

---

## Ссылки

- [Domain-Driven Design](https://domainlanguage.com/ddd/)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [Unit of Work](https://martinfowler.com/eaaCatalog/unitOfWork.html)
- [Основное README](../../README.md)
