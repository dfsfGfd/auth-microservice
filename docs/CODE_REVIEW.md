# Code Review — Auth Microservice

> Анализ кода на оверинженеринг и рекомендации по рефакторингу.

**Дата:** 2026-03-28  
**Версия:** v1.0.0

---

## 📊 Общая оценка

| Аспект | Оценка | Комментарий |
|--------|--------|-------------|
| **Структура проекта** | ✅ Отлично | Чёткое разделение на слои (DDD) |
| **DI через Wire** | ✅ Отлично | Явные зависимости, компиляция DI |
| **Value Objects** | ✅ Отлично | Email, PasswordHash с инвариантами |
| **Domain агрегат** | ✅ Отлично | Account с приватными полями |
| **Обработка ошибок** | ✅ Отлично | Доменные ошибки в `internal/errors` |
| **Безопасность** | ✅ Отлично | PasswordHash.String() → [REDACTED] |
| **Rate limiting** | ✅ Отлично | Redis-based sliding window |
| **Конфигурация** | ✅ Отлично | Только .env, нет YAML |

---

## ⚠️ Проблемы и оверинженеринг

### 1. **DI: Избыточная сложность** 🔴 Высокий приоритет

**Файл:** `internal/di/provider.go`

**Проблема:**
```go
type Application struct {
    Config      *config.Config      // ❌ Не используется в CleanUp
    Logger      *logger.Logger      // ❌ Не используется
    JWTService  *jwt.Service        // ❌ Не используется
    DB          *pgxpool.Pool       // ✅ Используется
    Redis       *goredis.Client     // ✅ Используется
    AccountRepo repository.AccountRepository // ❌ Не используется
    TokenCache  *token.RedisCache   // ❌ Не используется
    RateLimiter *middleware.RateLimiter      // ❌ Не используется
    AuthService *serviceAuth.AuthService     // ❌ Не используется
    AuthHandler authv1.AuthServiceServer     // ❌ Не используется
    // 11 полей, только 2 используются в CleanUp
}

func (a *Application) CleanUp() error {
    if a.DB != nil { a.DB.Close() }
    if a.Redis != nil { a.Redis.Close() }
    return nil
}
```

**Почему оверинженеринг:**
- ❌ 11 зависимостей в Application
- ❌ 9 полей не используются в CleanUp
- ❌ Wire генерирует 200+ строк кода
- ❌ Избыточная абстракция

**Решение:**
```go
// Упростить до:
type Application struct {
    DB    *pgxpool.Pool
    Redis *goredis.Client
}

func (a *Application) CleanUp() error {
    if a.DB != nil { a.DB.Close() }
    if a.Redis != nil { a.Redis.Close() }
    return nil
}

// Остальные зависимости инжектить напрямую в server.NewServer
```

**Файлы для изменения:**
- `internal/di/provider.go`
- `internal/di/wire.go` (перегенерировать)

---

### 2. **Model: SetID/SetCreatedAt/SetUpdatedAt** 🟡 Средний приоритет

**Файл:** `internal/model/account.go`

**Проблема:**
```go
func (a *Account) SetID(id int64) { a.id = id }
func (a *Account) SetCreatedAt(t time.Time) { a.createdAt = t }
func (a *Account) SetUpdatedAt(t time.Time) { a.updatedAt = t }
```

**Почему оверинженеринг:**
- ❌ Нарушают инкапсуляцию агрегата
- ❌ Используются **только** в конвертере из БД
- ❌ Могут быть заменены на конструктор

**Решение:**
```go
// Добавить новый конструктор:
func NewAccountFromDB(
    id int64,
    email *Email,
    passwordHash *PasswordHash,
    createdAt, updatedAt time.Time,
) *Account {
    return &Account{
        id:           id,
        email:        email,
        passwordHash: passwordHash,
        createdAt:    createdAt,
        updatedAt:    updatedAt,
    }
}

// Удалить методы SetID, SetCreatedAt, SetUpdatedAt
```

