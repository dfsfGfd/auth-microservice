# План рефакторинга: User → Account, удаление Username

## Обзор изменений

**Цель:** Переименовать сущность `User` в `Account` и удалить поле `Username` из доменной модели.

**Причина:** Упрощение модели данных — аутентификация только по email.

---

## Изменяемые файлы по слоям

### 1. Domain Layer (`internal/model/`)

| Файл | Действие | Описание |
|------|----------|----------|
| `internal/model/user.go` | ✏️ Переименовать → `account.go` | Rename struct `User` → `Account`, удалить поле `Username` |
| `internal/model/username.go` | ❌ Удалить | Value Object больше не нужен |
| `internal/model/email.go` | ✅ Оставить | Без изменений |
| `internal/model/password.go` | ✅ Оставить | Без изменений |
| `internal/model/password_hash.go` | ✅ Оставить | Без изменений |

**Изменения в `account.go`:**
```go
type Account struct {
    ID        AccountID
    Email     Email
    Password  PasswordHash
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

### 2. Repository Layer (`internal/repository/`)

#### 2.1 Интерфейс (`internal/repository/repository.go`)

✏️ Изменить:
- `UserRepository` → `AccountRepository`
- `Save(ctx, *model.User)` → `Save(ctx, *model.Account)`
- `GetByID(ctx, id)` → возвращает `*model.Account`
- `GetByEmail(ctx, email)` → возвращает `*model.Account`
- ❌ Удалить метод `GetByUsername`

#### 2.2 Реализация (`internal/repository/auth/`)

| Файл | Действие | Описание |
|------|----------|----------|
| `repository.go` | ✏️ Rename struct `UserRepository` → `AccountRepository` | Конструктор `NewAccountRepository` |
| `save.go` | ✏️ Изменить сигнатуру | `Save(ctx, *model.Account)` |
| `delete_by_id.go` | ✅ Без изменений | Логика та же |
| `get_by_id.go` | ✏️ Изменить return type | `*model.Account` |
| `get_by_email.go` | ✅ Без изменений | Логика та же |
| `get_by_username.go` | ❌ Удалить | Метод больше не нужен |
| `get_all.go` | ✏️ Изменить return type | `[]*model.Account` |

#### 2.3 Модель БД (`internal/repository/model/`)

| Файл | Действие | Описание |
|------|----------|----------|
| `user.go` | ✏️ Переименовать → `account.go` | Struct `UserDB` → `AccountDB`, удалить поле `Username` |

**Изменения в `account.go`:**
```go
type AccountDB struct {
    ID        uuid.UUID `db:"id"`
    Email     string    `db:"email"`
    Password  string    `db:"password"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}
```

#### 2.4 Конвертер (`internal/repository/converter/`)

| Файл | Действие | Описание |
|------|----------|----------|
| `user.go` | ✏️ Переименовать → `account.go` | Функции `UserToDB` → `AccountToDB`, `UserToDomain` → `AccountToDomain` |

**Изменения:**
```go
func AccountToDomain(acc *model.AccountDB) *model.Account
func AccountToDB(acc *model.Account) *model.AccountDB
```

---

### 3. DI Layer (`internal/di/`)

| Файл | Действие | Описание |
|------|----------|----------|
| `wire.go` | ✏️ Изменить | `UserRepository` → `AccountRepository`, `ProvideUserRepository` → `ProvideAccountRepository` |
| `wire_gen.go` | 🔄 Перегенерировать | `go generate ./...` |

---

### 4. Proto Layer (`pkg/proto/` / `proto/`)

| Файл | Действие | Описание |
|------|----------|----------|
| `proto/auth.proto` | ✏️ Изменить | Rename `User` → `Account`, удалить поле `username` |
| `pkg/proto/auth.pb.go` | 🔄 Перегенерировать | `buf generate` / `protoc` |
| `pkg/proto/auth_grpc.pb.go` | 🔄 Перегенерировать | `buf generate` / `protoc` |
| `pkg/proto/auth.pb.gw.go` | 🔄 Перегенерировать | `buf generate` / `protoc` |

**Изменения в `auth.proto`:**
```protobuf
message Account {
    string id = 1;
    string email = 2;
    // username удалён
}
```

---

### 5. Service Layer (`internal/service/`) — будущий слой

✏️ В реализации service layer использовать:
- `AccountService` вместо `UserService`
- Методы без упоминания username (только email + password)

---

### 6. Handler Layer (`internal/handler/`) — будущий слой

✏️ В реализации handler layer:
- GRPC handlers: `AccountServer` вместо `UserServer`
- REST mappings: обновить JSON поля (без `username`)

---

### 7. Errors (`internal/errors/`)

| Файл | Действие | Описание |
|------|----------|----------|
| `errors.go` | ✏️ Rename | `ErrUserNotFound` → `ErrAccountNotFound` |
| `repository.go` | ✏️ Rename | `ErrRepository` → оставить, но изменить контекст |

---

### 8. Database Migration

⚠️ **Требуется миграция БД:**

```sql
-- Вариант 1: Переименовать таблицу и удалить колонку
ALTER TABLE users RENAME TO accounts;
ALTER TABLE accounts DROP COLUMN username;

-- Вариант 2: Создать новую таблицу
CREATE TABLE accounts (...);
INSERT INTO accounts (id, email, password, created_at, updated_at)
SELECT id, email, password, created_at, updated_at FROM users;
DROP TABLE users;
```

---

### 9. Config (`internal/config/`)

✅ Без изменений (если нет специфичных ссылок на User)

---

### 10. Docs (`docs/`)

| Файл | Действие | Описание |
|------|----------|----------|
| `README.md` | ✏️ Обновить | User → Account |
| `api.md` | ✏️ Обновить | User → Account, удалить username из примеров |
| `repository_methods.md` | ✏️ Обновить | UserRepository → AccountRepository, удалить GetByUsername |

---

## Порядок выполнения

### Этап 1: Domain Layer
1. ✏️ Переименовать `user.go` → `account.go`
2. ✏️ Удалить `username.go`
3. ✅ Обновить импорты в других файлах

### Этап 2: Repository Layer
1. ✏️ Обновить интерфейс `repository.go`
2. ✏️ Переименовать файлы реализации
3. ✏️ Обновить модель и конвертер
4. ❌ Удалить `get_by_username.go`

### Этап 3: DI Layer
1. ✏️ Обновить `wire.go`
2. 🔄 Перегенерировать `wire_gen.go`

### Этап 4: Proto Layer
1. ✏️ Обновить `auth.proto`
2. 🔄 Перегенерировать Go-файлы

### Этап 5: Errors
1. ✏️ Переименовать ошибки

### Этап 6: Database Migration
1. ⚠️ Создать миграционный SQL-скрипт

### Этап 7: Docs
1. ✏️ Обновить документацию

### Этап 8: Verification
1. ✅ `go build ./...`
2. ✅ `go test ./...`
3. ✅ Проверка компиляции proto

---

## Чек-лист

- [x] Domain: `account.go` создан, `username.go` удалён
- [x] Repository: интерфейс обновлён
- [x] Repository: реализация обновлена
- [x] Repository: модель БД обновлена
- [x] Repository: конвертер обновлён
- [x] DI: `wire.go` обновлён, `wire_gen.go` перегенерирован
- [x] Proto: `.proto` файл обновлён
- [x] Proto: Go-файлы перегенерированы
- [x] Errors: переименованы
- [x] DB: миграция создана
- [x] Docs: обновлены
- [x] Build: `go build ./...` ✅
- [x] Tests: `go test ./...` ✅ (кроме logger_test.go — не связан с рефакторингом)

---

## Примечания

- После рефакторинга **аутентификация будет только по email**
- Refresh token логика остаётся без изменений
- JWT claims не должны содержать username
