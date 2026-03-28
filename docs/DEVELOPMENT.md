# Development Guide

Руководство для разработчиков auth-microservice.

---

## 📋 Оглавление

- [Архитектура](#-архитектура)
- [Структура проекта](#-структура-проекта)
- [Правила разработки](#-правила-разработки)
- [Генерация кода](#-генерация-кода)
- [Тестирование](#-тестирование)
- [Логирование](#-логирование)
- [Обработка ошибок](#-обработка-ошибок)
- [Docker](#-docker)

---

## 🏗 Архитектура

### Слои (DDD)

```
┌─────────────────────────────────────────┐
│           Handler (gRPC/REST)           │  ← Входные запросы
├─────────────────────────────────────────┤
│              Service Layer              │  ← Бизнес-логика
├─────────────────────────────────────────┤
│         Repository Layer                │  ← Доступ к данным
├─────────────────────────────────────────┤
│          Domain Layer (Model)           │  ← Бизнес-объекты
└─────────────────────────────────────────┘
```

### Поток запроса

```
HTTP/gRPC → Handler → Service → Repository → PostgreSQL
                         ↓
                    Cache (Redis)
```

---

## 📁 Структура проекта

```
auth-microservice/
├── cmd/
│   └── server/           # Точка входа приложения
│
├── internal/             # Приватный код (не экспортируется)
│   ├── model/            # Domain layer
│   │   ├── account.go    # Account агрегат
│   │   ├── email.go      # Email VO
│   │   └── password.go   # Password VO
│   │
│   ├── repository/       # Repository layer
│   │   ├── repository.go # Интерфейсы
│   │   ├── model/        # DB модели
│   │   ├── converter/    # Domain ↔ DB конвертеры
│   │   └── auth/         # PostgreSQL реализация
│   │
│   ├── service/          # Service layer (бизнес-логика)
│   │   ├── service.go    # Интерфейсы
│   │   └── auth/         # Реализация
│   │
│   ├── handler/          # Handler layer (gRPC)
│   │   └── auth/         # gRPC хендлеры
│   │
│   ├── cache/            # Кэш слой
│   │   └── token/        # Redis cache для токенов
│   │
│   ├── middleware/       # HTTP/gRPC middleware
│   │   ├── rate_limiter.go
│   │   └── logging.go
│   │
│   ├── di/               # Dependency Injection (Wire)
│   ├── config/           # Конфигурация из .env
│   └── errors/           # Доменные ошибки
│
├── pkg/                  # Публичные пакеты
│   ├── jwt/              # JWT сервис
│   ├── bcrypt/           # Хеширование паролей
│   ├── logger/           # Логирование
│   └── db/               # DB подключения
│       ├── postgresql/
│       └── redisdb/
│
├── proto/                # Proto контракты
│   └── auth/v1/
│       └── auth.proto
│
├── api/                  # Swagger спецификации
├── migrations/           # SQL миграции
├── deploy/               # Docker файлы
├── tests/                # E2E тесты
└── docs/                 # Документация
```

---

## 📝 Правила разработки

### 1. Структура пакетов

**internal/** — код, который нельзя импортировать из других модулей

**pkg/** — публичные пакеты, можно импортировать

```go
// ✅ Правильно
import "auth-microservice/pkg/logger"
import "auth-microservice/internal/model"

// ❌ Нельзя (из другого модуля)
import "auth-microservice/internal/service"
```

### 2. Именование

#### Файлы
- Один файл = один метод (в repository/)
- snake_case для файлов
- Группировка по функциональности

```
internal/repository/auth/
├── repository.go       # Конструктор
├── save.go             # Save метод
├── get_by_id.go        # GetByID метод
├── get_by_email.go     # GetByEmail метод
└── delete_by_id.go     # DeleteByID метод
```

#### Функции и методы
- CamelCase для экспортируемых
- camelCase для приватных
- Глаголы для действий (Get, Save, Delete)

```go
// ✅ Правильно
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Account, error)
func (r *Repository) save(ctx context.Context, account *Account) error

// ❌ Неправильно
func (r *Repository) getAccountById(...)  // избыточно
func (r *Repository) AccountSave(...)     // неясно
```

### 3. Обработка ошибок

#### Доменные ошибки в `internal/errors/`

```go
// internal/errors/errors.go
var (
    ErrAccountNotFound   = errors.New("account not found")
    ErrInvalidEmail      = errors.New("invalid email")
    ErrInvalidPassword   = errors.New("invalid password")
)
```

#### Использование

```go
// ✅ Правильно
account, err := r.repo.GetByID(ctx, id)
if err != nil {
    if errors.Is(err, errors.ErrAccountNotFound) {
        return nil, status.Error(codes.NotFound, "account not found")
    }
    return nil, status.Error(codes.Internal, "internal error")
}

// ❌ Неправильно
if err == errors.ErrAccountNotFound {  // не работает с wrapped errors
    return nil, status.Error(codes.NotFound, "account not found")
}
```

#### Конвертация ошибок (service → handler)

```go
// В handler
resp, err := h.service.Login(ctx, req.Email, req.Password)
if err != nil {
    switch {
    case errors.Is(err, errors.ErrInvalidCredentials):
        return nil, status.Error(codes.Unauthenticated, "invalid credentials")
    case errors.Is(err, errors.ErrAccountNotFound):
        return nil, status.Error(codes.NotFound, "account not found")
    default:
        return nil, status.Error(codes.Internal, "internal error")
    }
}
```

### 4. Логирование

#### Уровни

| Уровень | Когда использовать |
|---------|-------------------|
| `debug` | Отладка (SQL, детали выполнения) |
| `info` | Нормальная работа (старт, успешные операции) |
| `warn` | Предупреждения (rate limit, временные проблемы) |
| `error` | Ошибки (сбои БД, валидация) |
| `fatal` | Критические (невозможно запустить сервер) |

#### Формат

```go
// ✅ Правильно
log.Info("login", "user_id", userID, "status", 200, "dur_ms", duration)
log.Error("db_error", "err", err, "dur_ms", duration)
log.Warn("rate_limit", "ip", clientIP, "path", r.URL.Path)

// ❌ Неправильно
log.Info("user logged in")  // нет контекста
log.Error("error", err)     // неясно какая ошибка
```

#### Контекстное логирование

```go
// С request_id для трассировки
reqLog := log.WithRequestID(requestID)
reqLog.Info("login", "user_id", userID)

// С контекстом
ctx := logger.WithContext(ctx, requestID)
reqLog := log.WithContext(ctx)
reqLog.Info("refresh_token", "user_id", userID)
```

#### JSON формат (production, оптимизированный)

```json
{
  "ts": "2026-03-21T15:00:00Z",
  "lvl": "info",
  "msg": "login",
  "srv": "auth-service",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": 200,
  "dur_ms": 45
}
```

**Оптимизированные имена полей:**

| Поле | Описание |
|------|----------|
| `ts` | Timestamp (ISO 8601) |
| `lvl` | Level (info/warn/error) |
| `msg` | Message |
| `srv` | Service name |
| `user_id` | User ID |
| `rid` | Request ID |
| `trace` | Trace ID |
| `dur_ms` | Duration (миллисекунды) |
| `err` | Error |
| `status` | Status code |

---

## 🔧 Генерация кода

### Taskfile команды

```bash
# Форматирование
task format

# Линтинг
task lint

# Генерация кода
task proto:gen
task wire:gen

# Миграции
task migrate:install   # Установить golang-migrate CLI (один раз)
task migrate:up        # Применить все миграции
task migrate:down      # Откатить последнюю миграцию
task migrate:status    # Показать статус миграций
task migrate:force     # Принудительно установить версию

# Запуск сервера
task server:build    # Сборка
task server:run      # Запуск (локально)
task server:dev      # go run
task server:stop     # Остановка

# Интеграционные тесты
task test:integration:up     # Поднять контейнеры
task test:integration        # Запустить тесты
task test:integration:down   # Удалить контейнеры
```

---

## 🧪 Тестирование

### Unit тесты

```bash
# Запуск всех тестов
go test ./... -v

# Тесты конкретного пакета
go test ./pkg/jwt/... -v
```

### Интеграционные тесты

```bash
# 1. Поднять контейнеры (PostgreSQL, Redis)
task test:integration:up

# 2. Запустить тесты
task test:integration

# 3. Удалить контейнеры
task test:integration:down
```

Или напрямую:
```bash
go test -tags integration -v -timeout 5m ./tests/...
```

### Структура теста

```go
// tests/e2e_test.go
//go:build integration

func TestRegisterAndLogin(t *testing.T) {
    // 1. Регистрация
    resp := registerUser(t, email, password)
    require.Equal(t, 200, resp.StatusCode)

    // 2. Логин
    tokens := loginUser(t, email, password)
    require.NotEmpty(t, tokens.AccessToken)

    // 3. Обновление токена
    newTokens := refreshToken(t, tokens.RefreshToken)
    require.NotEmpty(t, newTokens.AccessToken)
}
```

---

## 🐳 Docker

### Локальная разработка

```bash
cd deploy
docker compose up -d --build
```

**Сервисы:**
- `auth-service` (8080, 9090)
- `postgres` (5432)
- `redis` (6379)

### Миграции

Для управления миграциями используется **golang-migrate CLI**.

**Установка (локально):**

```bash
task migrate:install
```

**Команды:**

```bash
# Применить все миграции
task migrate:up

# Откатить последнюю миграцию
task migrate:down

# Показать статус миграций
task migrate:status

# Принудительно установить версию
task migrate:force VERSION=3
```

**Напрямую через CLI:**

```bash
# Установить golang-migrate
go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1

# Применить миграции
./bin/migrate -path migrations -database "postgres://user:pass@localhost:5432/auth?sslmode=disable" up

# Откатить
./bin/migrate -path migrations -database "postgres://..." down

# Статус
./bin/migrate -path migrations -database "postgres://..." status
```

**В Docker:**

Миграции применяются **автоматически** при старте сервера. golang-migrate отслеживает применённые миграции в таблице `schema_migrations`, поэтому повторный запуск не применяет их заново.

```yaml
# deploy/docker-compose.yml
volumes:
  - postgres_data:/var/lib/postgresql/data  # Сохранение данных
```

**Важно:** Миграция проходит **один раз** при первом запуске. При повторном запуске (после остановки) миграция **не применяется**, данные сохраняются в volume.

📖 **Полное руководство:** [migrations/README.md](../migrations/README.md)

### Логи

```bash
# JSON логи
docker compose logs auth-service

# В реальном времени
docker compose logs -f auth-service
```

### Отладка

```bash
# Вход в контейнер
docker compose exec auth-service sh

# Проверка переменных окружения
docker compose exec auth-service env

# Подключение к PostgreSQL
docker compose exec postgres psql -U auth_user -d auth
```

---

## 📊 Конфигурация

### Переменные окружения (.env)

```bash
# Обязательно
DATABASE_URL=postgres://user:pass@localhost:5432/auth?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key-minimum-32-characters-long

# Опционально
LOG_LEVEL=info
LOG_FORMAT=json
APP_ENV=development
```

### Загрузка конфигурации

```go
// internal/di/provider.go
func loadConfig() (*config.Config, error) {
    // 1. Сначала .env (приоритет)
    cfg, err := config.LoadFromEnv()
    if err == nil {
        return cfg, nil
    }
    
    // 2. Потом config.yaml
    cfg, err = config.Load("config.yaml")
    if err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }
    
    return cfg, nil
}
```

---

## 🔒 Безопасность

### Пароли

Согласно **NIST 800-63B**, длина важнее сложности:

- ✅ Минимум **8 символов**
- ✅ Без требований к заглавным/строчным буквам или цифрам
- ✅ Поддержка passphrase

**Примеры:**
```
✅ password123    (валиден)
✅ PASSWORD123    (валиден)
✅ Password       (валиден)
❌ 1234567        (слишком короткий)
```

### JWT токены

| Токен | TTL | Хранение |
|-------|-----|----------|
| Access | 15 мин | Client (Authorization header) |
| Refresh | 14 дней | Redis (`refresh:{token}`) |

### Rate Limiting

| Endpoint | Лимит/мин |
|----------|-----------|
| Register | 5 |
| Login | 10 |
| Refresh | 30 |
| Logout | 60 |

### Bcrypt Cost

Стоимость хеширования: **12** (OWASP 2024 рекомендация)

```go
// pkg/bcrypt/bcrypt.go
const BCryptCost = 12
```

---

## 📚 Ссылки

- [API Documentation](docs/api.md)
- [Configuration Guide](docs/config.md)
- [Docker Guide](deploy/README.md)
- [Taskfile](Taskfile.yml)
