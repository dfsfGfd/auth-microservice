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
| `cmd/server/main.go` | Точка входа сервиса |
| `cmd/migrate/main.go` | Утилита миграций БД |
| `deploy/docker-compose.yml` | Docker окружение |
| `.env.example` | Шаблон переменных окружения |
| `proto/auth/v1/auth.proto` | gRPC/REST API контракт |
| `Taskfile.yml` | Build автоматизация |

---

## 🚀 Команды

### Taskfile (рекомендуется)

```bash
# Форматирование
task format

# Линтинг
task lint

# Запуск сервера
task server:build    # Сборка
task server:run      # Запуск (локально)
task server:dev      # go run
task server:stop     # Остановка

# Интеграционные тесты
task test:integration:up     # Поднять контейнеры
task test:integration        # Запустить тесты
task test:integration:down   # Удалить контейнеры

# Генерация кода
task proto:gen       # Proto → Go
task wire:gen        # DI код
```

### Docker

```bash
# Запуск в Docker
cd deploy && docker compose up -d --build

# Логи
docker compose logs -f auth-service

# Остановка
docker compose down
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

---

## 🔄 История изменений

### v1.3 Remove CORS

**Удалено:**
- CORS middleware (больше не нужен)
- CORS конфигурация из config
- Зависимость `github.com/rs/cors`

### v1.2 Config (последний коммит)

**Конфигурация:**
- Загрузка конфига только из .env файла (YAML удален)
- `config.Load()` - единая точка входа
- Упрощен DI provider

### v1.1 Refactor

**Исправленные баги:**
- Rate limiter для REST API теперь работает с правильными путями `/api/v1/auth/*`
- Email валидация: trim теперь выполняется перед проверкой длины

**Оптимизации:**
- Удален неиспользуемый код: `DeleteByID`, `GetByID` в репозитории
- Удалены неиспользуемые list конвертеры
- Упрощена PasswordHash валидация (убран хардкод bcrypt префиксов)
- Упрощен `RefreshTTLDuration()` (убрана лишняя обертка ошибки)
