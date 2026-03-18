# JWT Service

Пакет для работы с JWT токенами.

---

## 🚀 Использование

```go
import "auth-microservice/pkg/jwt"

// Создание сервиса
service, err := jwt.NewService(jwt.Config{
    SecretKey:       "your-secret-key-min-32-chars",
    AccessTokenTTL:  15 * time.Minute,
    RefreshTokenTTL: 14 * 24 * time.Hour,
    Issuer:          "auth-service",
})

// Генерация токенов
tokens, err := service.GenerateTokens(accountID, email)

// Валидация
claims, err := service.ValidateAccessToken(tokens.AccessToken)
```

---

## API

### Config

```go
type Config struct {
    SecretKey       string        // Минимум 32 символа
    AccessTokenTTL  time.Duration // 15 минут
    RefreshTokenTTL time.Duration // 14 дней
    Issuer          string        // "auth-service"
}
```

### TokenPair

```go
type TokenPair struct {
    AccessToken  string // JWT access токен
    RefreshToken string // JWT refresh токен
    ExpiresIn    int64  // Время жизни (сек)
    TokenType    string // "Bearer"
}
```

### Методы

```go
// Создание сервиса
func NewService(config Config) (*Service, error)

// Генерация пары токенов
func (s *Service) GenerateTokens(accountID, email string) (*TokenPair, error)

// Валидация токена
func (s *Service) ValidateToken(tokenString string) (*Claims, error)

// Валидация access токена
func (s *Service) ValidateAccessToken(tokenString string) (*Claims, error)

// Валидация refresh токена
func (s *Service) ValidateRefreshToken(tokenString string) (*Claims, error)
```

---

## Claims

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

| Claim | Описание |
|-------|----------|
| `iss` | Issuer (название сервиса) |
| `sub` | Account ID (UUID) |
| `email` | Email аккаунта |
| `iat` | Время выпуска |
| `exp` | Время истечения |
| `type` | Тип токена (access/refresh) |

---

## 🔒 Безопасность

### Генерация SecretKey

```bash
openssl rand -base64 32
```

### Требования

- Минимум 32 символа
- Хранить в secure storage
- Не коммитить в git

---

## 🧪 Тесты

```bash
go test ./pkg/jwt/... -v
```
