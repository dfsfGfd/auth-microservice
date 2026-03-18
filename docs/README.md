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

# 3. Настроить переменные (обязательно JWT_SECRET)
#    edit .env

# 4. Запустить
go run cmd/server/main.go
```

Сервер запустится на:
- **REST API:** `http://localhost:8080`
- **gRPC:** `localhost:9090`
- **Health:** `http://localhost:8080/health`

---

## 📡 API Endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/auth/register` | Регистрация |
| `POST` | `/api/auth/login` | Вход |
| `POST` | `/api/auth/logout` | Выход |
| `POST` | `/api/auth/refresh` | Обновление токена |

### Примеры

**Регистрация:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Password123!"}'
```

**Вход:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Password123!"}'
```

📖 **Полная документация API:** [docs/api.md](docs/api.md)

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
├── cmd/server/           # Точка входа
├── internal/
│   ├── model/           # Domain layer (агрегаты, VO)
│   ├── repository/      # Repository layer (PostgreSQL)
│   ├── service/         # Service layer (бизнес-логика)
│   ├── handler/         # Handler layer (gRPC)
│   ├── cache/           # Кэш (Redis)
│   ├── middleware/      # HTTP/gRPC middleware
│   ├── di/              # Dependency Injection (Wire)
│   ├── config/          # Конфигурация
│   └── errors/          # Доменные ошибки
├── pkg/                 # Общие пакеты (jwt, bcrypt, logger)
├── proto/               # Proto контракты
├── migrations/          # SQL миграции
├── deploy/              # Docker файлы
└── docs/                # Документация
```

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
| `task proto:gen` | Генерация Proto (gRPC + REST + Swagger) |
| `task format` | Форматирование Go кода |
| `task lint` | Линтинг Go кода |
| `task tidy` | Очистка зависимостей |

---

## 🐳 Docker

### Запуск (Development)

```bash
cd deploy
docker compose up -d --build
```

Сервисы:
- **auth-service:** `http://localhost:8080`
- **postgres:** `localhost:5432`
- **redis:** `localhost:6379`

### Остановка

```bash
docker compose down
```

📖 **Полное руководство:** [deploy/DEPLOY.md](deploy/DEPLOY.md)

---

## 🔒 Безопасность

### Требования к паролю

- Минимум 8 символов
- 1 заглавная буква (A-Z)
- 1 строчная буква (a-z)
- 1 цифра (0-9)

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
- [ ] CORS настроен для ваших доменов
- [ ] Rate limits под нагрузку

---

## 📚 Документация

| Документ | Описание |
|----------|----------|
| [API Documentation](docs/api.md) | Полное описание API endpoints |
| [Configuration Guide](docs/config.md) | Настройка .env переменных |
| [Docker Guide](deploy/DEPLOY.md) | Развёртывание в Docker |
| [JWT Package](pkg/jwt/README.md) | JWT сервис документация |
| [Migrations](migrations/README.md) | Управление миграциями БД |

---

## 📝 License

MIT
