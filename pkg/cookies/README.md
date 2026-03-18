# Cookies Package

Утилиты для работы с HTTP cookie.

---

## 🚀 Использование

```go
import "auth-microservice/pkg/cookies"

// Создание
cookieService := cookies.NewService(cookies.Config{
    Secure:   true,
    HTTPOnly: true,
    SameSite: http.SameSiteStrictMode,
    Domain:   "auth.example.com",
    Path:     "/",
    MaxAge:   1209600, // 14 дней
})

// Установка refresh токена
cookieService.SetRefreshToken(w, refreshToken)

// Получение
refreshToken, err := cookieService.GetRefreshToken(r)

// Удаление
cookieService.DeleteRefreshToken(w)
```

---

## Config

```go
type Config struct {
    Secure   bool              // HTTPS only
    HTTPOnly bool              // no JavaScript access
    SameSite http.SameSite     // Strict/Lax/None
    Domain   string            // домен (опционально)
    Path     string            // путь
    MaxAge   int               // секунды
}
```

---

## Методы

```go
// Создание
func NewService(config Config) *Service

// Refresh токен
func (s *Service) SetRefreshToken(w http.ResponseWriter, token string)
func (s *Service) GetRefreshToken(r *http.Request) (string, error)
func (s *Service) DeleteRefreshToken(w http.ResponseWriter)

// Access токен (опционально)
func (s *Service) SetAccessToken(w http.ResponseWriter, token string, expiresIn int)
func (s *Service) GetAccessToken(r *http.Request) (string, error)
func (s *Service) DeleteAccessToken(w http.ResponseWriter)

// Удалить все
func (s *Service) DeleteAll(w http.ResponseWriter)
```

---

## 🔒 Безопасность

### Production

```go
cookies.Config{
    Secure:   true,                    // HTTPS only
    HTTPOnly: true,                    // защита от XSS
    SameSite: http.SameSiteStrictMode, // защита от CSRF
}
```

### Development

```go
cookies.Config{
    Secure:   false,                   // HTTP локально
    HTTPOnly: true,
    SameSite: http.SameSiteLaxMode,
}
```

---

## 🧪 Тесты

```bash
go test ./pkg/cookies/... -v
```