**Файлы для изменения:**
- `internal/model/account.go`
- `internal/repository/converter/account.go`

---

### 3. **Converter: Дублирование кода** 🟡 Средний приоритет

**Файл:** `internal/repository/converter/account.go`

**Проблема:**
```go
func AccountToDomain(db *dbmodel.Account) (*model.Account, error) {
    email := model.NewEmailFromDB(db.Email)
    passwordHash := model.NewPasswordHashFromString(db.PasswordHash)

    account, err := model.NewAccount(id, email, passwordHash)  // ← Создаёт с ID
    if err != nil {
        return nil, err
    }

    account.SetCreatedAt(db.CreatedAt) // ← Перезаписываем время
    account.SetUpdatedAt(db.UpdatedAt) // ← Перезаписываем время
    return account, nil
}
```

**Почему проблема:**
- ❌ `NewAccount` создаёт createdAt=now, потом перезаписываем
- ❌ Лишние аллокации (3 вызова вместо 1)

**Решение:**
```go
func AccountToDomain(db *dbmodel.Account) (*model.Account, error) {
    email := model.NewEmailFromDB(db.Email)
    passwordHash := model.NewPasswordHashFromString(db.PasswordHash)
    // Один конструктор вместо 3 вызовов
    return model.NewAccountFromDB(db.ID, email, passwordHash, db.CreatedAt, db.UpdatedAt), nil
}
```

**Файлы для изменения:**
- `internal/repository/converter/account.go`

---

### 4. **PasswordHash: Избыточный интерфейс** 🟡 Средний приоритет

**Файл:** `internal/model/password_hash.go`

**Проблема:**
```go
type PasswordHash string

func NewPasswordHash(hash string) (*PasswordHash, error) {
    if hash == "" { return nil, errs.ErrPasswordInvalid }
    if len(hash) < 53 || len(hash) > 72 { return nil, errs.ErrPasswordInvalid }
    h := PasswordHash(hash)
    return &h, nil  // ← Возвращаем pointer на string
}

func (p PasswordHash) Value() string { return string(p) }
func (p PasswordHash) String() string { return "[REDACTED]" }
```

**Почему оверинженеринг:**
- ❌ Pointer на string (`*PasswordHash`) — лишняя косвенность
- ❌ Метод `Value()` просто делает `string(p)`
- ❌ `String()` возвращает "[REDACTED]", но это не нужно (не используется в логах)

**Решение:**
```go
type PasswordHash string

func NewPasswordHash(hash string) (PasswordHash, error) {
    if len(hash) < 53 || len(hash) > 72 {
        return "", errs.ErrPasswordInvalid
    }
    return PasswordHash(hash), nil
}

// Value() не нужен — используем напрямую: string(passwordHash)
// String() можно оставить для безопасного логирования
```

**Файлы для изменения:**
- `internal/model/password_hash.go`
- `internal/model/account.go` (изменить тип поля)
- `internal/service/auth/*.go` (обновить вызовы)

---

### 5. **RateLimiter: Сложная логика для 4 endpoint'ов** 🟡 Средний приоритет

**Файл:** `internal/middleware/rate_limiter.go`

**Проблема:**
```go
type RateLimiterConfig struct {
    Window time.Duration
    Limit  int
    Prefix string
}

func ProvideRateLimitConfigs(cfg *config.Config) map[string]RateLimiterConfig {
    return map[string]RateLimiterConfig{
        "register": {Window: time.Minute, Limit: cfg.RateLimit.Register, Prefix: "ratelimit:register:"},
        "login": {Window: time.Minute, Limit: cfg.RateLimit.Login, Prefix: "ratelimit:login:"},
        "refresh": {Window: time.Minute, Limit: cfg.RateLimit.Refresh, Prefix: "ratelimit:refresh:"},
        "logout": {Window: time.Minute, Limit: cfg.RateLimit.Logout, Prefix: "ratelimit:logout:"},
    }
}

// Allow проверяет, разрешён ли запрос для данного ключа и endpoint
func (rl *RateLimiter) Allow(ctx context.Context, endpoint, key string) (bool, int, time.Time, error) {
    config, ok := rl.configs[endpoint]
    if !ok { /* ... */ }
    // ... сложная логика с Redis pipeline
}
```

