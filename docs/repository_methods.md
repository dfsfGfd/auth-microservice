# Repository Methods — DDD Implementation Guide

Полный список методов для реализации репозиториев в соответствии с DDD паттернами.

---

## 📋 Оглавление

- [Структура пакета](#структура-пакета)
- [UserRepository методы](#userrepository-методы)
- [Схема именования](#схема-именования)
- [Доменные ошибки](#доменные-ошибки)

---

## Структура пакета

```
internal/repository/
├── errors/
│   └── errors.go           # Доменные ошибки
├── model/
│   └── user.go             # DB модели
├── converter/
│   └── user.go             # Конвертеры domain ↔ DB
├── auth/
│   ├── repository.go       # Конструктор
│   ├── save.go             # Save метод
│   ├── delete_by_id.go     # DeleteByID метод
│   ├── get_by_id.go        # GetByID метод
│   ├── get_by_email.go     # GetByEmail метод
│   ├── get_by_username.go  # GetByUsername метод
│   └── get_all.go          # GetAll метод
└── repository.go           # Интерфейс UserRepository
```

---

## UserRepository методы

| Метод | Файл | Сигнатура | Описание |
|-------|------|-----------|----------|
| `Save` | `save.go` | `Save(ctx, user) error` | Сохранение (создание или обновление) |
| `DeleteByID` | `delete_by_id.go` | `DeleteByID(ctx, id) error` | Удаление по ID |
| `GetByID` | `get_by_id.go` | `GetByID(ctx, id) (*User, error)` | Получение по ID |
| `GetByEmail` | `get_by_email.go` | `GetByEmail(ctx, email) (*User, error)` | Получение по email |
| `GetByUsername` | `get_by_username.go` | `GetByUsername(ctx, username) (*User, error)` | Получение по username |
| `GetAll` | `get_all.go` | `GetAll(ctx) ([]*User, error)` | Получить всех |

---

## Схема именования

### Формат имён методов

```
{Operation}{By}{Criteria}
```

**Операции:**
- `Get` — получение одного объекта
- `Save` — сохранение (создание или обновление)
- `Delete` — удаление

**Критерии:**
- `ByID` — по идентификатору
- `ByEmail` — по email
- `ByUsername` — по username

### Примеры

```go
// Правильно
repo.Save(ctx, user)
repo.GetByID(ctx, id)
repo.GetByEmail(ctx, email)
repo.GetByUsername(ctx, username)
repo.GetAll(ctx)
repo.DeleteByID(ctx, id)

// Неправильно
repo.GetUserById(ctx, id)      // избыточно "User"
repo.FindUsers(ctx, spec)      // избыточно "Users"
repo.Delete(ctx, id)           // неоднозначно
```

---

## Доменные ошибки

### Базовые ошибки (`internal/errors/errors.go`)

```go
// Общие ошибки
ErrInternal        // внутренняя ошибка
ErrNotFound        // не найдено
ErrAlreadyExists   // уже существует
ErrInvalidArgument // невалидный аргумент

// Ошибки аутентификации
ErrUnauthorized      // не авторизован
ErrForbidden         // доступ запрещён
ErrTokenInvalid      // невалидный токен
ErrTokenExpired      // токен истёк

// Ошибки пользователя
ErrUserNotFound      // пользователь не найден
ErrUserExists        // пользователь уже существует
ErrInvalidCredentials // невалидные учётные данные

// Ошибки пароля
ErrPasswordInvalid   // невалидный пароль
ErrPasswordTooShort  // пароль слишком короткий

// Ошибки email
ErrEmailInvalid      // невалидный email
ErrEmailTooLong      // email слишком длинный

// Ошибки username
ErrUsernameInvalid   // невалидный username
ErrUsernameTooShort  // username слишком короткий
ErrUsernameTooLong   // username слишком длинный

// Ошибки репозитория
ErrRepository        // ошибка репозитория
ErrDBConnection      // ошибка подключения к БД
ErrDBQuery           // ошибка запроса к БД
```

### Использование ошибок

```go
user, err := repo.GetByID(ctx, id)
if err != nil {
    if errors.Is(err, errors.ErrUserNotFound) {
        // Обработка "не найден"
    }
    return err
}

if err := repo.Save(ctx, user); err != nil {
    return err
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

if err := userRepo.Save(ctx, user); err != nil {
    return err
}

// Получение пользователя
user, err := userRepo.GetByID(ctx, userID)
if err != nil {
    return err
}

// Обновление
user.UpdateEmail(newEmail)
if err := userRepo.Save(ctx, user); err != nil {
    return err
}

// Удаление
if err := userRepo.DeleteByID(ctx, userID); err != nil {
    return err
}
```

### Получение всех пользователей

```go
users, err := userRepo.GetAll(ctx)
if err != nil {
    return err
}

for _, user := range users {
    // Обработка пользователя
}
```

---

## Ссылки

- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [DDD Repository](https://domainlanguage.com/ddd/repositories/)
