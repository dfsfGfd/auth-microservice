# Repository Methods — DDD Implementation Guide

Полный список методов для реализации репозиториев в соответствии с DDD паттернами.

---

## 📋 Оглавление

- [Структура пакета](#структура-пакета)
- [UserRepository методы](#userrepository-методы)
- [UnitOfWork методы](#unitofwork-методы)
- [Схема именования](#схема-именования)
- [Доменные ошибки](#доменные-ошибки)

---

## Структура пакета

```
internal/repository/
├── errors/
│   └── errors.go           # Доменные ошибки репозитория
├── model/
│   └── user.go             # DB модели (маппинг)
├── converter/
│   └── user.go             # Конвертеры domain ↔ DB
├── auth/                   # PostgreSQL реализация
│   ├── user_repository.go  # UserRepository реализация
│   ├── transaction.go      # UnitOfWork реализация
│   └── doc.go
├── user_repository.go      # UserRepository интерфейс
└── token_repository.go     # TokenRepository интерфейс (Redis, отдельно)
```

---

## UserRepository методы

### Write Operations (Операции записи)

| Метод | Сигнатура | Описание | Ошибки |
|-------|-----------|----------|--------|
| `Create` | `Create(ctx, user) error` | Создание нового пользователя | `ErrUserAlreadyExists`, `ErrEmailAlreadyExists`, `ErrUsernameAlreadyExists` |
| `Update` | `Update(ctx, user) error` | Обновление существующего | `ErrUserNotFound`, `ErrEmailAlreadyExists`, `ErrUsernameAlreadyExists` |
| `Delete` | `Delete(ctx, id) error` | Удаление по ID | `ErrUserNotFound` |
| `DeleteByEmail` | `DeleteByEmail(ctx, email) error` | Удаление по email | `ErrUserNotFound` |
| `DeleteByUsername` | `DeleteByUsername(ctx, username) error` | Удаление по username | `ErrUserNotFound` |

### Read Operations (Операции чтения)

| Метод | Сигнатура | Описание | Возврат | Ошибки |
|-------|-----------|----------|---------|--------|
| `GetByID` | `GetByID(ctx, id) (*User, error)` | Получение по ID | `*model.User` | `ErrUserNotFound` |
| `GetByEmail` | `GetByEmail(ctx, email) (*User, error)` | Получение по email | `*model.User` | `ErrUserNotFound` |
| `GetByUsername` | `GetByUsername(ctx, username) (*User, error)` | Получение по username | `*model.User` | `ErrUserNotFound` |
| `Find` | `Find(ctx, spec) ([]*User, error)` | Поиск по спецификации | `[]*model.User` | — |
| `List` | `List(ctx, limit, offset) ([]*User, error)` | Список с пагинацией | `[]*model.User` | — |

### Check Operations (Проверки)

| Метод | Сигнатура | Описание | Возврат |
|-------|-----------|----------|---------|
| `Exists` | `Exists(ctx, id) (bool, error)` | Проверка по ID | `bool` |
| `ExistsByEmail` | `ExistsByEmail(ctx, email) (bool, error)` | Проверка по email | `bool` |
| `ExistsByUsername` | `ExistsByUsername(ctx, username) (bool, error)` | Проверка по username | `bool` |

### Count Operations (Подсчёт)

| Метод | Сигнатура | Описание | Возврат |
|-------|-----------|----------|---------|
| `Count` | `Count(ctx) (int64, error)` | Общее количество | `int64` |
| `CountByEmail` | `CountByEmail(ctx, email) (int64, error)` | Количество по email | `int64` |
| `CountByUsername` | `CountByUsername(ctx, username) (int64, error)` | Количество по username | `int64` |

---

## UnitOfWork методы

| Метод | Сигнатура | Описание |
|-------|-----------|----------|
| `Begin` | `Begin(ctx) (context.Context, error)` | Начало транзакции |
| `Commit` | `Commit(ctx) error` | Фиксация транзакции |
| `Rollback` | `Rollback(ctx) error` | Откат транзакции |
| `Users` | `Users() UserRepository` | UserRepository в транзакции |

---

## Схема именования

### Формат имён методов

```
{Operation}{By}{Criteria}
```

**Операции:**
- `Get` — получение одного объекта
- `Find` — поиск нескольких объектов
- `Create` — создание
- `Update` — обновление
- `Delete` — удаление
- `Exists` — проверка существования
- `Count` — подсчёт количества

**Критерии:**
- `ByID` — по идентификатору
- `ByEmail` — по email
- `ByUsername` — по username
- `ByUserID` — по ID пользователя

### Примеры

```go
// Правильно
repo.GetByID(ctx, id)
repo.GetByEmail(ctx, email)
repo.Find(ctx, spec)
repo.Exists(ctx, id)
repo.Count(ctx)

// Неправильно
repo.GetUserById(ctx, id)      // избыточно "User"
repo.FindUsers(ctx, spec)      // избыточно "Users"
repo.CheckExists(ctx, id)      // избыточно "Check"
```

---

## Доменные ошибки

### Базовые ошибки (`repository/errors/errors.go`)

```go
// Общие ошибки репозитория
ErrNotFound              // Агрегат не найден
ErrAlreadyExists         // Агрегат уже существует
ErrUniqueViolation       // Нарушение уникальности
ErrForeignKeyViolation   // Нарушение внешнего ключа
ErrTransactionActive     // Активная транзакция
ErrNoTransaction         // Нет активной транзакции
```

### Ошибки пользователя

```go
ErrUserNotFound          // Пользователь не найден
ErrUserAlreadyExists     // Пользователь уже существует
ErrEmailAlreadyExists    // Email уже зарегистрирован
ErrUsernameAlreadyExists // Username уже занят
```

### Ошибки токена

```go
ErrTokenNotFound         // Токен не найден
ErrTokenExpired          // Токен истёк
```

### Использование ошибок

```go
user, err := repo.GetByID(ctx, id)
if err != nil {
    if errors.Is(err, repoerrors.ErrUserNotFound) {
        // Обработка "не найден"
    }
    return err
}

exists, err := repo.ExistsByEmail(ctx, email)
if err != nil {
    return err
}
if exists {
    return repoerrors.ErrEmailAlreadyExists
}
```

---

## Примеры использования

### Базовое использование

```go
// Создание пользователя
user, err := model.NewUser(email, username, password)
if err != nil {
    return err
}

if err := userRepo.Create(ctx, user); err != nil {
    return err
}

// Получение пользователя
user, err := userRepo.GetByID(ctx, userID)
if err != nil {
    return err
}

// Обновление
user.UpdateEmail(newEmail)
if err := userRepo.Update(ctx, user); err != nil {
    return err
}
```

### Транзакции

```go
// Начало транзакции
txCtx, err := uow.Begin(ctx)
if err != nil {
    return err
}
defer func() {
    if r := recover(); r != nil {
        uow.Rollback(txCtx)
    }
}()

// Операции в транзакции
user, _ := uow.Users().GetByID(txCtx, userID)
user.UpdateEmail(newEmail)
if err := uow.Users().Update(txCtx, user); err != nil {
    uow.Rollback(txCtx)
    return err
}

// Фиксация
if err := uow.Commit(txCtx); err != nil {
    return err
}
```

### Спецификация поиска

```go
spec := repository.UserSpec{
    Email:        ptr.String("user@example.com"),
    CreatedAfter: ptr.Time(time.Now().Add(-24 * time.Hour)),
    Limit:        10,
    Offset:       0,
    OrderBy:      "created_at",
    OrderDir:     "DESC",
}

users, err := userRepo.Find(ctx, spec)
```

---

## Ссылки

- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [Unit of Work](https://martinfowler.com/eaaCatalog/unitOfWork.html)
- [DDD Repository](https://domainlanguage.com/ddd/repositories/)
