# Integration Tests

Интеграционные тесты для auth-microservice с использованием testcontainers.

## 📁 Структура

```
tests/
├── integration_test.go      # Основные интеграционные тесты
├── database_test.go         # Тесты базы данных и Redis
└── README.md                # Этот файл
```

## 🚀 Запуск тестов

### Все тесты

```bash
# Из корня проекта
go test ./tests -v
```

### Только интеграционные тесты

```bash
# С тегом integration
go test -tags integration ./tests -v
```

### Отдельный тест

```bash
# Конкретный тест
go test ./tests -run TestDatabaseConnection -v

# Тесты базы данных
go test ./tests -run TestDatabase -v

# Тесты Redis
go test ./tests -run TestRedis -v
```

### С очисткой контейнеров

```bash
# Автоматическая очистка после тестов
go test ./tests -v -cleanup
```

### С выводом логов контейнеров

```bash
# Вывод логов PostgreSQL и Redis
go test ./tests -v -args -container-logs
```

## 📦 Зависимости

### Требования

- **Docker** — для запуска testcontainers
- **Go 1.26+** — версия Go проекта
- **testcontainers-go** — библиотека для управления контейнерами

### Установка зависимостей

```bash
# Автоматическая установка через go mod
go mod tidy

# Проверка Docker
docker --version
docker ps
```

## 🏗 Архитектура тестов

### TestMain

`TestMain` запускается перед всеми тестами:

1. **Создаёт контейнеры:**
   - PostgreSQL 18-alpine
   - Redis 7.4-alpine

2. **Применяет миграции:**
   - Запускает `cmd/migrate` для создания таблиц

3. **Устанавливает переменные окружения:**
   - `DATABASE_URL`
   - `REDIS_URL`
   - `JWT_SECRET`
   - `APP_ENV=test`

4. **Очищает контейнеры** после всех тестов

### Контейнеры

#### PostgreSQL

```go
postgres.Run(ctx,
    "postgres:18-alpine",
    postgres.WithDatabase("auth_test"),
    postgres.WithUsername("test_user"),
    postgres.WithPassword("test_password"),
)
```

**Параметры:**
- Image: `postgres:18-alpine`
- Database: `auth_test`
- User: `test_user`
- Password: `test_password`
- Wait strategy: log "database system is ready"

#### Redis

```go
redis.Run(ctx,
    "redis:7.4-alpine",
)
```

**Параметры:**
- Image: `redis:7.4-alpine`
- Wait strategy: log "Ready to accept connections"

## 📝 Тесты

### Database Tests (`database_test.go`)

| Тест | Описание |
|------|----------|
| `TestDatabaseConnection` | Проверка подключения к PostgreSQL |
| `TestRedisConnection` | Проверка подключения к Redis |
| `TestDatabaseTables` | Проверка существования таблиц |
| `TestDatabaseIndexes` | Проверка индексов |
| `TestDatabaseTrigger` | Проверка триггеров |
| `TestDatabaseFunction` | Проверка функций |
| `TestRedisOperations` | Тест операций Redis (SET/GET/DELETE) |
| `TestConcurrentDatabaseAccess` | Конкурентный доступ к БД |
| `TestConcurrentRedisAccess` | Конкурентный доступ к Redis |

### Integration Tests (`integration_test.go`)

| Тест | Описание |
|------|----------|
| `TestIntegration_Registration` | Регистрация пользователя |
| `TestIntegration_Login` | Вход пользователя |
| `TestIntegration_TokenRefresh` | Обновление токенов |
| `TestIntegration_Logout` | Выход пользователя |
| `TestIntegration_RateLimiting` | Rate limiting |

## 🔧 Отладка

### Вывод логов контейнеров

```bash
# Включить логи в TestMain
fmt.Println(container.Logs(ctx))
```

### Сохранение контейнеров после тестов

```bash
# Не удалять контейнеры (для отладки)
export TESTCONTAINERS_RYUK_DISABLED=true
go test ./tests -v
```

### Подключение к контейнеру во время тестов

```bash
# Найти контейнер
docker ps | grep auth

# Подключиться к PostgreSQL
docker exec -it <container_id> psql -U test_user -d auth_test

# Подключиться к Redis
docker exec -it <container_id> redis-cli
```

## 🐛 Troubleshooting

### Ошибка: "Cannot connect to the Docker daemon"

```bash
# Проверить Docker
docker ps

# Если не работает, запустить Docker
sudo systemctl start docker
```

### Ошибка: "Container failed to start"

```bash
# Проверить логи Docker
docker logs <container_id>

# Увеличить таймаут
export TEST_STARTUP_TIMEOUT=120s
go test ./tests -v
```

### Ошибка: "Connection refused"

```bash
# Проверить, что контейнер запущен
docker ps

# Проверить порты
docker port <container_id>
```

### Ошибка: "No space left on device"

```bash
# Очистить Docker
docker system prune -a

# Очистить testcontainers
docker volume prune
```

## 📊 Покрытие

### Проверка покрытия

```bash
# Запуск с покрытием
go test -tags integration ./tests -coverprofile=coverage.out

# Просмотр покрытия
go tool cover -html=coverage.out
```

### Целевое покрытие

- Database tests: 80%+
- Integration tests: 70%+
- Overall: 75%+

## 📚 Ссылки

- [testcontainers-go documentation](https://golang.testcontainers.org/)
- [PostgreSQL testcontainer](https://golang.testcontainers.org/modules/postgres/)
- [Redis testcontainer](https://golang.testcontainers.org/modules/redis/)
- [Go testing package](https://pkg.go.dev/testing)
