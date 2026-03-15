# Configuration Guide

Настройка микросервиса аутентификации.

## 📋 Быстрый старт

### 1. Создание .env файла

```bash
# Для локальной разработки
cp .env.example .env
```

### 2. Минимальная конфигурация

Для запуска достаточно настроить 3 переменные:

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

## 🔧 Переменные окружения

### Обязательные

| Переменная | Описание | Пример |
|------------|----------|--------|
| `DATABASE_URL` | DSN PostgreSQL | `postgres://user:pass@localhost:5432/auth?sslmode=disable` |
| `REDIS_URL` | DSN Redis | `redis://localhost:6379` |
| `JWT_SECRET` | Секретный ключ JWT (мин. 32 символа) | `super-secret-key-change-in-production` |

### Рекомендуемые

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `SERVER_PORT` | `8080` | Порт HTTP сервера |
| `ENV` | `development` | Режим: `development`, `staging`, `production` |
| `LOG_LEVEL` | `info` | Уровень логов: `debug`, `info`, `warn`, `error` |
| `LOG_FORMAT` | `json` | Формат: `json`, `console` |

### JWT

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `JWT_ACCESS_TTL` | `15m` | Время жизни access токена |
| `JWT_REFRESH_TTL` | `336h` (14 дней) | Время жизни refresh токена |
| `JWT_ISSUER` | `auth-service` | Название сервиса (iss claim) |

### Cookie

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `COOKIE_SECURE` | `true` | Передача только по HTTPS |
| `COOKIE_HTTP_ONLY` | `true` | Защита от XSS |
| `COOKIE_SAME_SITE` | `Strict` | CSRF защита: `Strict`, `Lax`, `None` |
| `COOKIE_PATH` | `/` | Путь cookie |
| `COOKIE_MAX_AGE` | `1209600` (14 дней) | Время жизни в секундах |

## 🖥 Конфигурация по окружениям

### Development (локально)

```bash
ENV=development
LOG_LEVEL=debug
LOG_FORMAT=console
COOKIE_SECURE=false
DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth?sslmode=disable
REDIS_URL=redis://localhost:6379
```

### Production

```bash
ENV=production
LOG_LEVEL=warn
LOG_FORMAT=json
COOKIE_SECURE=true
COOKIE_SAME_SITE=Strict
DATABASE_URL=postgres://user:pass@db.prod:5432/auth?sslmode=require
REDIS_URL=redis://redis.prod:6379
JWT_SECRET=<crypto-random-32-chars>
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

- [ ] `JWT_SECRET` — криптографически случайный, минимум 32 символа
- [ ] `COOKIE_SECURE=true` — только HTTPS
- [ ] `COOKIE_HTTP_ONLY=true` — защита от XSS
- [ ] `COOKIE_SAME_SITE=Strict` — защита от CSRF
- [ ] `LOG_LEVEL=warn` или `error` — не логировать лишнего
- [ ] `DATABASE_URL` — с `sslmode=require`
- [ ] `.env` файл не закоммичен в git

## 📁 Файлы конфигурации

| Файл | Описание |
|------|----------|
| `.env` | Локальная конфигурация (не коммитить) |
| `.env.example` | Шаблон конфигурации (коммитить) |

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

- [Шаблон конфигурации](../.env.example)
- [README](../README.md)
- [API Documentation](api.md)
