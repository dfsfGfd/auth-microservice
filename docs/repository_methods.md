# Repository Methods — DDD Implementation Guide

Полный список методов для реализации репозиториев в соответствии с DDD паттернами.

---

## 📋 Оглавление

- [Структура пакета](#структура-пакета)
- [AccountRepository методы](#accountrepository-методы)
- [Схема именования](#схема-именования)
- [Доменные ошибки](#доменные-ошибки)

---

## Структура пакета

```
internal/repository/
├── errors/
│   └── errors.go           # Доменные ошибки
├── model/
│   └── account.go          # DB модели
├── converter/
│   └── account.go          # Конвертеры domain ↔ DB
├── auth/
│   ├── repository.go       # Конструктор
│   ├── save.go             # Save метод
│   ├── delete_by_id.go     # DeleteByID метод
│   ├── get_by_id.go        # GetByID метод
│   ├── get_by_email.go     # GetByEmail метод
│   └── get_all.go          # GetAll метод
└── repository.go           # Интерфейс AccountRepository
```

---

## AccountRepository методы

| Метод | Файл | Сигнатура | Описание |
|-------|------|-----------|----------|
| `Save` | `save.go` | `Save(ctx, account) error` | Сохранение (создание или обновление) |
| `DeleteByID` | `delete_by_id.go` | `DeleteByID(ctx, id) error` | Удаление по ID |
| `GetByID` | `get_by_id.go` | `GetByID(ctx, id) (*Account, error)` | Получение по ID |
| `GetByEmail` | `get_by_email.go` | `GetByEmail(ctx, email) (*Account, error)` | Получение по email |
| `GetAll` | `get_all.go` | `GetAll(ctx) ([]*Account, error)` | Получить все |

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

### Примеры

```go
// Правильно
repo.Save(ctx, account)
repo.GetByID(ctx, id)
repo.GetByEmail(ctx, email)
repo.GetAll(ctx)
repo.DeleteByID(ctx, id)

// Неправильно
repo.GetAccountById(ctx, id)  // избыточно "Account"
repo.FindAccounts(ctx, spec)  // избыточно "Accounts"
repo.Delete(ctx, id)          // неоднозначно
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

// Ошибки аккаунта
ErrAccountNotFound   // аккаунт не найден
ErrAccountExists     // аккаунт уже существует
ErrInvalidCredentials // невалидные учётные данные

// Ошибки пароля
ErrPasswordInvalid   // невалидный пароль
ErrPasswordTooShort  // пароль слишком короткий

// Ошибки email
ErrEmailInvalid      // невалидный email
ErrEmailTooLong      // email слишком длинный

// Ошибки репозитория
ErrRepository        // ошибка репозитория
ErrDBConnection      // ошибка подключения к БД
ErrDBQuery           // ошибка запроса к БД
```

### Использование ошибок

```go
account, err := repo.GetByID(ctx, id)
if err != nil {
    if errors.Is(err, errors.ErrAccountNotFound) {
        // Обработка "не найден"
    }
    return err
}

if err := repo.Save(ctx, account); err != nil {
    return err
}
```

---

## Примеры использования

### Базовое использование

```go
// Создание аккаунта
account, err := model.NewAccount(email, password)
if err != nil {
    return err
}

if err := accountRepo.Save(ctx, account); err != nil {
    return err
}

// Получение аккаунта
account, err := accountRepo.GetByID(ctx, accountID)
if err != nil {
    return err
}

// Обновление
account.UpdateEmail(newEmail)
if err := accountRepo.Save(ctx, account); err != nil {
    return err
}

// Удаление
if err := accountRepo.DeleteByID(ctx, accountID); err != nil {
    return err
}
```

### Получение всех аккаунтов

```go
accounts, err := accountRepo.GetAll(ctx)
if err != nil {
    return err
}

for _, account := range accounts {
    // Обработка аккаунта
}
```

---

## Ссылки

- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [DDD Repository](https://domainlanguage.com/ddd/repositories/)
