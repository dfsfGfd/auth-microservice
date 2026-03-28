# Configuration Guide

Настройка микросервиса через переменные окружения (.env).

---

## 🚀 Быстрый старт

```bash
# 1. Скопировать шаблон
cp .env.example .env

# 2. Настроить минимум переменных
# .env:
DATABASE_URL=postgres://user:pass@localhost:5432/auth?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key-minimum-32-characters-long

# 3. Запустить
go run cmd/server/main.go
```

---

## 📋 Все переменные

### Server

```bash
SERVER_HTTP_PORT=8080          # REST API порт
SERVER_GRPC_PORT=9090          # gRPC порт
SERVER_ENV=development         # development, staging, production
SERVER_READ_TIMEOUT=10         # Таймаут чтения (сек)
SERVER_WRITE_TIMEOUT=10        # Таймаут записи (сек)
SERVER_IDLE_TIMEOUT=60         # Таймаут простоя (сек)
```

### Database (PostgreSQL)

```bash
DATABASE_URL=postgres://user:pass@localhost:5432/auth?sslmode=disable
DATABASE_MAX_CONNECTIONS=25
DATABASE_CONNECTION_TIMEOUT=10
```

### Redis

```bash
REDIS_URL=redis://localhost:6379
REDIS_DB=0
REDIS_CONNECTION_TIMEOUT=5
```

### JWT

```bash
# Обязательно для production! openssl rand -base64 32
JWT_SECRET=your-secret-key-minimum-32-characters-long
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h           # 30 дней
JWT_ISSUER=auth-service
```

### Logging

```bash
LOG_LEVEL=debug                # debug, info, warn, error, fatal
LOG_FORMAT=console             # json, console
LOG_SERVICE_NAME=auth-service
```

### Rate Limiting

```bash
RATE_LIMIT_REGISTER=5          # запросов/мин
RATE_LIMIT_LOGIN=10
RATE_LIMIT_REFRESH=30
RATE_LIMIT_LOGOUT=60
```

### Health & Shutdown

```bash
HEALTH_PATH=/health
SHUTDOWN_TIMEOUT=30
```

---

## 🔧 Окружения

### Development

```bash
SERVER_ENV=development
LOG_LEVEL=debug
LOG_FORMAT=console
DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth?sslmode=disable
REDIS_URL=redis://localhost:6379
JWT_SECRET=super-secret-key-minimum-32-characters-long
```

### Production

```bash
SERVER_ENV=production
LOG_LEVEL=warn
LOG_FORMAT=json
DATABASE_URL=postgres://user:pass@db.prod:5432/auth?sslmode=require
DATABASE_MAX_CONNECTIONS=100
REDIS_URL=redis://redis.prod:6379
JWT_SECRET=<openssl rand -base64 32>
```

---

## 🔒 Безопасность

### Требования к паролю

Согласно **NIST 800-63B**:

- ✅ Минимум **8 символов**
- ✅ Без требований к заглавным/строчным буквам или цифрам
- ✅ Поддержка passphrase

**Примеры:**
```
password123    ✅
PASSWORD123    ✅
Password       ✅
abcdefgh       ✅
12345678       ❌ (слишком простой)
1234567        ❌ (короткий)
```

### Генерация JWT_SECRET

```bash
openssl rand -base64 32
```

### Production Checklist

- [ ] `JWT_SECRET` ≥ 32 символа
- [ ] `APP_ENV=production`
- [ ] `DATABASE_URL` с `sslmode=require`
- [ ] `.env` не в git

---

## 📁 Файлы

| Файл | Описание |
|------|----------|
| `.env` | Переменные (не коммитить) |
| `.env.example` | Шаблон (коммитить) |

---

## 📚 Ссылки

- [README](docs/README.md)
- [API](docs/api.md)
- [Docker](deploy/DEPLOY.md)
