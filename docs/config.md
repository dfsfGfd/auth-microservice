# Configuration Guide

Настройка микросервиса аутентификации.

## 📋 Быстрый старт

### 1. Создание .env файла

```bash
# Для локальной разработки
cp .env.example .env
```

### 2. Минимальная конфигурация

Для запуска достаточно настроить 3 переменные в `.env`:

```bash
# .env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key-minimum-32-characters-long
```

### 3. Запуск

```bash
go run cmd/server/main.go
```

> **Примечание:** Конфигурация загружается из `.env` файла. YAML конфигурация (`config.yaml`) поддерживается опционально.

## 🔧 Конфигурация

### Server

Микросервис поддерживает **параллельную работу REST и gRPC**:

```bash
# .env
APP_ENV=development          # development, staging, production
HTTP_PORT=8080               # Порт HTTP (REST + grpc-gateway)
GRPC_PORT=9090               # Порт gRPC (прямой доступ)
READ_TIMEOUT=10              # Таймаут чтения (сек)
WRITE_TIMEOUT=10             # Таймаут записи (сек)
IDLE_TIMEOUT=60              # Таймаут простоя соединения (сек)
```

**Архитектура:**

```
┌─────────────────────────────────────────────────┐
│                  Auth Microservice              │
│                                                 │
│  ┌─────────────┐    ┌─────────────────────┐    │
│  │ HTTP:8080   │───▶│  grpc-gateway       │    │
│  │ (REST API)  │    │  (REST → gRPC)      │    │
│  └─────────────┘    └──────────┬──────────┘    │
│                                │                │
│  ┌─────────────┐               │                │
│  │ gRPC:9090   │───────────────┼────────────────┤
│  │ (direct)    │               │                │
│  └─────────────┘               ▼                │
│                         ┌─────────────┐         │
│                         │  gRPC       │         │
│                         │  Service    │         │
│                         └─────────────┘         │
└─────────────────────────────────────────────────┘
```

| Порт | Протокол | Описание |
|------|----------|----------|
| `8080` | HTTP/REST | REST API через grpc-gateway |
| `9090` | gRPC | Прямой gRPC доступ для микросервисов |

### Database (PostgreSQL)

```bash
# .env
DATABASE_URL=postgres://user:pass@localhost:5432/auth?sslmode=disable
DATABASE_MAX_CONNECTIONS=25       # Максимум подключений в пуле
DATABASE_CONNECTION_TIMEOUT=10    # Таймаут подключения (сек)
```

### Redis

```bash
# .env
REDIS_URL=redis://localhost:6379
REDIS_DB=0                        # DB номер (0-15)
REDIS_CONNECTION_TIMEOUT=5        # Таймаут подключения (сек)
```

### JWT

```bash
# .env
# Обязательно для production! Генерируйте через: openssl rand -base64 32
JWT_SECRET=your-secret-key-minimum-32-characters-long
JWT_ACCESS_TTL=15m                # Время жизни access токена
JWT_REFRESH_TTL=336h              # Время жизни refresh токена (14 дней)
JWT_ISSUER=auth-service           # Название сервиса (iss claim)
```

> **Важно:** Для production используйте `JWT_SECRET` через environment variable, а не храните в файле.

### Cookie

```bash
# .env
COOKIE_SECURE=false               # true для HTTPS (auto в production)
COOKIE_HTTP_ONLY=true             # Защита от XSS
COOKIE_SAME_SITE=Lax              # Strict, Lax, None
COOKIE_DOMAIN=                    # Домен (опционально)
COOKIE_PATH=/
COOKIE_MAX_AGE=1209600            # 14 дней в секундах
```

### Logging

```bash
# .env
LOG_LEVEL=debug                   # debug, info, warn, error, fatal
LOG_FORMAT=console                # json, console
LOG_SERVICE_NAME=auth-service
```

### CORS

