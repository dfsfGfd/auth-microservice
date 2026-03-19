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
│   ├── server/           # Точка входа приложения
│   └── migrate/          # Утилита миграций БД
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
│   │   └── cors.go
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
log.Info("user logged in", "user_id", userID, "email", email)
log.Error("database error", "error", err)
log.Warn("rate limit exceeded", "ip", clientIP, "attempts", attempts)

// ❌ Неправильно
log.Info("user logged in")  // нет контекста
log.Error("error", err)     // неясно какая ошибка
```

#### Контекстное логирование

```go
// С request_id для трассировки
reqLog := log.WithRequestID(requestID)
reqLog.Info("processing request", "method", "Login")

// С контекстом
ctx := logger.WithContext(ctx, requestID)
reqLog := log.WithContext(ctx)
reqLog.Info("processing request")
```

#### JSON формат (production)

```json
{
  "level": "info",
  "service": "auth-service",
  "message": "user logged in",
  "user_id": "13f9c5ac-...",
  "email": "user@example.com",
  "time": "2026-03-18T14:43:59Z"
}
```

---

## 🔧 Генерация кода

### Proto (gRPC + REST + Swagger)

```bash
# Генерация из .proto файлов
task proto:gen
```

**Что генерирует:**
- `pkg/proto/auth/v1/` — Go код для gRPC
- `api/auth/v1/` — Swagger/OpenAPI спецификации
- `internal/handler/` — gRPC хендлеры (частично)

### Dependency Injection (Wire)

```bash
# Генерация DI кода
task wire:gen

# Проверка без генерации
task wire:check
```

**Что генерирует:**
- `internal/di/wire_gen.go` — код внедрения зависимостей

### Миграции

```bash
# Применение всех миграций
go run cmd/migrate/main.go -dsn "postgres://..." up

# Откат последней
go run cmd/migrate/main.go -dsn "postgres://..." down
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

### E2E тесты

```bash
# Запуск через docker-compose
task test:e2e

# Только поднять контейнеры
task test:e2e:up

# Запустить тесты против running контейнеров
task test:e2e:run

# Очистка
task test:e2e:down
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

### Разработка

```bash
cd deploy
docker compose up -d --build
```

**Сервисы:**
- `auth-service` (8080, 9090)
- `postgres` (5432)
- `redis` (6379)

### Логи

```bash
# JSON логи
docker compose logs auth-service | grep -v cors

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

- Минимум 8 символов
- 1 заглавная (A-Z)
- 1 строчная (a-z)
- 1 цифра (0-9)

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

---

## 📚 Ссылки

- [API Documentation](docs/api.md)
- [Configuration Guide](docs/config.md)
- [Docker Guide](deploy/README.md)
- [Taskfile](Taskfile.yml)
