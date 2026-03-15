# JWT Service

Пакет для работы с JWT токенами. Может быть экспортирован и использован в других проектах.

## 📦 Установка

```bash
go get github.com/golang-jwt/jwt/v5
```

## 🚀 Быстрый старт

### Создание сервиса

```go
import "auth-microservice/pkg/jwt"

// Конфигурация
config := jwt.Config{
    SecretKey:       "your-secret-key",
    AccessTokenTTL:  15 * time.Minute,
    RefreshTokenTTL: 14 * 24 * time.Hour,
    Issuer:          "auth-service",
}

// Создание сервиса
service, err := jwt.NewService(config)
if err != nil {
    log.Fatal(err)
}
```

### Генерация токенов

```go
userID := "550e8400-e29b-41d4-a716-446655440000"
email := "user@example.com"
username := "testuser"

// Генерация пары токенов
tokens, err := service.GenerateTokens(userID, email, username)
if err != nil {
    log.Fatal(err)
}

// Использование
fmt.Println(tokens.AccessToken)   // Access токен
fmt.Println(tokens.RefreshToken)  // Refresh токен
fmt.Println(tokens.ExpiresIn)     // Время жизни в секундах (900)
fmt.Println(tokens.TokenType)     // "Bearer"
```

### Валидация токенов

```go
// Валидация любого токена
claims, err := service.ValidateToken(tokenString)
if err != nil {
    log.Fatal(err)
}

fmt.Println(claims.UserID)   // ID пользователя
fmt.Println(claims.Email)    // Email
fmt.Println(claims.Username) // Username
fmt.Println(claims.Type)     // Тип токена (access/refresh)
fmt.Println(claims.Issuer)   // Issuer сервиса
```

### Валидация access токена

```go
claims, err := service.ValidateAccessToken(accessToken)
if err != nil {
    // Токен невалиден или истёк
    log.Fatal(err)
}

// Токен валиден, это access токен
if claims.Type != jwt.AccessToken {
    log.Fatal("invalid token type")
}
```

### Валидация refresh токена

```go
claims, err := service.ValidateRefreshToken(refreshToken)
if err != nil {
    // Токен невалиден или истёк
    log.Fatal(err)
}

// Токен валиден, это refresh токен
if claims.Type != jwt.RefreshToken {
    log.Fatal("invalid token type")
}
```

## 📖 API

### Типы

#### `Config`

Конфигурация JWT сервиса.

```go
type Config struct {
    SecretKey       string        // Секретный ключ для подписи
    AccessTokenTTL  time.Duration // Время жизни access токена
    RefreshTokenTTL time.Duration // Время жизни refresh токена
    Issuer          string        // Название сервиса (iss claim)
}
```

#### `Claims`

Claims JWT токена.

```go
type Claims struct {
    UserID   string    `json:"sub"`
    Email    string    `json:"email,omitempty"`
    Username string    `json:"username,omitempty"`
    Type     TokenType `json:"type"`
    jwt.RegisteredClaims
}
```

#### `TokenPair`

Пара токенов.

```go
type TokenPair struct {
    AccessToken  string // Access токен
    RefreshToken string // Refresh токен
    ExpiresIn    int64  // Время жизни в секундах
    TokenType    string // "Bearer"
}
```

#### `TokenType`

Тип токена.

```go
type TokenType string

const (
    AccessToken  TokenType = "access"
    RefreshToken TokenType = "refresh"
)
```

### Функции

#### `NewService(config Config) (*Service, error)`

Создаёт новый JWT сервис.

```go
service, err := jwt.NewService(jwt.Config{
    SecretKey:       "secret",
    AccessTokenTTL:  15 * time.Minute,
    RefreshTokenTTL: 14 * 24 * time.Hour,
    Issuer:          "auth-service",
})
```

**Ошибки:**
- `SecretKey` пуст
- `AccessTokenTTL` <= 0
- `RefreshTokenTTL` <= 0
- `Issuer` пуст

### Методы сервиса

#### `GenerateTokens(userID, email, username string) (*TokenPair, error)`

Генерирует пару access и refresh токенов.

```go
tokens, err := service.GenerateTokens(userID, email, username)
```

#### `ValidateToken(tokenString string) (*Claims, error)`

Валидирует JWT токен и возвращает claims.

```go
claims, err := service.ValidateToken(tokenString)
```

**Ошибки:**
- `ErrInvalidToken` — токен невалиден
- `ErrExpiredToken` — токен истёк
- `ErrInvalidClaims` — неверные claims

#### `ValidateAccessToken(tokenString string) (*Claims, error)`

Валидирует access токен.

```go
claims, err := service.ValidateAccessToken(tokenString)
```

#### `ValidateRefreshToken(tokenString string) (*Claims, error)`

Валидирует refresh токен.

```go
claims, err := service.ValidateRefreshToken(tokenString)
```

#### `GetConfig() Config`

Возвращает конфигурацию сервиса.

```go
config := service.GetConfig()
```

## 🔒 Безопасность

### Требования к SecretKey

- Минимум 32 символа
- Используйте криптографически стойкую генерацию
- Храните в secure storage (Vault, AWS Secrets Manager)

```go
// Пример генерации
import "crypto/rand"

func generateSecretKey() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.StdEncoding.EncodeToString(bytes), nil
}
```

### JWT Claims структура

```json
{
  "iss": "auth-service",
  "sub": "{user_id}",
  "email": "{email}",
  "username": "{username}",
  "iat": 1705312200,
  "exp": 1705313100,
  "nbf": 1705312200,
  "type": "access"
}
```

| Claim | Описание |
|-------|----------|
| `iss` | Issuer (название сервиса) |
| `sub` | Subject (ID пользователя) |
| `email` | Email пользователя |
| `username` | Username пользователя |
| `iat` | Issued at (время выпуска) |
| `exp` | Expiration time (время истечения) |
| `nbf` | Not before (не действителен до) |
| `type` | Тип токена (`access` или `refresh`) |

## 🧪 Тесты

```bash
go test ./pkg/jwt/... -v
```

## 📝 Примеры использования

### Middleware для аутентификации

```go
func AuthMiddleware(service *jwt.Service) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "missing authorization header", http.StatusUnauthorized)
                return
            }

            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            
            claims, err := service.ValidateAccessToken(tokenString)
            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }

            // Добавляем claims в контекст
            ctx := context.WithValue(r.Context(), "claims", claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Refresh токенов

```go
func RefreshTokens(w http.ResponseWriter, r *http.Request, service *jwt.Service) {
    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    // Валидируем refresh токен
    claims, err := service.ValidateRefreshToken(req.RefreshToken)
    if err != nil {
        http.Error(w, "invalid refresh token", http.StatusUnauthorized)
        return
    }

    // Генерируем новую пару токенов
    tokens, err := service.GenerateTokens(claims.UserID, claims.Email, claims.Username)
    if err != nil {
        http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(tokens)
}
```

## 📦 Экспорт в другие проекты

Пакет может быть использован в других проектах:

```go
import "github.com/dfsfGfd/auth-microservice/pkg/jwt"

// Использование в другом проекте
service, _ := jwt.NewService(jwt.Config{
    SecretKey:       os.Getenv("JWT_SECRET"),
    AccessTokenTTL:  15 * time.Minute,
    RefreshTokenTTL: 14 * 24 * time.Hour,
    Issuer:          "my-service",
})
```
