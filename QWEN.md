# Auth Microservice — Project Context

> Go микросервис аутентификации с gRPC + REST (grpc-gateway), JWT токенами, PostgreSQL и Redis.

---

## 🏗 Архитектура

**Стек:**
- Go 1.26
- gRPC + REST (grpc-gateway)
- PostgreSQL (pgx)
- Redis
- JWT (access + refresh)
- DDD (Domain/Repository/Service/Handler)

**Структура:**
```
internal/
├── model/          # Domain layer (агрегаты, VO)
├── repository/     # Repository layer (PostgreSQL)
├── service/        # Service layer (бизнес-логика)
├── handler/        # Handler layer (gRPC)
├── cache/          # Redis cache для токенов
├── middleware/     # HTTP/gRPC middleware
├── di/             # Google Wire DI
├── config/         # Конфигурация из .env
└── errors/         # Доменные ошибки
```

---

## 📁 Важные файлы

| Файл | Описание |
|------|----------|
| `cmd/server/main.go` | Точка входа |
| `cmd/migrate/main.go` | Утилита миграций |
| `deploy/docker-compose.yml` | Docker окружение |
| `.env` | Переменные окружения |

---

## 🚀 Команды

```bash
# Запуск в Docker
cd deploy && docker compose up -d --build

# Локальный запуск
go run cmd/server/main.go

# Форматирование
task format

# Линтинг
task lint

# Генерация Proto
task proto:gen

# Генерация DI
task wire:gen
```

---

## 📚 Документация

- [docs/README.md](docs/README.md) — основная документация
- [docs/api.md](docs/api.md) — API endpoints
- [docs/config.md](docs/config.md) — настройка .env
- [deploy/README.md](deploy/README.md) — Docker guide

---

## GitHub

Username: `dfsfGfd`
