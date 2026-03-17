# Security Fixes — P0 Critical Vulnerabilities

Этот документ описывает исправленные критические уязвимости безопасности для высоконагруженного сервиса аутентификации.

---

## ✅ Исправленные уязвимости

### 1. Rate Limiting Middleware

**Проблема:** Отсутствовало ограничение количества запросов, что позволяло:
- Брутфорсить пароли
- Проводить DoS-атаки
- Массово регистрировать аккаунты

**Решение:**
- Создан пакет `internal/middleware/rate_limiter.go`
- Используется Redis-based sliding window алгоритм
- Настройка лимитов через `.env` (переменные `RATE_LIMIT_*`)
- Отдельные лимиты для каждого endpoint'а:
  - `register`: 5 запросов/мин
  - `login`: 10 запросов/мин
  - `refresh`: 30 запросов/мин
  - `logout`: 60 запросов/мин

**Файлы:**
- `internal/middleware/rate_limiter.go` — ядро rate limiting
- `internal/middleware/rate_limiter_http.go` — HTTP/gRPC middleware
- `internal/di/provider.go` — DI интеграция
- `cmd/server/main.go` — применение middleware

**Использование:**
```bash
# Настройка лимитов в .env
RATE_LIMIT_REGISTER=5
RATE_LIMIT_LOGIN=10
RATE_LIMIT_REFRESH=30
RATE_LIMIT_LOGOUT=60
```

**Заголовки ответа:**
```
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 5
X-RateLimit-Reset: 1647389400
Retry-After: 60
```

---

### 2. Fix User Enumeration

**Проблема:** Разные ошибки для "аккаунт не найден" и "неверный пароль" позволяли злоумышленнику определять существование email в системе.

**Решение:**
- Все ошибки аутентификации возвращают одинаковый ответ `"invalid credentials"`
- Детальное логирование сохраняется для внутреннего использования
- Изменён `internal/handler/auth/login.go`
- Изменён `internal/service/auth/login.go`

**Файлы:**
- `internal/handler/auth/login.go` — одинаковые ошибки в handler
- `internal/service/auth/login.go` — одинаковые ошибки в service

**До:**
```go
case stderrors.Is(err, errors.ErrAccountNotFound):
    return nil, status.Error(codes.NotFound, "account not found")
case stderrors.Is(err, errors.ErrInvalidCredentials):
    return nil, status.Error(codes.Unauthenticated, "invalid credentials")
```

**После:**
```go
// Всегда возвращаем одинаковую ошибку
return nil, status.Error(codes.Unauthenticated, "invalid credentials")
```

---

### 3. JWT Secret из Environment Variables

**Проблема:** JWT secret хранился в файле конфигурации, что создавало риск утечки через git.

**Решение:**
- Добавлена поддержка переменной окружения `JWT_SECRET`
- Загрузка из `.env` файла (приоритет)
- Валидация в `internal/config/config.go`
- Обновлён `.env.example`

**Файлы:**
- `internal/config/config.go` — чтение из `JWT_SECRET`
- `.env.example` — шаблон переменных

**Использование:**
```bash
# Production
export JWT_SECRET="your-super-secret-key-minimum-32-characters-long"

# Или через .env
cp .env.example .env
# Отредактируйте .env, установив JWT_SECRET
```

---

### 4. HTTPS Enforcement в Production

**Проблема:** Cookie флаг `Secure` был отключён по умолчанию, что позволяло передавать токены по HTTP.

**Решение:**
- Автоматическое включение `Secure: true` для `APP_ENV=production`
- Валидация в `internal/config/config.go`
- Обновлён `config.example.yaml`

**Файлы:**
- `internal/config/config.go` — автоматическое включение Secure
- `config.example.yaml` — документация

**Использование:**
```bash
# Production
export APP_ENV=production

# Cookie автоматически будет установлен с Secure=true
```

---

### 5. CORS Middleware

**Проблема:** CORS был настроен в конфиге, но не применялся, что позволяло любым сайтам делать запросы к API.

**Решение:**
- Создан `internal/middleware/cors.go` с использованием `github.com/rs/cors`
- Настройка из `.env` (переменные `CORS_*`)
- Применение в `cmd/server/main.go`

**Файлы:**
- `internal/middleware/cors.go` — CORS middleware
- `cmd/server/main.go` — применение middleware
- `.env.example` — конфигурация

**Использование:**
```bash
# .env
CORS_ALLOWED_ORIGINS=https://your-domain.com,https://app.your-domain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Authorization,Content-Type,X-Request-ID
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=86400
```

---

## 📊 Тестирование

### Rate Limiting
```bash
# Быстрые запросы к login endpoint
for i in {1..15}; do
  curl -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@test.com","password":"password"}' \
    -w "\nHTTP Status: %{http_code}\n"
done

# После 10 запросов получите 429 Too Many Requests
```

### User Enumeration
```bash
# Несуществующий email
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"notexist@test.com","password":"password"}'

# Ответ: {"error":"invalid credentials"} (401)

# Существующий email, неверный пароль
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"exist@test.com","password":"wrongpassword"}'

# Ответ: {"error":"invalid credentials"} (401)
# Одинаковые ответы!
```

### JWT Secret
```bash
# Без JWT_SECRET
unset JWT_SECRET
./server  # Ошибка: secret is required

# С JWT_SECRET
export JWT_SECRET="super-secret-key-min-32-chars-long"
./server  # Успешный запуск
```

### CORS
```bash
# Запрос с другого origin
curl -X OPTIONS http://localhost:8080/api/auth/login \
  -H "Origin: https://evil.com" \
  -H "Access-Control-Request-Method: POST" \
  -i

# Проверьте заголовки:
# Access-Control-Allow-Origin: (только разрешённые origins)
```

---

## 🚀 Развёртывание

### Docker Compose (production)
```yaml
version: '3.8'
services:
  auth-service:
    image: auth-microservice:latest
    env_file:
      - .env.production
    ports:
      - "8080:8080"
```

### Environment Variables
```bash
# Обязательные для production
export JWT_SECRET="your-super-secret-key-minimum-32-characters-long"
export APP_ENV=production

# Опциональные
export DATABASE_URL="postgres://..."
export REDIS_URL="redis://..."
```

### Через .env файл
```bash
# 1. Скопируйте шаблон
cp .env.example .env

# 2. Отредактируйте .env
# Установите: JWT_SECRET, DATABASE_URL, REDIS_URL

# 3. Запуск
go run cmd/server/main.go
```

---

## 📝 Checklist для Production

- [ ] `.env` создан и настроен (не закоммичен в git!)
- [ ] `JWT_SECRET` установлен через environment variable или .env
- [ ] `APP_ENV=production`
- [ ] Настроен CORS `CORS_ALLOWED_ORIGINS` для вашего домена
- [ ] Rate limits настроены под вашу нагрузку
- [ ] Включить HTTPS (reverse proxy: nginx, traefik)
- [ ] Проверить, что cookie передаются только по HTTPS
- [ ] Протестировать rate limiting
- [ ] Настроить мониторинг (метрики, алерты)

---

## 🔗 Связанные документы

- [API Documentation](docs/api.md)
- [Configuration Guide](docs/config.md)
- [Security Best Practices](docs/security.md)
