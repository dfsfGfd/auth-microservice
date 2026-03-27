# Auth Microservice

> Микросервис аутентификации с поддержкой gRPC + REST API (grpc-gateway), JWT access/refresh токенов, PostgreSQL и Redis.

---

## 📋 Оглавление

- [Быстрый старт](#-быстрый-старт)
- [API Endpoints](#-api-endpoints)
- [Архитектура](#-архитектура)
- [Разработка](#-разработка)
- [Docker](#-docker)
- [Документация](#-документация)

---

## 🚀 Быстрый старт

### Требования

- Go 1.25+
- PostgreSQL 15+
- Redis 7+

### Установка

```bash
# 1. Клонировать репозиторий
git clone <repository-url>
cd auth-microservice

# 2. Скопировать .env
cp .env.example .env

# 3. Настроить переменные (обязательно JWT_SECRET, DATABASE_URL, REDIS_URL)
#    edit .env

# 4. Запустить
go run cmd/server/main.go
```

> **Примечание:** Конфигурация загружается только из .env файла. YAML формат не поддерживается.

Сервер запустится на:
- **REST API:** `http://localhost:8080`
- **gRPC:** `localhost:9090`
- **Health:** `http://localhost:8080/health`

---

## 📡 API Endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/auth/register` | Регистрация |
| `POST` | `/api/v1/auth/login` | Вход |
| `POST` | `/api/v1/auth/logout` | Выход |
| `POST` | `/api/v1/auth/refresh` | Обновление токена |
| `GET` | `/health` | Health check |

### Примеры

**Регистрация:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Password123!"}'
```

**Вход:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Password123!"}'
```

**Обновление токена:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<token>"}'
```

**Выход:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{"refresh_token":"<token>"}'
```

📖 **Полная документация API:** [docs/api.md](docs/api.md)  
📦 **Postman коллекция:** [docs/POSTMAN_COLLECTION.md](docs/POSTMAN_COLLECTION.md)

---

## 🏗 Архитектура

### Стек

| Компонент | Технология |
|-----------|------------|
| **RPC** | gRPC + REST (grpc-gateway) |
| **Токены** | JWT (access + refresh) |
| **Кэш** | Redis (refresh токены) |
| **БД** | PostgreSQL (pgx) |
| **DDD** | Domain/Repository/Service/Handler |

### Структура

```
.
├── cmd/
│   ├── server/          # Точка входа сервиса
│   └── migrate/         # Утилита миграций
├── internal/
│   ├── model/           # Domain layer (агрегаты, VO)
│   ├── repository/      # Repository layer (PostgreSQL)
│   ├── service/         # Service layer (бизнес-логика)
│   ├── handler/         # Handler layer (gRPC)
│   ├── cache/           # Redis cache для токенов
│   ├── middleware/      # HTTP/gRPC middleware (rate limiter, logging)
│   ├── di/              # Dependency Injection (Wire)
│   ├── config/          # Конфигурация из .env
│   └── errors/          # Доменные ошибки
├── pkg/                 # Общие пакеты (jwt, bcrypt, logger, db)
├── proto/               # Proto контракты (Buf)
├── api/                 # Swagger/OpenAPI спецификация
├── migrations/          # SQL миграции (golang-migrate)
├── deploy/              # Docker файлы
├── tests/               # Интеграционные тесты
└── docs/                # Документация
```

### Оптимизации (v1.1)

- ✅ Удален неиспользуемый код (DeleteByID, GetByID, list конвертеры)
- ✅ Удалена CORS middleware
- ✅ Исправлен rate limiter для REST API (пути `/api/v1/auth/*`)
- ✅ Исправлена валидация email (trim перед проверкой длины)
- ✅ Упрощена валидация PasswordHash (только проверка длины)

### Время жизни токенов

| Токен | TTL | Хранение |
|-------|-----|----------|
| **Access** | 15 мин | Client (Authorization header) |
| **Refresh** | 14 дней | Redis (`refresh:{token}`) |

---

## 🛠 Разработка

### Команды

```bash
# Форматирование
task format

# Линтинг
task lint

# Генерация Proto
task proto:gen

# Генерация DI
task wire:gen

# Тесты
go test ./... -v
```

### Taskfile

| Команда | Описание |
|---------|----------|
| `task format` | Форматирование Go кода |
| `task lint` | Линтинг Go кода |
| `task tidy` | Очистка зависимостей |
| `task proto:gen` | Генерация Proto (gRPC + REST + Swagger) |
| `task wire:gen` | Генерация DI кода |

### Запуск сервера

| Команда | Описание |
|---------|----------|
| `task server:build` | Сборка бинарного файла |
| `task server:run` | Запуск сервера (локально, требуется `.env`) |
| `task server:dev` | Запуск через `go run` |
| `task server:stop` | Остановка сервера |

### Интеграционные тесты

| Команда | Описание |
|---------|----------|
| `task test:integration:up` | Поднять контейнеры (PostgreSQL, Redis) |
| `task test:integration` | Запустить интеграционные тесты |
| `task test:integration:down` | Удалить контейнеры тестов |

---

## 🐳 Docker

### Быстрый старт

```bash
# 1. Создать .env файл
cp deploy/.env.example deploy/.env

# 2. Запустить через Taskfile
task server:docker:up

# Или напрямую через docker compose
cd deploy && docker compose up -d --build
```

Сервисы:
- **auth-service:** `http://localhost:8080`
- **postgres:** `localhost:5432`
- **redis:** `localhost:6379`

### Остановка

```bash
# Через Taskfile
task server:docker:down

# Или напрямую
docker compose down
```

### Логи

```bash
task server:docker:logs
```

> ⚠️ **Важно:** Docker команды доступны только для развёртывания в production-like окружении.
> Для локальной разработки используйте `task server:dev` или `task server:run`.

📖 **Полное руководство:** [deploy/README.md](deploy/README.md)

---

## 🔒 Безопасность

### Требования к паролю

Согласно рекомендациям **NIST 800-63B**, длина пароля важнее сложности:

- ✅ Минимум **8 символов**
- ✅ Без требований к заглавным/строчным буквам или цифрам
- ✅ Поддержка passphrase (например, `correct horse battery staple`)

**Примеры валидных паролей:**
```
✅ password123      (нет заглавной)
✅ PASSWORD123      (нет строчной)
✅ Password         (нет цифры)
✅ abcdefgh         (только буквы)
✅ MyPass123        (все символы)
```

**Примеры невалидных паролей:**
```
❌ 1234567          (слишком короткий)
❌ Pass1            (слишком короткий)
❌ (пустой)         (пустой пароль)
```

### Rate Limiting

| Endpoint | Лимит/мин |
|----------|-----------|
| Register | 5 |
| Login | 10 |
| Refresh | 30 |
| Logout | 60 |

### Production Checklist

- [ ] `JWT_SECRET` ≥ 32 символов
- [ ] `APP_ENV=production`
- [ ] HTTPS включён
- [ ] Rate limits под нагрузку

---

## 📚 Документация

| Документ | Описание |
|----------|----------|
| [Development Guide](docs/DEVELOPMENT.md) | **Для разработчиков** — архитектура, правила, генерация кода |
| [API Documentation](docs/api.md) | Полное описание API endpoints |
| [Configuration Guide](docs/config.md) | Настройка .env переменных |
| [Docker Guide](deploy/README.md) | Развёртывание в Docker |
| [JWT Package](pkg/jwt/README.md) | JWT сервис документация |
| [Migrations](migrations/README.md) | Управление миграциями БД |

---

## 📝 License

MIT
