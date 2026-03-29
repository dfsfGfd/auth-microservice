# API Documentation

> Auth Microservice API: gRPC + REST (grpc-gateway)

---

## 📋 Endpoints

| Метод | HTTP | gRPC | Описание |
|-------|------|------|----------|
| `Register` | `POST /api/v1/auth/register` | `Register` | Регистрация |
| `Login` | `POST /api/v1/auth/login` | `Login` | Вход |
| `Logout` | `POST /api/v1/auth/logout` | `Logout` | Выход |
| `Refresh` | `POST /api/v1/auth/refresh` | `Refresh` | Обновление токена |

---

## Формат ответов

```json
{
  "status_code": 200,
  "message": "Success",
  "data": { ... }
}
```

---

## 1. Регистрация

```http
POST /api/v1/auth/register
Content-Type: application/json
```

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "status_code": 200,
  "message": "Account registered successfully",
  "data": {
    "account_id": "296494707175849985",
    "email": "user@example.com",
    "created_at": "2024-01-15T10:30:00Z"
  }
}
```

**Ошибки:**

| Код | HTTP | Сообщение |
|-----|------|-----------|
| `INVALID_EMAIL` | 400 | Неверный формат email |
| `EMAIL_EXISTS` | 409 | Email уже зарегистрирован |
| `PASSWORD_TOO_SHORT` | 400 | Пароль короче 8 символов |

---

## Требования к паролю

Согласно **NIST 800-63B**, длина важнее сложности:

- ✅ Минимум **8 символов**
- ✅ Без требований к заглавным/строчным буквам или цифрам

**Примеры валидных паролей:**
```
password123    ✅ (нет заглавной)
PASSWORD123    ✅ (нет строчной)
Password       ✅ (нет цифры)
abcdefgh       ✅ (только буквы)
```

---

## 2. Вход

```http
POST /api/v1/auth/login
Content-Type: application/json
```

**Request:**
```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

**Response (200):**
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

**Ошибки:**

| Код | HTTP | Сообщение |
|-----|------|-----------|
| `INVALID_CREDENTIALS` | 401 | Неверный email или пароль |

> **Примечание:** Все ошибки аутентификации возвращают `invalid credentials` для защиты от user enumeration.

**Rate Limiting:**
- Лимит: 10 запросов/мин
- Заголовки: `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `X-RateLimit-Reset`

---

## 3. Выход

```http
POST /api/v1/auth/logout
Content-Type: application/json
Authorization: Bearer <access_token>
```

**Request:**
```json
{
  "refresh_token": "dGhpc2lzYXJlZnJlc2h0b2tlbg..."
}
```

**Response (200):**
```json
{
  "status_code": 200,
  "message": "Logout successful",
  "data": {
    "success": true
  }
}
```

**Ошибки:**

| Код | HTTP | Сообщение |
|-----|------|-----------|
| `UNAUTHORIZED` | 401 | Access токен недействителен |
| `TOKEN_EXPIRED` | 401 | Access токен истёк |

---

## 4. Обновление токена

```http
POST /api/v1/auth/refresh
Content-Type: application/json
```

**Request:**
```json
{
  "refresh_token": "dGhpc2lzYXJlZnJlc2h0b2tlbg..."
}
```

**Response (200):**
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

**Ошибки:**

| Код | HTTP | Сообщение |
|-----|------|-----------|
| `INVALID_REFRESH_TOKEN` | 401 | Refresh токен недействителен |
| `REFRESH_TOKEN_EXPIRED` | 401 | Refresh токен истёк |
| `ACCOUNT_NOT_FOUND` | 404 | Аккаунт не найден |

---

## Требования к данным

### Пароль

- Минимум 8 символов
- 1 заглавная буква (A-Z)
- 1 строчная буква (a-z)
- 1 цифра (0-9)

### Email

- Формат: RFC 5321
- Максимум: 254 символа

---

## JWT Claims

**Access Token:**
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

**Refresh Token:**
```json
{
  "iss": "auth-service",
  "sub": "{account_id}",
  "iat": 1705312200,
  "exp": 1706521800,
  "type": "refresh"
}
```

| Claim | Описание |
|-------|----------|
| `iss` | Issuer (auth-service) |
| `sub` | Account ID (Snowflake ID строкой) |
| `email` | Email аккаунта |
| `iat` | Время выпуска |
| `exp` | Время истечения |
| `type` | Тип токена |

---

## Rate Limiting

| Endpoint | Лимит/мин |
|----------|-----------|
| `/api/v1/auth/register` | 5 |
| `/api/v1/auth/login` | 10 |
| `/api/v1/auth/refresh` | 30 |
| `/api/v1/auth/logout` | 60 |

**Заголовки ответа:**
- `X-RateLimit-Limit` — максимум запросов
- `X-RateLimit-Remaining` — осталось запросов
- `X-RateLimit-Reset` — timestamp сброса
- `Retry-After` — секунд до следующего запроса (429)

**429 Too Many Requests:**
```json
{
  "status_code": 429,
  "message": "rate limit exceeded",
  "data": null
}
```

---

## Ссылки

- [README](docs/README.md) — основная документация
- [Config](docs/config.md) — настройка .env
- [Proto](proto/auth/v1/auth.proto) — gRPC контракты
