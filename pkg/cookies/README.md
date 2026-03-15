# Cookies Service

Пакет для работы с HTTP cookie. Может быть экспортирован и использован в других проектах.

## 📦 Установка

Пакет использует стандартную библиотеку `net/http`:

```bash
# Нет внешних зависимостей
```

## 🚀 Быстрый старт

### Создание сервиса

```go
import "auth-microservice/pkg/cookies"

// Конфигурация для продакшена
config := cookies.Config{
    Secure:   true,                              // Только HTTPS
    HTTPOnly: true,                              // Недоступны для JavaScript
    SameSite: http.SameSiteStrictMode,           // Защита от CSRF
    Domain:   "auth.example.com",                // Домен
    Path:     "/",                               // Путь
    MaxAge:   1209600,                           // 14 дней в секундах
}

// Создание сервиса
cookieService := cookies.NewService(config)
```

### Установка refresh токена

```go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // ... аутентификация ...
    
    refreshToken := "generated_refresh_token"
    
    // Установка cookie с refresh токеном
    cookieService.SetRefreshToken(w, refreshToken)
    
    // Отправка access токена в ответе
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": accessToken,
    })
}
```

### Получение refresh токена

```go
func RefreshHandler(w http.ResponseWriter, r *http.Request) {
    // Получение refresh токена из cookie
    refreshToken, err := cookieService.GetRefreshToken(r)
    if err != nil {
        http.Error(w, "refresh token not found", http.StatusUnauthorized)
        return
    }
    
    // ... обновление токенов ...
}
```

### Удаление cookie (logout)

```go
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // ... отзыв токена ...
    
    // Удаление всех cookie
    cookieService.DeleteAll(w)
    
    w.WriteHeader(http.StatusOK)
}
```

## 📖 API

### Типы

#### `Config`

Конфигурация cookie сервиса.

```go
type Config struct {
    Secure   bool            // Флаг Secure (HTTPS only)
    HTTPOnly bool            // Флаг HttpOnly (no JavaScript access)
    SameSite http.SameSite   // SameSite policy (Strict/Lax/None)
    Domain   string          // Домен cookie (опционально)
    Path     string          // Путь cookie
    MaxAge   int             // Время жизни в секундах
}
```

**Поля:**

| Поле | Тип | По умолчанию | Описание |
|------|-----|--------------|----------|
| `Secure` | `bool` | `false` | Передача только по HTTPS |
| `HTTPOnly` | `bool` | `false` | Недоступны для JavaScript |
| `SameSite` | `http.SameSite` | `http.SameSiteDefaultMode` | CSRF защита |
| `Domain` | `string` | `""` | Домен cookie |
| `Path` | `string` | `"/"` | Путь cookie |
| `MaxAge` | `int` | `1209600` (14 дней) | Время жизни в секундах |

### Функции

#### `NewService(config Config) *Service`

Создаёт новый cookie сервис.

```go
cookieService := cookies.NewService(cookies.Config{
    Secure:   true,
    HTTPOnly: true,
    SameSite: http.SameSiteStrictMode,
})
```

### Методы сервиса

#### `SetRefreshToken(w http.ResponseWriter, token string)`

Устанавливает refresh токен в cookie.

```go
cookieService.SetRefreshToken(w, refreshToken)
```

#### `GetRefreshToken(r *http.Request) (string, error)`

Получает refresh токен из cookie.

```go
refreshToken, err := cookieService.GetRefreshToken(r)
if err != nil {
    // Cookie не найдена
}
```

#### `DeleteRefreshToken(w http.ResponseWriter)`

Удаляет refresh токен из cookie.

```go
cookieService.DeleteRefreshToken(w)
```

#### `SetAccessToken(w http.ResponseWriter, token string, expiresIn int)`

Устанавливает access токен в cookie (опционально).

```go
cookieService.SetAccessToken(w, accessToken, 900) // 15 минут
```

#### `GetAccessToken(r *http.Request) (string, error)`

Получает access токен из cookie.

```go
accessToken, err := cookieService.GetAccessToken(r)
```

#### `DeleteAccessToken(w http.ResponseWriter)`

Удаляет access токен из cookie.

```go
cookieService.DeleteAccessToken(w)
```

#### `DeleteAll(w http.ResponseWriter)`

Удаляет все cookie.

```go
cookieService.DeleteAll(w)
```

#### `GetConfig() Config`

Возвращает конфигурацию сервиса.

```go
config := cookieService.GetConfig()
```

## 🔒 Безопасность

### Рекомендации для продакшена

```go
config := cookies.Config{
    Secure:   true,                    // Обязательно для HTTPS
    HTTPOnly: true,                    // Защита от XSS
    SameSite: http.SameSiteStrictMode, // Защита от CSRF
    Domain:   "auth.example.com",      // Ограничить доменом
    Path:     "/",
    MaxAge:   1209600,                 // 14 дней
}
```

### Флаги безопасности

