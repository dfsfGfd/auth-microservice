# Auth Microservice

> Микросервис аутентификации с поддержкой gRPC и REST API, использующий JWT access/refresh токены.

---

## 📋 Оглавление

- [Возможности](#-возможности)
- [Архитектура](#-архитектура)
- [Структура проекта](#-структура-проекта)
- [Быстрый старт](#-быстрый-старт)
- [API Endpoints](#-api-endpoints)
- [Разработка](#-разработка)
- [Безопасность](#-безопасность)
- [Зависимости](#-зависимости)

---

## 📋 Возможности

| Функция | Описание |
|---------|----------|
| ✅ **Регистрация** | Создание нового аккаунта |
| ✅ **Вход/Выход** | Аутентификация и завершение сессии |
| ✅ **Обновление токенов** | Ротация JWT access/refresh токенов |
| ✅ **gRPC + REST** | Единый сервис для обоих протоколов (grpc-gateway) |
| ✅ **Swagger/OpenAPI** | Автогенерируемая документация API |
| ✅ **Rate Limiting** | Redis-based ограничение запросов (sliding window) |
| ✅ **CORS Middleware** | Настройка跨origin запросов |
| ✅ **Защита от User Enumeration** | Одинаковые ошибки аутентификации |

---

## 🏗 Архитектура

### Стек технологий

| Компонент | Технология |
|-----------|------------|
| **RPC** | gRPC + REST (grpc-gateway) |
| **Токены** | JWT (access + refresh) |
| **Кэш** | Redis |
| **БД** | PostgreSQL (pgx) |
| **Протокол** | Protocol Buffers v3 |
| **UUID** | Генерация на бэкенде (Go) |

### Время жизни токенов

| Токен | TTL | Хранение |
|-------|-----|----------|
| **Access Token** | 15 минут | Клиент (Authorization header) |
| **Refresh Token** | 2 недели | Redis (`refresh:{token}` → `account_id`) |

---

## 📁 Структура проекта

```
.
├── cmd/
│   └── server/                 # Точка входа приложения
│
├── internal/
│   ├── model/                  # Доменные модели (агрегаты, VO)
│   │   ├── account.go          # Account агрегат
│   │   ├── email.go            # Email VO
│   │   ├── password.go         # PlainPassword VO
│   │   └── password_hash.go    # PasswordHash VO
│   │
│   ├── cache/                  # Кэш слой
│   │   └── token/              # Кэш для токенов (Redis)
│   │
│   ├── middleware/             # HTTP/gRPC middleware
│   │   ├── rate_limiter.go     # Redis-based rate limiting
│   │   ├── rate_limiter_http.go # HTTP/gRPC адаптеры
│   │   └── cors.go             # CORS middleware
│   │
│   ├── repository/             # Репозитории (PostgreSQL)
│   │   ├── repository.go       # Интерфейсы репозиториев
│   │   ├── model/              # DB модели
│   │   ├── converter/          # Конвертеры domain ↔ DB
│   │   └── auth/               # PostgreSQL реализация
│   │
│   ├── config/                 # Загрузка конфигурации
│   ├── di/                     # Dependency Injection (Google Wire)
│   ├── errors/                 # Доменные ошибки
│   ├── service/                # Бизнес-логика (сервисный слой)
│   └── handler/                # gRPC хендлеры
│
├── pkg/
│   ├── proto/                  # Сгенерированный Proto код (gRPC + REST)
│   ├── bcrypt/                 # Хеширование паролей
│   ├── jwt/                    # JWT утилиты
│   ├── cookies/                # Cookie утилиты
│   ├── logger/                 # Логирование (zerolog)
│   └── db/                     # Подключения к БД
│       ├── postgresql/         # PostgreSQL подключение
│       └── redisdb/            # Redis подключение
│
├── api/
│   └── auth/v1/                # Swagger/OpenAPI спецификации
├── proto/
│   └── auth/v1/
│       └── auth.proto          # Proto контракты
│
├── migrations/
│   └── 001_create_accounts_table.sql  # Миграция БД
│
├── docs/
│   ├── README.md               # Основная документация
│   ├── api.md                  # API документация
│   ├── config.md               # Руководство по конфигурации
│   └── repository_methods.md   # Repository layer документация
│
├── config.yaml                 # Локальная конфигурация
├── config.example.yaml         # Шаблон конфигурации
├── Taskfile.yml                # Taskfile команды
└── go.mod
```

---

## 🚀 Быстрый старт

### Требования

| Зависимость | Версия |
|-------------|--------|
| Go | 1.25+ |
| PostgreSQL | 15+ |
| Redis | 7+ |
| Task | latest |

### Установка зависимостей

```bash
# Установка инструментов разработки
task proto:install-plugins
task install-buf
task install-formatters
task install-golangci-lint
```

### Генерация Proto

```bash
task proto:gen
```

### Применение миграций

```bash
# Пример с goose (если используется)
goose -dir migrations postgres "DATABASE_URL" up
```

### Запуск

```bash
# Переменные окружения
export DATABASE_URL="postgres://user:pass@localhost:5432/auth?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="your-secret-key-minimum-32-characters-long"  # Обязательно!
export APP_ENV="development"  # production для продакшена

# Запуск сервера
go run cmd/server/main.go
```

### Production запуск

```bash
# Обязательно установите JWT_SECRET
export JWT_SECRET=$(openssl rand -base64 32)
export APP_ENV=production

# Запуск
./server
```

---

## 📡 API Endpoints

### Основные методы

| Метод | gRPC | HTTP | Описание |
|-------|------|------|----------|
| `Register` | `Register` | `POST /api/auth/register` | Регистрация |
| `Login` | `Login` | `POST /api/auth/login` | Вход |
| `Logout` | `Logout` | `POST /api/auth/logout` | Выход |
| `Refresh` | `Refresh` | `POST /api/auth/refresh` | Обновление токенов |

### Формат ответов

Все ответы API имеют единую структуру:

```json
{
  "status_code": 200,
  "message": "Success",
  "data": { ... }
}
```

| Поле | Тип | Описание |
|------|-----|----------|
| `status_code` | `int` | HTTP статус код |
| `message` | `string` | Сообщение статуса |
| `data` | `object` | Тело ответа (может быть `null`) |

📖 **Подробное описание API:** [docs/api.md](api.md)

---

## 🛠 Разработка

### Форматирование

```bash
task format
```

### Линтинг

```bash
# Go код
task lint

# Proto файлы
task proto:lint
```

### Тесты

```bash
go test ./... -v
```

### Taskfile команды

| Команда | Описание |
|---------|----------|
| `task proto:gen` | Генерация Proto (gRPC + REST + Swagger) |
| `task proto:lint` | Линтинг Proto файлов |
| `task proto:deps` | Обновление Proto зависимостей |
| `task format` | Форматирование Go кода |
| `task lint` | Линтинг Go кода |
| `task tidy` | Очистка зависимостей |
| `task wire:gen` | Генерация DI кода |

---

## 🔒 Безопасность

### JWT Claims

Токены содержат название сервиса в поле `iss` (issuer):

```json
{
  "iss": "auth-service",
  "sub": "{account_id}",
  "email": "{email}",
  "iat": 1705312200,
  "exp": 1705313100,
  "type": "access"
}
```

### Rate Limiting

| Endpoint | Лимит (запросов/мин) |
|----------|---------------------|
| `/api/auth/register` | 5 |
| `/api/auth/login` | 10 |
| `/api/auth/refresh` | 30 |
| `/api/auth/logout` | 60 |

**Заголовки ответа:**
```
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 5
X-RateLimit-Reset: 1647389400
Retry-After: 60
```

### Защита от User Enumeration

Все ошибки аутентификации возвращают одинаковый ответ:
```json
{
  "status_code": 401,
  "message": "invalid credentials",
  "data": null
}
```

Это предотвращает определение существования email в системе.

### Требования к паролю

| Требование | Значение |
|------------|----------|
| Минимальная длина | 8 символов |
| Заглавные буквы | Минимум 1 (A-Z) |
| Строчные буквы | Минимум 1 (a-z) |
| Цифры | Минимум 1 (0-9) |

### Требования к email

| Требование | Значение |
|------------|----------|
| Формат | RFC 5321 |
| Максимальная длина | 254 символа |

### Production Checklist

- [ ] `JWT_SECRET` установлен через environment variable
- [ ] `APP_ENV=production` для автоматического включения Secure cookies
- [ ] Настроен CORS для ваших доменов
- [ ] Rate limits настроены под вашу нагрузку
- [ ] HTTPS включён (reverse proxy: nginx, traefik)

---

## 📦 Зависимости

### Внешние библиотеки

```go
// Стандартные библиотеки
github.com/google/uuid                  // UUID генерация
github.com/google/wire                  // Dependency Injection
github.com/golang-jwt/jwt/v5            // JWT токены
github.com/rs/zerolog                   // Логирование
github.com/rs/cors                      // CORS middleware
github.com/redis/go-redis/v9            // Redis клиент
github.com/jackc/pgx/v5                 // PostgreSQL драйвер
golang.org/x/crypto                     // bcrypt
google.golang.org/grpc                  // gRPC
google.golang.org/protobuf              // Protocol Buffers
github.com/grpc-ecosystem/grpc-gateway/v2 // gRPC → REST
```

### Proto зависимости

```yaml
deps:
  - buf.build/googleapis/googleapis
  - buf.build/grpc-ecosystem/grpc-gateway
```

---

## 📚 Документация

| Документ | Описание |
|----------|----------|
| [API Documentation](api.md) | Полное описание API endpoints |
| [Configuration Guide](config.md) | Руководство по настройке |
| [Repository Methods](repository_methods.md) | Repository layer с DDD паттернами |
| [Swagger UI](http://localhost:8080/swagger/) | Интерактивная документация (после запуска) |

---

## 📝 License

MIT
