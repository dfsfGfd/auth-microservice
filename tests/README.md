# Integration Tests

Интеграционные тесты для auth-microservice с использованием Docker Compose.

## 📁 Структура

```
tests/
├── docker-compose.integration.yml  # Docker Compose для тестов
├── integration_api_test.go         # API интеграционные тесты
├── database_test.go                # Тесты базы данных и Redis
└── README.md                       # Этот файл
```

## 🚀 Быстрый старт

### 1. Поднять инфраструктуру

```bash
cd tests

# Запуск всех сервисов (auth-service + postgres + redis)
docker-compose -f docker-compose.integration.yml up -d --build

# Проверка статуса
docker-compose -f docker-compose.integration.yml ps
```

### 2. Запустить тесты

```bash
# Все интеграционные тесты
go test -tags integration -v ./tests

# Отдельный тест
go test -tags integration -v ./tests -run TestIntegration_HealthCheck

# Тесты регистрации
go test -tags integration -v ./tests -run TestIntegration_Registration

# Тесты входа
go test -tags integration -v ./tests -run TestIntegration_Login
```

### 3. Очистить

```bash
# Остановка и удаление контейнеров
docker-compose -f docker-compose.integration.yml down -v

# Удаление образов
docker-compose -f docker-compose.integration.yml down -v --rmi all
```

## 📦 Зависимости

### Требования

- **Docker** + **Docker Compose** — для запуска инфраструктуры
- **Go 1.26+** — версия Go проекта

### Проверка

```bash
# Проверка Docker
docker --version
docker compose version

# Проверка Go
go version
```

## 🏗 Архитектура тестов

### Инфраструктура

`docker-compose.integration.yml` поднимает:

1. **auth-service** — тестируемый микросервис (порт 8080)
2. **postgres:18-alpine** — база данных (порт 5432)
3. **redis:7.4-alpine** — кэш для токенов (порт 6379)
4. **migrate** — утилита для применения миграций

```yaml
services:
  auth-service:     # Тестируемый сервис
    ports: [8080, 9090]
    depends_on: [postgres, redis, migrate]
  
  migrate:          # Применяет миграции
    depends_on: [postgres]
  
  postgres:         # База данных
    ports: [5432]
  
  redis:            # Кэш
    ports: [6379]
```

### Переменные окружения

```bash
APP_ENV=test
DATABASE_URL=postgres://test_user:test_password@postgres:5432/auth_test
REDIS_URL=redis://redis:6379
JWT_SECRET=test-secret-key-minimum-32-characters-long
```

## 📝 Тесты

### API Integration Tests (`integration_api_test.go`)

| Тест | Описание |
|------|----------|
| `TestIntegration_HealthCheck` | Проверка health check endpoint |
| `TestIntegration_Registration_Success` | Успешная регистрация |
| `TestIntegration_Registration_DuplicateEmail` | Регистрация с существующим email |
| `TestIntegration_Registration_InvalidEmail` | Регистрация с невалидным email |
| `TestIntegration_Registration_WeakPassword` | Регистрация со слабым паролем |
| `TestIntegration_Login_Success` | Успешный вход |
| `TestIntegration_Login_InvalidCredentials` | Вход с неверными данными |
| `TestIntegration_TokenRefresh_Success` | Обновление токена |
| `TestIntegration_Logout_Success` | Выход из системы |

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

## 🔧 Отладка

### Логи сервисов

```bash
# Логи auth-service
docker-compose -f docker-compose.integration.yml logs auth-service

# Логи PostgreSQL
docker-compose -f docker-compose.integration.yml logs postgres

# Логи Redis
docker-compose -f docker-compose.integration.yml logs redis

# Логи migrate
docker-compose -f docker-compose.integration.yml logs migrate

# Все логи в реальном времени
docker-compose -f docker-compose.integration.yml logs -f
```

### Подключение к контейнерам

```bash
# PostgreSQL
docker-compose -f docker-compose.integration.yml exec postgres psql -U test_user -d auth_test

# Redis
docker-compose -f docker-compose.integration.yml exec redis redis-cli

# Auth Service
docker-compose -f docker-compose.integration.yml exec auth-service sh
```

### Проверка здоровья

```bash
# Health check
curl http://localhost:8080/health

# Статус контейнеров
docker-compose -f docker-compose.integration.yml ps
```

## 🐛 Troubleshooting

### Ошибка: "Cannot connect to the Docker daemon"

```bash
# Проверить Docker
docker ps

# Если не работает, запустить Docker
sudo systemctl start docker
```

### Ошибка: "Service not ready"

```bash
# Проверить логи
docker-compose -f docker-compose.integration.yml logs auth-service

# Проверить, что миграции применились
docker-compose -f docker-compose.integration.yml logs migrate

# Перезапустить
docker-compose -f docker-compose.integration.yml down -v
docker-compose -f docker-compose.integration.yml up -d --build
```

### Ошибка: "Connection refused"

```bash
# Проверить, что сервис запущен
docker-compose -f docker-compose.integration.yml ps

# Проверить порты
docker-compose -f docker-compose.integration.yml port auth-service 8080
```

### Тесты не находят контейнеры

```bash
# Убедиться, что docker-compose запущен
docker-compose -f docker-compose.integration.yml ps

# Если сервисов нет, запустить
docker-compose -f docker-compose.integration.yml up -d --build
```

### Ошибка: "No space left on device"

```bash
# Очистить Docker
docker system prune -a

# Очистить volumes
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