**Почему оверинженеринг:**
- ❌ Map для 4 endpoint'ов
- ❌ Отдельный конфиг для каждого
- ❌ Сложная валидация `Validate()`
- ❌ Возврат 4 значений (bool, int, time.Time, error)

**Решение:**
```go
type RateLimiter struct {
    client *redis.Client
}

func (rl *RateLimiter) AllowRegister(ctx context.Context, key string) (bool, error) {
    return rl.allow(ctx, "ratelimit:register:", key, 5, time.Minute)
}

func (rl *RateLimiter) AllowLogin(ctx context.Context, key string) (bool, error) {
    return rl.allow(ctx, "ratelimit:login:", key, 10, time.Minute)
}

func (rl *RateLimiter) AllowRefresh(ctx context.Context, key string) (bool, error) {
    return rl.allow(ctx, "ratelimit:refresh:", key, 30, time.Minute)
}

func (rl *RateLimiter) AllowLogout(ctx context.Context, key string) (bool, error) {
    return rl.allow(ctx, "ratelimit:logout:", key, 60, time.Minute)
}

// Внутренний метод
func (rl *RateLimiter) allow(ctx context.Context, prefix, key string, limit int, window time.Duration) (bool, error) {
    // ... упрощённая логика
}
```

**Файлы для изменения:**
- `internal/middleware/rate_limiter.go`
- `internal/middleware/rate_limiter_http.go` (обновить вызовы)
- `internal/di/provider.go` (упростить ProvideRateLimitConfigs)

---

### 6. **JWT Service: Дублирование Validate** 🟡 Средний приоритет

**Файл:** `pkg/jwt/jwt.go`

**Проблема:**
```go
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
    // ... парсинг и валидация токена
}

func (s *Service) ValidateRefreshToken(tokenString string) (*Claims, error) {
    c, err := s.ValidateToken(tokenString)  // ← Дублирование
    if err != nil { return nil, err }
    if c.Type != refreshToken { return nil, ErrInvalidClaims }
    return c, nil
}
```

**Почему оверинженеринг:**
- ❌ `ValidateRefreshToken` просто обёртка
- ❌ Проверка типа токена могла быть в `ValidateToken`

**Решение:**
```go
func (s *Service) ValidateToken(tokenString string, requiredType tokenType) (*Claims, error) {
    // ... парсинг
    if c.Type != requiredType {
        return nil, ErrInvalidClaims
    }
    return c, nil
}

// Использование:
claims, err := jwt.ValidateToken(token, refreshToken)
```

**Файлы для изменения:**
- `pkg/jwt/jwt.go`
- `internal/service/auth/*.go` (обновить вызовы)

---

### 7. **Config: Избыточные методы Duration** 🟢 Низкий приоритет

**Файл:** `internal/config/config.go`

**Проблема:**
```go
func (c *JWTConfig) RefreshTTLDuration() (time.Duration, error) {
    return time.ParseDuration(c.RefreshTTL)
}

func (c *JWTConfig) AccessTTLDuration() (time.Duration, error) {
    return time.ParseDuration(c.AccessTTL)
}

func (c *ServerConfig) ReadTimeoutDuration() time.Duration {
    return time.Duration(c.ReadTimeout) * time.Second
}

func (c *ServerConfig) WriteTimeoutDuration() time.Duration {
    return time.Duration(c.WriteTimeout) * time.Second
}

func (c *ServerConfig) IdleTimeoutDuration() time.Duration {
    return time.Duration(c.IdleTimeout) * time.Second
}

func (c *ShutdownConfig) TimeoutDuration() time.Duration {
    return time.Duration(c.Timeout) * time.Second
}
```

