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

> **Примечание:** Конфигурация загружается из `.env` файла (приоритет) или `config.yaml`.

## 🔧 Конфигурация

### Server

Микросервис поддерживает **параллельную работу REST и gRPC**:

```yaml
server:
  http_port: 8080          # Порт HTTP (REST + grpc-gateway)
  grpc_port: 9090          # Порт gRPC (прямой доступ)
  env: development         # development, staging, production
  read_timeout: 10         # Таймаут чтения (сек)
  write_timeout: 10        # Таймаут записи (сек)
  idle_timeout: 60         # Таймаут простоя соединения (сек)
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

```yaml
database:
  url: postgres://user:pass@localhost:5432/auth?sslmode=disable
  max_connections: 25          # Максимум подключений в пуле
  connection_timeout: 10       # Таймаут подключения (сек)
```

### Redis

```yaml
redis:
  url: redis://localhost:6379
  db: 0                        # DB номер (0-15)
  connection_timeout: 5        # Таймаут подключения (сек)
```

### JWT

```yaml
jwt:
  # Секретный ключ (минимум 32 символа)
  # ПРИОРИТЕТ: Если не указан в config.yaml, берётся из переменной окружения JWT_SECRET
  secret: ""  # Рекомендуется: export JWT_SECRET="your-secret-key..."
  access_ttl: 15m              # Время жизни access токена
  refresh_ttl: 336h            # Время жизни refresh токена (14 дней)
  issuer: auth-service         # Название сервиса (iss claim)
```

> **Важно:** Для production используйте переменную окружения `JWT_SECRET` вместо хранения в файле конфигурации.
>
> ```bash
> export JWT_SECRET=$(openssl rand -base64 32)
> ```

### Cookie

```yaml
cookie:
  secure: false                # true для HTTPS
  http_only: true              # Защита от XSS
  same_site: Lax               # Strict, Lax, None
  domain: ""                   # Домен (опционально)
  path: /
  max_age: 1209600             # 14 дней в секундах
```

### Logging

```yaml
logging:
  level: debug                 # debug, info, warn, error, fatal
  format: console              # json, console
  service_name: auth-service
```

### CORS

```yaml
cors:
  allowed_origins:
    - https://your-domain.com
    - https://app.your-domain.com
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
  allowed_headers:
    - Authorization
    - Content-Type
    - X-Request-ID
    - X-RateLimit-Limit
    - X-RateLimit-Remaining
    - X-RateLimit-Reset
  allow_credentials: true      # Разрешить cookies/credentials
  max_age: 86400               # Pre-flight cache (сек)
```

> **Примечание:** Поддерживаются wildcard поддомены (например, `*.example.com`).

### Rate Limiting

```yaml
rate_limit:
  register: 5    # Лимит запросов в минуту на endpoint
  login: 10
  refresh: 30
  logout: 60
```

> **Алгоритм:** Sliding window с Redis хранением.
>
> **При превышении:** HTTP 429 Too Many Requests с заголовком `Retry-After`.

### Health Check

```yaml
health:
  path: /health  # Endpoint для проверок здоровья
```

### Graceful Shutdown

```yaml
shutdown:
  timeout: 30    # Таймаут завершения работы (сек)
```

## 🖥 Конфигурация по окружениям

### Development (локально)

```yaml
server:
  env: development
  port: 8080

logging:
  level: debug
  format: console

cookie:
  secure: false

database:
  url: postgres://postgres:postgres@localhost:5432/auth?sslmode=disable

redis:
  url: redis://localhost:6379
```

### Production

```yaml
server:
  env: production
  port: 8080

logging:
  level: warn
  format: json

cookie:
  secure: true         # Автоматически включается при APP_ENV=production
  same_site: Strict

database:
  url: postgres://user:pass@db.prod:5432/auth?sslmode=require
  max_connections: 100

redis:
  url: redis://redis.prod:6379

jwt:
  secret: ""  # Используйте JWT_SECRET env var!

cors:
  allowed_origins:
    - https://your-domain.com
  allow_credentials: true

rate_limit:
  register: 5
  login: 10
  refresh: 30
  logout: 60
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

- [ ] `JWT_SECRET` установлен через environment variable (не в config.yaml!)
- [ ] `APP_ENV=production` для автоматического включения Secure cookies
- [ ] `cookie.secure: true` — только HTTPS
- [ ] `cookie.http_only: true` — защита от XSS
- [ ] `cookie.same_site: Strict` — защита от CSRF
- [ ] `logging.level: warn` или `error` — не логировать лишнего
- [ ] `database.url` — с `sslmode=require`
- [ ] `config.yaml` не закоммичен в git
- [ ] CORS настроен только для ваших доменов
- [ ] Rate limits настроены под вашу нагрузку
- [ ] HTTPS включён (reverse proxy: nginx, traefik)

## 📁 Файлы конфигурации

| Файл | Описание |
|------|----------|
| `config.yaml` | Локальная конфигурация (не коммитить) |
| `config.example.yaml` | Шаблон конфигурации (коммитить) |

## 🧪 Тестирование

Для тестов используется отдельная конфигурация:

```yaml
# config.test.yaml
database:
  url: postgres://postgres:postgres@localhost:5432/auth_test?sslmode=disable

redis:
  url: redis://localhost:6379

jwt:
  secret: test-secret-key-for-testing-only

logging:
  level: error
```

## 📚 Ссылки

- [Шаблон конфигурации](../config.example.yaml)
- [README](../README.md)
- [API Documentation](api.md)
