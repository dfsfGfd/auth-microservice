# Configuration Guide

Настройка микросервиса аутентификации.

## 📋 Быстрый старт

### 1. Создание config.yaml файла

```bash
# Для локальной разработки
cp config.example.yaml config.yaml
```

### 2. Минимальная конфигурация

Для запуска достаточно настроить 3 секции:

```yaml
# config.yaml
database:
  url: postgres://postgres:postgres@localhost:5432/auth?sslmode=disable

redis:
  url: redis://localhost:6379

jwt:
  secret: your-secret-key-minimum-32-characters-long
```

### 3. Запуск

```bash
go run cmd/server/main.go
```

## 🔧 Конфигурация

### Server

```yaml
server:
  port: 8080           # Порт HTTP сервера
  grpc_port: 9090      # Порт gRPC сервера
  env: development     # development, staging, production
```

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
  secret: super-secret-key-minimum-32-characters
  access_ttl: 15m              # Время жизни access токена
  refresh_ttl: 336h            # Время жизни refresh токена (14 дней)
  issuer: auth-service         # Название сервиса (iss claim)
```

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
    - http://localhost:3000
    - http://localhost:8080
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
  max_age: 86400
```

### Rate Limiting

```yaml
rate_limit:
  register: 5    # Лимит запросов в минуту
  login: 10
  refresh: 30
  logout: 60
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
  secure: true
  same_site: Strict

database:
  url: postgres://user:pass@db.prod:5432/auth?sslmode=require
  max_connections: 100

redis:
  url: redis://redis.prod:6379

jwt:
  secret: <crypto-random-32-chars>
```

## 🔒 Безопасность

### Генерация JWT_SECRET

```bash
# OpenSSL
openssl rand -base64 32

# Go
go run -e 'package main; import "crypto/rand"; import "encoding/base64"; func main() { b := make([]byte, 32); rand.Read(b); println(base64.StdEncoding.EncodeToString(b)) }'
```

### Чеклист для продакшена

- [ ] `jwt.secret` — криптографически случайный, минимум 32 символа
- [ ] `cookie.secure: true` — только HTTPS
- [ ] `cookie.http_only: true` — защита от XSS
- [ ] `cookie.same_site: Strict` — защита от CSRF
- [ ] `logging.level: warn` или `error` — не логировать лишнего
- [ ] `database.url` — с `sslmode=require`
- [ ] `config.yaml` не закоммичен в git

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
