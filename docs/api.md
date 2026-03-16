# Auth Microservice — API Documentation

> Микросервис аутентификации с поддержкой JWT access/refresh токенов, gRPC и REST API.

---

## 📋 Оглавление

- [Обзор](#обзор)
- [Архитектура](#архитектура)
- [Формат ответов API](#формат-ответов-api)
- [API Endpoints](#api-endpoints)
- [Требования к данным](#требования-к-данным)
- [Хранение данных](#хранение-данных)
- [Безопасность](#безопасность)
- [Разработка](#разработка)

---

## Обзор

Микросервис предоставляет следующие возможности:

| Функция | Описание |
|---------|----------|
| 🔐 **Регистрация** | Создание нового аккаунта |
| 🔑 **Вход/Выход** | Аутентификация и завершение сессии |
| 🔄 **Обновление токенов** | Ротация JWT access/refresh токенов |
| 🌐 **gRPC + REST** | Единый сервис для обоих протоколов |

---

## Архитектура

### Стек технологий

| Компонент | Технология |
|-----------|------------|
| **RPC** | gRPC + REST (grpc-gateway) |
| **Токены** | JWT (access + refresh) |
| **Кэш** | Redis (refresh токены) |
| **БД** | PostgreSQL (pgx) |
| **Протокол** | Protocol Buffers v3 |
| **UUID** | Генерация на бэкенде (Go) |

### Время жизни токенов

| Токен | TTL | Хранение |
|-------|-----|----------|
| **Access Token** | 15 минут | Клиент (Authorization header) |
| **Refresh Token** | 2 недели | Redis (`refresh:{token}` → `account_id`) |

---

## Формат ответов API

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

### Пример успешного ответа

```json
{
  "status_code": 200,
  "message": "Account registered successfully",
  "data": {
    "account_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

### Пример ответа с ошибкой

```json
{
  "status_code": 400,
  "message": "Invalid email format",
  "data": null
}
```

---

## API Endpoints

### gRPC Service: `AuthService`

| Метод | HTTP | Описание |
|-------|------|----------|
| `Register` | `POST /api/auth/register` | Регистрация аккаунта |
| `Login` | `POST /api/auth/login` | Вход (получение токенов) |
| `Logout` | `POST /api/auth/logout` | Выход (отзыв refresh токена) |
| `Refresh` | `POST /api/auth/refresh` | Обновление access токена |

---

## Детали эндпоинтов

### 1. Регистрация аккаунта

```http
POST /api/auth/register
```

**gRPC:** `Register(RegisterRequest) returns (RegisterResponse)`

#### Request

```json
{
  "email": "user@example.com",
  "password": "Password123"
}
```

#### Response (200 OK)

```json
{
  "status_code": 200,
  "message": "Account registered successfully",
  "data": {
    "account_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

#### Ошибки

| Код | HTTP | Описание |
|-----|------|----------|
| `INVALID_EMAIL` | 400 | Неверный формат email |
| `EMAIL_EXISTS` | 409 | Email уже зарегистрирован |
| `INVALID_PASSWORD` | 400 | Пароль не соответствует требованиям |

#### Пример ошибки (400 Bad Request)

```json
{
  "status_code": 400,
  "message": "Invalid email format",
  "data": null
}
```

---

### 2. Вход

```http
POST /api/auth/login
```

**gRPC:** `Login(LoginRequest) returns (LoginResponse)`

#### Request

```json
{
  "email": "user@example.com",
  "password": "Password123"
}
```

#### Response (200 OK)

```json
{
  "status_code": 200,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "dGhpc2lzYXJlZnJlc2h0b2tlbg...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

#### Ошибки

| Код | HTTP | Описание |
|-----|------|----------|
| `INVALID_CREDENTIALS` | 401 | Неверный email или пароль |
| `ACCOUNT_NOT_FOUND` | 404 | Аккаунт не найден |

#### Пример ошибки (401 Unauthorized)

```json
{
  "status_code": 401,
  "message": "Invalid credentials",
  "data": null
}
```

---

### 3. Выход

```http
POST /api/auth/logout
Authorization: Bearer <access_token>
```

**gRPC:** `Logout(LogoutRequest) returns (LogoutResponse)`

#### Request

```json
{
  "refresh_token": "dGhpc2lzYXJlZnJlc2h0b2tlbg..."
}
```

#### Response (200 OK)

```json
{
  "status_code": 200,
  "message": "Logout successful",
  "data": {
    "success": true
  }
}
```

#### Ошибки

| Код | HTTP | Описание |
|-----|------|----------|
| `UNAUTHORIZED` | 401 | Access токен недействителен |
| `TOKEN_EXPIRED` | 401 | Access токен истёк |

---

### 4. Обновление токенов

```http
POST /api/auth/refresh
```

**gRPC:** `Refresh(RefreshRequest) returns (RefreshResponse)`

#### Request

```json
{
  "refresh_token": "dGhpc2lzYXJlZnJlc2h0b2tlbg..."
}
```

#### Response (200 OK)

```json
{
  "status_code": 200,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "bmV3cmVmcmVzaHRva2Vu...",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

#### Ошибки

| Код | HTTP | Описание |
|-----|------|----------|
| `INVALID_REFRESH_TOKEN` | 401 | Refresh токен недействителен |
| `REFRESH_TOKEN_EXPIRED` | 401 | Refresh токен истёк |
| `ACCOUNT_NOT_FOUND` | 404 | Аккаунт не найден |

---

## Требования к данным

### Пароль

| Требование | Значение |
|------------|----------|
| Минимальная длина | 8 символов |
| Заглавные буквы | Минимум 1 (A-Z) |
| Строчные буквы | Минимум 1 (a-z) |
| Цифры | Минимум 1 (0-9) |

### Email

| Требование | Значение |
|------------|----------|
| Формат | RFC 5321 |
| Максимальная длина | 254 символа |

---

## Хранение данных

### Redis

#### Refresh Tokens

```
Ключ:   refresh:{refresh_token}
Значение: {account_id}
TTL:    14 дней
```

### PostgreSQL

#### Таблица `accounts`

```sql
CREATE TABLE accounts (
    id              UUID PRIMARY KEY,
    email           VARCHAR(254) NOT NULL UNIQUE,
    password        VARCHAR(72) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_accounts_email ON accounts(email);
```

> **Примечание:** UUID генерируется на бэкенде (Go) перед вставкой.

---

## Безопасность

### JWT Claims

#### Access Token

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

#### Refresh Token

```json
{
  "iss": "auth-service",
  "sub": "{account_id}",
  "iat": 1705312200,
  "exp": 1706521800,
  "type": "refresh"
}
```

#### Описание claim'ов

| Claim | Описание |
|-------|----------|
| `iss` | Название сервиса (`auth-service`) |
| `sub` | ID аккаунта (UUID) |
| `email` | Email аккаунта |
| `iat` | Время выпуска токена (Unix timestamp) |
| `exp` | Время истечения токена (Unix timestamp) |
| `type` | Тип токена (`access` или `refresh`) |

### Headers

| Header | Значение |
|--------|----------|
| `Authorization` | `Bearer <access_token>` |
| `Content-Type` | `application/json` |

### CORS

| Настройка | Значение |
|-----------|----------|
| Allowed Origins | Конфигурируемо |
| Allowed Methods | `GET, POST, DELETE` |
| Allowed Headers | `Authorization, Content-Type` |
| Max Age | `86400` (24 часа) |

### Rate Limiting

| Endpoint | Лимит |
|----------|-------|
| `/api/auth/register` | 5 запросов/минуту |
| `/api/auth/login` | 10 запросов/минуту |
| `/api/auth/refresh` | 30 запросов/минуту |
| `/api/auth/logout` | 60 запросов/минуту |

---

## Разработка

### Taskfile Commands

```bash
# Генерация Proto (gRPC + REST + Swagger)
task proto:gen

# Линтинг Proto
task proto:lint

# Обновление Proto зависимостей
task proto:deps

# Форматирование кода
task format

# Линтинг Go кода
task lint

# Очистка зависимостей
task tidy

# Генерация DI кода
task wire:gen
```

### JWT Сервис

Пакет `pkg/jwt` предоставляет сервис для работы с JWT токенами. Может быть экспортирован в другие проекты.

**Пример использования:**

```go
import "auth-microservice/pkg/jwt"

// Создание сервиса
service, _ := jwt.NewService(jwt.Config{
    SecretKey:       "your-secret-key",
    AccessTokenTTL:  15 * time.Minute,
    RefreshTokenTTL: 14 * 24 * time.Hour,
    Issuer:          "auth-service",
})

// Генерация токенов
tokens, _ := service.GenerateTokens(accountID, email)

// Валидация токена
claims, _ := service.ValidateAccessToken(tokens.AccessToken)
```

📖 **Полная документация:** [pkg/jwt/README.md](../pkg/jwt/README.md)

### Зависимости

#### Внешние библиотеки

```go
github.com/dfsfGfd/redis-connect        // Redis клиент
github.com/dfsfGfd/postgresql-connect   // PostgreSQL клиент (pgx)
github.com/google/uuid                  // UUID генерация (Go)
```

#### Proto зависимости

```yaml
deps:
  - buf.build/googleapis/googleapis
  - buf.build/grpc-ecosystem/grpc-gateway
```

### Proto Files Structure

```
proto/
├── auth/
│   └── v1/
│       └── auth.proto          # Service definition
├── buf.yaml                    # Buf configuration
└── buf.gen.yaml                # Generation plugins
```

---

## Ссылки

- [README](../README.md) — Основная документация проекта
- [Proto файл](../proto/auth/v1/auth.proto) — gRPC контракты
