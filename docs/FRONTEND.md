# Frontend Integration Guide

> Быстрый старт для подключения Auth API к вашему frontend-приложению

---

## 📋 Оглавление

- [Быстрый старт](#-быстрый-старт)
- [API Endpoints](#-api-endpoints)
- [Примеры кода](#-примеры-кода)
- [Обработка ошибок](#-обработка-ошибок)
- [Хранение токенов](#-хранение-токенов)
- [TypeScript типы](#-typescript-типы)

---

## 🚀 Быстрый старт

### 1. Базовый URL

```typescript
const API_BASE = 'http://localhost:8080/api/v1/auth'
// Production: https://your-domain.com/api/v1/auth
```

### 2. Первый запрос (Регистрация)

```typescript
const response = await fetch(`${API_BASE}/register`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'Password123'
  })
})

const data = await response.json()
console.log(data.data.accountId) // "296502646347399169"
```

### 3. Вход и сохранение токена

```typescript
const response = await fetch(`${API_BASE}/login`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'Password123'
  })
})

const { data } = await response.json()

// Сохраняем токены
localStorage.setItem('accessToken', data.accessToken)
localStorage.setItem('refreshToken', data.refreshToken)
```

### 4. Авторизованный запрос

```typescript
const accessToken = localStorage.getItem('accessToken')

const response = await fetch(`${API_BASE}/refresh`, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${accessToken}`
  },
  body: JSON.stringify({
    refresh_token: localStorage.getItem('refreshToken')
  })
})
```

---

## 📡 API Endpoints

### POST `/register` — Регистрация

**Request:**
```typescript
{
  email: string,      // 1-254 символов, формат email
  password: string    // минимум 8 символов
}
```

**Response (200):**
```typescript
{
  statusCode: 200,
  message: "Account registered successfully",
  data: {
    accountId: string,    // Snowflake ID (строка)
    email: string,
    createdAt: string     // ISO 8601
  }
}
```

**Ошибки:**
```typescript
// 400 Bad Request
{
  code: 3,
  message: "invalid email"  // или "invalid password"
}

// 409 Conflict
{
  code: 6,
  message: "account already exists"
}
```

---

### POST `/login` — Вход

**Request:**
```typescript
{
  email: string,
  password: string
}
```

**Response (200):**
```typescript
{
  statusCode: 200,
  message: "Login successful",
  data: {
    accessToken: string,    // JWT, истекает через 15 мин
    refreshToken: string,   // JWT, истекает через 30 дней
    expiresIn: number,      // 900 (секунд)
    tokenType: "Bearer"
  }
}
```

**Ошибки:**
```typescript
// 401 Unauthorized
{
  code: 16,
  message: "invalid credentials"  // Всегда одинаковое сообщение
}
```

---

### POST `/logout` — Выход

**Request:**
```typescript
{
  refresh_token: string  // Refresh токен из localStorage
}
```

**Response (200):**
```typescript
{
  statusCode: 200,
  message: "Logout successful",
  data: {
    success: true
  }
}
```

---

### POST `/refresh` — Обновление токена

**Request:**
```typescript
{
  refresh_token: string
}
```

**Response (200):**
```typescript
{
  statusCode: 200,
  message: "Token refreshed successfully",
  data: {
    accessToken: string,
    refreshToken: string,   // Новый refresh токен
    expiresIn: number,
    tokenType: "Bearer"
  }
}
```

**Ошибки:**
```typescript
// 401 Unauthorized
{
  code: 16,
  message: "invalid token"       // Токен недействителен
  // или "token expired"         // Истёк срок действия
  // или "refresh token not found"
}
```

---

## 💻 Примеры кода

### React Hook

```typescript
// hooks/useAuth.ts
import { useState, useCallback } from 'react'

const API_BASE = '/api/v1/auth'

interface Tokens {
  accessToken: string
  refreshToken: string
  expiresIn: number
  tokenType: string
}

export function useAuth() {
  const [tokens, setTokens] = useState<Tokens | null>(() => {
    const saved = localStorage.getItem('tokens')
    return saved ? JSON.parse(saved) : null
  })

  const register = useCallback(async (email: string, password: string) => {
    const res = await fetch(`${API_BASE}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    })
    
    const data = await res.json()
    if (!res.ok) throw new Error(data.message)
    return data.data
  }, [])

  const login = useCallback(async (email: string, password: string) => {
    const res = await fetch(`${API_BASE}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    })
    
    const data = await res.json()
    if (!res.ok) throw new Error(data.message)
    
    setTokens(data.data)
    localStorage.setItem('tokens', JSON.stringify(data.data))
    return data.data
  }, [])

  const logout = useCallback(async () => {
    if (!tokens?.refreshToken) return
    
    await fetch(`${API_BASE}/logout`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: tokens.refreshToken })
    })
    
    setTokens(null)
    localStorage.removeItem('tokens')
  }, [tokens])

  const refreshToken = useCallback(async () => {
    if (!tokens?.refreshToken) throw new Error('No refresh token')
    
    const res = await fetch(`${API_BASE}/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: tokens.refreshToken })
    })
    
    const data = await res.json()
    if (!res.ok) throw new Error(data.message)
    
    setTokens(data.data)
    localStorage.setItem('tokens', JSON.stringify(data.data))
    return data.data
  }, [tokens])

  return { tokens, register, login, logout, refreshToken }
}
```

### Использование в компоненте

```typescript
// components/LoginForm.tsx
import { useAuth } from '@/hooks/useAuth'

export function LoginForm() {
  const { login } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await login(email, password)
      // Перенаправление на главную
    } catch (err) {
      setError('Неверный email или пароль')
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Пароль"
      />
      {error && <p className="error">{error}</p>}
      <button type="submit">Войти</button>
    </form>
  )
}
```

### Axios Interceptor (авто-обновление токена)

```typescript
// api/auth.ts
import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1/auth',
  headers: { 'Content-Type': 'application/json' }
})

// Авто-обновление токена
let isRefreshing = false
let failedQueue: Array<{
  resolve: (value: unknown) => void
  reject: (reason?: unknown) => void
}> = []

const processQueue = (error: Error | null, token: string | null = null) => {
  failedQueue.forEach(prom => {
    if (error) prom.reject(error)
    else prom.resolve(token)
  })
  failedQueue = []
}

api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        })
          .then(token => {
            originalRequest.headers['Authorization'] = `Bearer ${token}`
            return api(originalRequest)
          })
          .catch(err => Promise.reject(err))
      }

      originalRequest._retry = true
      isRefreshing = true

      const refreshToken = localStorage.getItem('refreshToken')
      
      try {
        const response = await api.post('/refresh', { refresh_token: refreshToken })
        const { accessToken, refreshToken: newRefreshToken } = response.data.data

        localStorage.setItem('accessToken', accessToken)
        localStorage.setItem('refreshToken', newRefreshToken)

        processQueue(null, accessToken)
        
        originalRequest.headers['Authorization'] = `Bearer ${accessToken}`
        return api(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError as Error, null)
        localStorage.removeItem('accessToken')
        localStorage.removeItem('refreshToken')
        window.location.href = '/login'
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)

export default api
```

### Vue 3 Composable

```typescript
// composables/useAuth.ts
import { ref, computed } from 'vue'

export function useAuth() {
  const tokens = ref<{
    accessToken: string
    refreshToken: string
  } | null>(JSON.parse(localStorage.getItem('tokens') || 'null'))

  const isAuthenticated = computed(() => !!tokens.value)

  const login = async (email: string, password: string) => {
    const res = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    })
    
    const data = await res.json()
    if (!res.ok) throw new Error(data.message)
    
    tokens.value = data.data
    localStorage.setItem('tokens', JSON.stringify(data.data))
    return data.data
  }

  const logout = async () => {
    if (tokens.value?.refreshToken) {
      await fetch('/api/v1/auth/logout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refresh_token: tokens.value.refreshToken })
      })
    }
    
    tokens.value = null
    localStorage.removeItem('tokens')
  }

  return { tokens, isAuthenticated, login, logout }
}
```

---

## ❌ Обработка ошибок

### gRPC коды → HTTP статусы

| gRPC Code | HTTP | Значение |
|-----------|------|----------|
| `OK (0)` | 200 | Успех |
| `INVALID_ARGUMENT (3)` | 400 | Неверные данные |
| `NOT_FOUND (5)` | 404 | Не найдено |
| `ALREADY_EXISTS (6)` | 409 | Уже существует |
| `UNAUTHENTICATED (16)` | 401 | Не авторизован |
| `INTERNAL (13)` | 500 | Внутренняя ошибка |

### Пример обработки

```typescript
async function handleApiCall() {
  try {
    const response = await fetch('/api/v1/auth/login', { ... })
    const data = await response.json()
    
    if (!response.ok) {
      switch (response.status) {
        case 400:
          showError('Неверный email или пароль')
          break
        case 401:
          showError('Неверные учётные данные')
          break
        case 409:
          showError('Аккаунт уже существует')
          break
        case 500:
          showError('Ошибка сервера, попробуйте позже')
          break
        default:
          showError(data.message || 'Произошла ошибка')
      }
      return
    }
    
    // Успех
    console.log(data.data)
  } catch (error) {
    console.error('Network error:', error)
    showError('Нет соединения с сервером')
  }
}
```

---

## 🔐 Хранение токенов

### Варианты

| Способ | Безопасность | Удобство | Рекомендация |
|--------|--------------|----------|--------------|
| **localStorage** | ⚠️ Средняя | ✅ Удобно | Для MVP |
| **sessionStorage** | ⚠️ Средняя | ✅ Удобно | Для сессий |
| **HttpOnly Cookie** | ✅ Высокая | ⚠️ Сложнее | Для production |
| **Memory + Refresh** | ✅ Высокая | ⚠️ Сложнее | Для SPA |

### Recommended: Memory + Refresh Token

```typescript
// stores/auth.ts
class AuthStore {
  private accessToken: string | null = null
  private refreshToken: string | null = null

  constructor() {
    // Загружаем refresh токен при старте
    this.refreshToken = localStorage.getItem('refreshToken')
  }

  async login(email: string, password: string) {
    const res = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password })
    })
    
    const { data } = await res.json()
    this.accessToken = data.accessToken
    this.refreshToken = data.refreshToken
    
    // Сохраняем только refresh токен
    localStorage.setItem('refreshToken', data.refreshToken)
  }

  getAccessToken() {
    return this.accessToken
  }

  async ensureValidToken() {
    // Если access токен истёк, обновляем
    if (!this.accessToken) {
      await this.refresh()
    }
    return this.accessToken
  }

  private async refresh() {
    if (!this.refreshToken) throw new Error('No refresh token')
    
    const res = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: this.refreshToken })
    })
    
    const { data } = await res.json()
    this.accessToken = data.accessToken
    this.refreshToken = data.refreshToken
    localStorage.setItem('refreshToken', data.refreshToken)
  }

  logout() {
    this.accessToken = null
    this.refreshToken = null
    localStorage.removeItem('refreshToken')
  }
}

export const authStore = new AuthStore()
```

---

## 📝 TypeScript Типы

```typescript
// types/auth.ts

export interface RegisterRequest {
  email: string
  password: string
}

export interface RegisterResponse {
  statusCode: number
  message: string
  data: {
    accountId: string
    email: string
    createdAt: string  // ISO 8601
  }
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  statusCode: number
  message: string
  data: {
    accessToken: string
    refreshToken: string
    expiresIn: number
    tokenType: 'Bearer'
  }
}

export interface LogoutRequest {
  refresh_token: string
}

export interface LogoutResponse {
  statusCode: number
  message: string
  data: {
    success: boolean
  }
}

export interface RefreshRequest {
  refresh_token: string
}

export interface RefreshResponse {
  statusCode: number
  message: string
  data: {
    accessToken: string
    refreshToken: string
    expiresIn: number
    tokenType: 'Bearer'
  }
}

export interface ApiError {
  code: number
  message: string
  details?: Array<{
    type: string
    value: string
  }>
}
```

---

## 🧪 Тестирование

### cURL примеры

```bash
# Регистрация
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Password123"}'

# Вход
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"Password123"}'

# Обновление токена
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbGciOiJIUzI1NiIs..."}'

# Выход
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"eyJhbGciOiJIUzI1NiIs..."}'
```

### Postman Collection

Импортируйте в Postman:

```json
{
  "info": {
    "name": "Auth API",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Register",
      "request": {
        "method": "POST",
        "header": [{ "key": "Content-Type", "value": "application/json" }],
        "body": {
          "mode": "raw",
          "raw": "{\"email\":\"test@example.com\",\"password\":\"Password123\"}"
        },
        "url": {
          "raw": "{{baseUrl}}/api/v1/auth/register",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "auth", "register"]
        }
      }
    },
    {
      "name": "Login",
      "request": {
        "method": "POST",
        "header": [{ "key": "Content-Type", "value": "application/json" }],
        "body": {
          "mode": "raw",
          "raw": "{\"email\":\"test@example.com\",\"password\":\"Password123\"}"
        },
        "url": {
          "raw": "{{baseUrl}}/api/v1/auth/login",
          "host": ["{{baseUrl}}"],
          "path": ["api", "v1", "auth", "login"]
        }
      }
    }
  ],
  "variable": [
    {
      "key": "baseUrl",
      "value": "http://localhost:8080"
    }
  ]
}
```

---

## 📚 Ссылки

- [API Documentation](api.md) — Полная документация API
- [Configuration](config.md) — Настройка окружения
- [Docker Guide](../deploy/README.md) — Запуск в Docker

---

## 🆘 Поддержка

**Проблемы с подключением?**

1. Проверьте, что сервер запущен: `http://localhost:8080/health`
2. Убедитесь в правильности CORS настроек
3. Проверьте логи: `docker logs auth-service`

**Частые ошибки:**

| Ошибка | Причина | Решение |
|--------|---------|---------|
| `401 Unauthorized` | Токен истёк | Используйте `/refresh` endpoint |
| `400 Bad Request` | Неверный email/password | Проверьте валидацию на клиенте |
| `500 Internal Error` | Ошибка сервера | Проверьте логи сервера |
| `Network Error` | CORS / сервер недоступен | Проверьте настройки CORS |