| Флаг | Описание | Рекомендация |
|------|----------|--------------|
| `Secure` | Передача только по HTTPS | ✅ Всегда в продакшене |
| `HttpOnly` | Недоступны для JavaScript | ✅ Всегда включён |
| `SameSite=Strict` | Защита от CSRF | ✅ Для refresh токенов |
| `SameSite=Lax` | Разрешает навигацию GET | ⚠️ Для access токенов |

### Development конфигурация

```go
// Только для локальной разработки (localhost)
config := cookies.Config{
    Secure:   false,                   // HTTP локально
    HTTPOnly: true,                    // Всё равно защищено
    SameSite: http.SameSiteLaxMode,
    Path:     "/",
    MaxAge:   1209600,
}
```

## 📝 Примеры использования

### Middleware для аутентификации

```go
func AuthMiddleware(cookieService *cookies.Service, jwtService *jwt.Service) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Попытка получить токен из cookie
            refreshToken, err := cookieService.GetRefreshToken(r)
            if err != nil {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            
            // Валидация токена
            claims, err := jwtService.ValidateRefreshToken(refreshToken)
            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }
            
            // Добавление в контекст
            ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Полный handler входа

```go
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
    TokenType   string `json:"token_type"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    
    // Аутентификация
    tokens, err := h.authService.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        h.logger.Error("login failed", "error", err)
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }
    
    // Установка refresh токена в cookie
    h.cookieService.SetRefreshToken(w, tokens.RefreshToken)
    
    // Отправка access токена в ответе
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(LoginResponse{
        AccessToken: tokens.AccessToken,
        ExpiresIn:   tokens.ExpiresIn,
        TokenType:   "Bearer",
    })
}
```

### Handler обновления токенов

```go
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
    // Получение refresh токена из cookie
    refreshToken, err := h.cookieService.GetRefreshToken(r)
    if err != nil {
        h.logger.Warn("refresh token not found", "error", err)
        http.Error(w, "refresh token required", http.StatusUnauthorized)
        return
    }
    
    // Обновление токенов
    tokens, err := h.authService.Refresh(r.Context(), refreshToken)
    if err != nil {
        h.logger.Error("refresh failed", "error", err)
        http.Error(w, "invalid refresh token", http.StatusUnauthorized)
        return
    }
    
    // Установка нового refresh токена
    h.cookieService.SetRefreshToken(w, tokens.RefreshToken)
    
    // Отправка нового access токена
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(LoginResponse{
        AccessToken: tokens.AccessToken,
        ExpiresIn:   tokens.ExpiresIn,
        TokenType:   "Bearer",
    })
}
```

### Handler выхода

```go
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
    // Получение refresh токена
    refreshToken, err := h.cookieService.GetRefreshToken(r)
    if err == nil {
        // Отзыв токена в Redis
        _ = h.authService.Logout(r.Context(), refreshToken)
    }
    
    // Удаление cookie
    h.cookieService.DeleteAll(w)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
```

## 🎯 Best Practices

### ✅ Делайте

- Используйте `HttpOnly` для защиты от XSS
- Используйте `Secure` в продакшене
- Используйте `SameSite=Strict` для refresh токенов
- Ограничивайте `Domain` и `Path`
- Устанавливайте разумный `MaxAge`

### ❌ Не делайте

- Не храните чувствительные данные в cookie (кроме токенов)
- Не используйте `SameSite=None` без необходимости
- Не устанавливайте слишком долгий `MaxAge`
- Не передавайте refresh токены в теле ответа

## 🧪 Тесты

```bash
go test ./pkg/cookies/... -v
```

### Пример тестирования

```go
func TestCookieService(t *testing.T) {
    svc := cookies.NewService(cookies.Config{
        Secure:   true,
        HTTPOnly: true,
        SameSite: http.SameSiteStrictMode,
    })
    
    t.Run("set and get refresh token", func(t *testing.T) {
        w := httptest.NewRecorder()
        svc.SetRefreshToken(w, "test_token")
        
        r := httptest.NewRequest(http.MethodPost, "/refresh", nil)
        r.AddCookie(w.Result().Cookies()[0])
        
        token, err := svc.GetRefreshToken(r)
        require.NoError(t, err)
        assert.Equal(t, "test_token", token)
    })
}
```

## 📦 Экспорт в другие проекты

Пакет может быть использован в других проектах:

```go
import "github.com/dfsfGfd/auth-microservice/pkg/cookies"

// Использование в другом проекте
cookieService := cookies.NewService(cookies.Config{
    Secure:   true,
    HTTPOnly: true,
    SameSite: http.SameSiteStrictMode,
    Domain:   "my-app.com",
})
```

## 📚 Ссылки

- [RFC 6265 - HTTP State Management Mechanism](https://tools.ietf.org/html/rfc6265)
- [OWASP Cookie Security](https://cheatsheetseries.owasp.org/cheatsheets/Cookies_Cheat_Sheet.html)
- [Основное README](../../README.md)