```bash
# .env
CORS_ALLOWED_ORIGINS=https://your-domain.com,https://app.your-domain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Authorization,Content-Type,X-Request-ID
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=86400
```

> **Примечание:** CORS_ALLOWED_ORIGINS поддерживает несколько origin через запятую.

### Rate Limiting

```bash
# .env
RATE_LIMIT_REGISTER=5    # Лимит запросов в минуту на endpoint
RATE_LIMIT_LOGIN=10
RATE_LIMIT_REFRESH=30
RATE_LIMIT_LOGOUT=60
```

> **Алгоритм:** Sliding window с Redis хранением.
> **При превышении:** HTTP 429 Too Many Requests с заголовком `Retry-After`.

### Health Check

```bash
# .env
HEALTH_PATH=/health      # Endpoint для проверок здоровья
```

### Graceful Shutdown

```bash
# .env
SHUTDOWN_TIMEOUT=30      # Таймаут завершения работы (сек)
```

## 🖥 Конфигурация по окружениям

### Development (локально)

```bash
# .env
APP_ENV=development
HTTP_PORT=8080

LOG_LEVEL=debug
LOG_FORMAT=console

COOKIE_SECURE=false

DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth?sslmode=disable
REDIS_URL=redis://localhost:6379
```

### Production

```bash
# .env
APP_ENV=production
HTTP_PORT=8080

LOG_LEVEL=warn
LOG_FORMAT=json

# Автоматически включается при APP_ENV=production
COOKIE_SECURE=true
COOKIE_SAME_SITE=Strict

DATABASE_URL=postgres://user:pass@db.prod:5432/auth?sslmode=require
DATABASE_MAX_CONNECTIONS=100

REDIS_URL=redis://redis.prod:6379

# Обязательно! Генерируйте через: openssl rand -base64 32
JWT_SECRET=<crypto-random-32-chars>

CORS_ALLOWED_ORIGINS=https://your-domain.com
CORS_ALLOW_CREDENTIALS=true

RATE_LIMIT_REGISTER=5
RATE_LIMIT_LOGIN=10
RATE_LIMIT_REFRESH=30
RATE_LIMIT_LOGOUT=60
```

## 🔒 Безопасность

### Генерация JWT_SECRET

```bash
# OpenSSL
openssl rand -base64 32

# Или с помощью Go
go run -e 'package main; import "crypto/rand"; import "encoding/base64"; func main() { b := make([]byte, 32); rand.Read(b); println(base64.StdEncoding.EncodeToString(b)) }'
```

### Чеклист для продакшена

- [ ] `JWT_SECRET` установлен через `.env` или environment variable
- [ ] `APP_ENV=production` для автоматического включения Secure cookies
- [ ] `COOKIE_SECURE=true` — только HTTPS
- [ ] `COOKIE_HTTP_ONLY=true` — защита от XSS
- [ ] `COOKIE_SAME_SITE=Strict` — защита от CSRF
- [ ] `LOG_LEVEL=warn` или `error` — не логировать лишнего
- [ ] `DATABASE_URL` — с `sslmode=require`
- [ ] `.env` не закоммичен в git
- [ ] CORS настроен только для ваших доменов
- [ ] Rate limits настроены под вашу нагрузку
- [ ] HTTPS включён (reverse proxy: nginx, traefik)

## 📁 Файлы конфигурации

| Файл | Описание |
|------|----------|
| `.env` | Переменные окружения (не коммитить, создать из .env.example) |
| `.env.example` | Шаблон переменных окружения (коммитить) |
| `config.example.yaml` | YAML шаблон (опционально, для совместимости) |

## 🧪 Тестирование

Для тестов используется отдельная конфигурация:

```bash
# .env.test
DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth_test?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=test-secret-key-for-testing-only
LOG_LEVEL=error
```

## 📚 Ссылки

- [Шаблон конфигурации](../config.example.yaml)
- [README](../README.md)
- [API Documentation](api.md)