**Почему оверинженеринг:**
- ❌ 6 методов для простого парсинга
- ❌ Вызываются 1 раз при старте
- ❌ Можно парсить напрямую в DI

**Решение:**
```go
// Парсить напрямую в DI:
accessTTL, _ := time.ParseDuration(cfg.JWT.AccessTTL)
refreshTTL, _ := time.ParseDuration(cfg.JWT.RefreshTTL)

return jwt.NewService(jwt.Config{
    SecretKey:       cfg.JWT.Secret,
    AccessTokenTTL:  accessTTL,
    RefreshTokenTTL: refreshTTL,
    Issuer:          cfg.JWT.Issuer,
})
```

**Файлы для изменения:**
- `internal/config/config.go` (удалить методы)
- `internal/di/provider.go` (парсить напрямую)

---

### 8. **Handler Converter: 1 функция** 🟢 Низкий приоритет

**Файл:** `internal/handler/converter/account.go`

**Проблема:**
```go
// internal/handler/converter/account.go
func AccountToProto(account *model.Account) *authv1.RegisterData {
    return &authv1.RegisterData{
        AccountId: account.ID().String(),
        Email:     account.Email().String(),
        CreatedAt: timestamppb.New(account.CreatedAt()),
    }
}
```

**Почему оверинженеринг:**
- ❌ Отдельный файл для 1 функции
- ❌ Можно inline в handler

**Решение:**
```go
// internal/handler/auth/register.go
return &authv1.RegisterResponse{
    StatusCode: 200,
    Message:    "Account registered successfully",
    Data: &authv1.RegisterData{
        AccountId: account.ID().String(),
        Email:     account.Email().String(),
        CreatedAt: timestamppb.New(account.CreatedAt()),
    },
}
```

**Файлы для изменения:**
- `internal/handler/converter/account.go` (удалить)
- `internal/handler/auth/register.go` (inline)

---

## 📋 План рефакторинга

### 🔴 Критичные (улучшить архитектуру)

| # | Задача | Файлы | Строк | Приоритет |
|---|--------|-------|-------|-----------|
| 1 | Упростить Application struct | `internal/di/provider.go` | ~50 | Высокий |
| 2 | Заменить Setters на конструктор | `internal/model/account.go` | ~20 | Высокий |
| 3 | Упростить Converter | `internal/repository/converter/account.go` | ~30 | Высокий |

### 🟡 Средние (уменьшить код)

| # | Задача | Файлы | Строк | Приоритет |
|---|--------|-------|-------|-----------|
| 4 | PasswordHash — value type | `internal/model/password_hash.go` | ~15 | Средний |
| 5 | RateLimiter — отдельные методы | `internal/middleware/rate_limiter.go` | ~80 | Средний |
| 6 | JWT Validate — один метод | `pkg/jwt/jwt.go` | ~20 | Средний |

### 🟢 Низкие (очистка)

| # | Задача | Файлы | Строк | Приоритет |
|---|--------|-------|-------|-----------|
| 7 | Config Duration — удалить методы | `internal/config/config.go` | ~30 | Низкий |
| 8 | Handler Converter — inline | `internal/handler/converter/account.go` | ~10 | Низкий |

---

## ✅ Чек-лист выполнения

- [ ] 1. Упростить Application struct (DI)
- [ ] 2. Заменить Setters на NewAccountFromDB
- [ ] 3. Упростить Converter (использовать NewAccountFromDB)
- [ ] 4. PasswordHash — value type (удалить pointer)
- [ ] 5. RateLimiter — отдельные методы для endpoint'ов
- [ ] 6. JWT Validate — один метод с requiredType
- [ ] 7. Config Duration — удалить методы
- [ ] 8. Handler Converter — inline в handler

---

## 📚 Ссылки

- [DDD Architecture](README.md)
- [Development Guide](DEVELOPMENT.md)
- [API Documentation](api.md)
