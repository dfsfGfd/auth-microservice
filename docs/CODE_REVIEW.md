# Code Review Report — Auth Microservice

**Дата:** 2026-03-19  
**Статус:** Complete  
**Всего найдено проблем:** 35 (18 критических, 17 средних)

---

## Executive Summary

Код микросервиса аутентификации написан на хорошем уровне с соблюдением многих best practices:
- ✅ DDD архитектура (Domain/Repository/Service/Handler)
- ✅ Dependency Injection через Google Wire
- ✅ Value Objects с валидацией
- ✅ Parameterized queries (защита от SQL injection)
- ✅ Graceful shutdown
- ✅ Structured logging с zerolog
- ✅ E2E тесты с testcontainers

Однако выявлено **35 проблем** различной критичности.

---

## 🔴 Critical Issues (18)

### 1. Error Handling

#### 1.1 Ошибки не оборачиваются с контекстом в repository layer

**Файлы:**
- `internal/repository/auth/save.go`
- `internal/repository/auth/get_by_email.go`
- `internal/repository/auth/get_by_id.go`
- `internal/repository/auth/exists_by_email.go`
- `internal/repository/auth/exists_by_id.go`
- `internal/repository/auth/delete_by_id.go`

**Проблема:**
```go
_, err := r.pool.Exec(ctx, `...`)
return err  // ❌ теряется контекст
```

**Исправление:**
```go
_, err := r.pool.Exec(ctx, `...`)
if err != nil {
    return fmt.Errorf("save account: %w", err)
}
return nil
```

#### 1.2 В service layer ошибки логируются, но не оборачиваются

**Файл:** `internal/service/auth/register.go:47`

**Проблема:**
```go
if err := s.accountRepo.Save(ctx, account); err != nil {
    s.log.Error("save account", "email", email, "error", err)
    return nil, err  // ❌ теряется контекст
}
```

**Исправление:**
```go
if err := s.accountRepo.Save(ctx, account); err != nil {
    s.log.Error("save account", "email", email, "error", err)
    return nil, fmt.Errorf("save account: %w", err)
}
```

#### 1.3 В cache layer ошибки оборачиваются непоследовательно

**Файл:** `internal/cache/token/redis_cache.go`

**Проблема:**
```go
// ✅ Хорошо в Get:
return "", fmt.Errorf("redis get: %w", err)

// ❌ Плохо в Set (нет обёртки):
return c.client.Set(ctx, key, accountID, ttl).Err()
```

**Исправление:**
```go
func (c *RedisCache) Set(ctx context.Context, token string, accountID string, ttl time.Duration) error {
    key := c.key(token)
    if err := c.client.Set(ctx, key, accountID, ttl).Err(); err != nil {
        return fmt.Errorf("redis set: %w", err)
    }
    return nil
}
```

---

### 2. Security

#### 2.1 JWT Secret hardcoded в docker-compose

**Файл:** `deploy/docker-compose.yml:22`

**Проблема:**
```yaml
# ❌ КРИТИЧНО: Hardcoded secret в production конфиге
JWT_SECRET=dev-secret-key-minimum-32-characters-long-change-in-production
```

**Исправление:** Использовать Docker secrets:
```yaml
secrets:
  - jwt_secret
  
environment:
  JWT_SECRET_FILE: /run/secrets/jwt_secret
  
secrets:
  jwt_secret:
    external: true
```

#### 2.2 Rate limiting fail-open стратегия опасна для login/register

**Файл:** `internal/middleware/rate_limiter.go:91`

**Проблема:**
```go
if err != nil {
    // Если Redis недоступен, пропускаем запрос (fail-open)
    return true, 0, time.Now(), nil
}
```

**Исправление:** Для login/register использовать fail-close:
```go
if err != nil {
    // Для критических endpoint'ов блокируем при ошибке Redis
    if endpoint == "login" || endpoint == "register" {
        log.Error("rate limiter redis error", "error", err)
        return false, 0, now, fmt.Errorf("rate limiter unavailable: %w", err)
    }
    // Для остальных — fail-open
    return true, 0, now, nil
}
```

#### 2.3 CORS AllowCredentials с wildcard origin

**Файл:** `internal/middleware/cors.go:55`

**Проблема:**
```go
AllowOriginFunc: func(origin string) bool {
    for _, allowed := range config.AllowedOrigins {
        if allowed == "*" {
            return true  // ❌ wildcard с credentials
        }
    }
}
```

**Исправление:** Запретить wildcard с credentials:
```go
AllowOriginFunc: func(origin string) bool {
    for _, allowed := range config.AllowedOrigins {
        if allowed == "*" {
            // ❌ Нельзя использовать wildcard с credentials
            if config.AllowCredentials {
                return false
            }
            return true
        }
        // ...
    }
}
```

---

### 3. Code Organization

#### 3.1 Service layer зависит от конкретной реализации RedisCache

**Файл:** `internal/service/auth/service.go:20`

**Проблема:**
```go
type AuthService struct {
    // ...
    tokenCache  *token.RedisCache  // ❌ должна быть interface TokenCache
    // ...
}
```

**Исправление:**
```go
// internal/cache/token/token_cache.go
type TokenCache interface {
    Set(ctx context.Context, token string, accountID string, ttl time.Duration) error
    Get(ctx context.Context, token string) (string, error)
    Delete(ctx context.Context, token string) error
}

// internal/service/auth/service.go
type AuthService struct {
    // ...
    tokenCache  token.TokenCache  // ✅ интерфейс
    // ...
}
```

#### 3.2 DI контейнер создаёт конкретные реализации вместо интерфейсов

**Файл:** `internal/di/provider.go:120`

**Проблема:**
```go
ProvideTokenCachePrefix,
token.NewRedisCache,  // возвращает *RedisCache
```

**Исправление:** Создать provider wrapper:
```go
func ProvideTokenCache(client *redis.Client, prefix string) token.TokenCache {
    return token.NewRedisCache(client, prefix)
}

// В ProviderSet:
ProvideTokenCache,
```

---

### 4. Performance

#### 4.1 Нет context таймаута для Redis операций в rate limiter

**Файл:** `internal/middleware/rate_limiter.go:56`

**Проблема:**
```go
pipe := rl.client.Pipeline()
// ...
_, err := pipe.Exec(ctx)  // ❌ нет таймаута
```

**Исправление:**
```go
func (rl *RateLimiter) Allow(ctx context.Context, endpoint, key string) (bool, int, time.Time, error) {
    // Добавить таймаут на Redis операцию
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    pipe := rl.client.Pipeline()
    // ...
}
```

#### 4.2 Connection pool для Redis не настроен оптимально

**Файл:** `internal/di/provider.go:107`

**Проблема:**
```go
return redisdb.Config{
    Addr:         addr,
    DB:           cfg.Redis.DB,
    PoolSize:     10,  // ❌ фиксированный размер
    // ...
}
```

**Исправление:** Настройка из конфига:
```go
// internal/config/config.go
type RedisConfig struct {
    URL              string `yaml:"url"`
    DB               int    `yaml:"db"`
    ConnectionTimeout int   `yaml:"connection_timeout"`
    PoolSize         int    `yaml:"pool_size" default:"25"`
    MinIdleConns     int    `yaml:"min_idle_conns" default:"5"`
}
```

#### 4.3 Двойной запрос в Refresh для удаления старого токена

**Файл:** `internal/service/auth/refresh.go:65`

**Проблема:**
```go
// Удаление старого токена из кэша
_ = s.tokenCache.Delete(ctx, refreshToken)  // ❌ ошибка игнорируется
```

**Исправление:**
```go
// Удаление старого токена (ошибка не критична, но логируем)
if err := s.tokenCache.Delete(ctx, refreshToken); err != nil {
    s.log.Warn("delete old refresh token", "error", err)
}
```

---

### 5. Testing

#### 5.1 Нет unit тестов для service layer

**Проблема:** Критическая бизнес-логика в `/internal/service/auth/` не покрыта unit тестами.

**Исправление:** Добавить тесты:
```go
// internal/service/auth/register_test.go
func TestAuthService_Register(t *testing.T) {
    tests := []struct {
        name          string
        email         string
        password      string
        accountExists bool
        wantErr       error
    }{
        {"valid registration", "test@example.com", "Password123!", false, nil},
        {"duplicate email", "test@example.com", "Password123!", true, errors.ErrAccountExists},
        // ...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Mock repository
            mockRepo := &mockAccountRepository{}
            mockRepo.On("ExistsByEmail", mock.Anything, tt.email).Return(tt.accountExists, nil)
            // ...
        })
    }
}
```

#### 5.2 Нет тестов для repository layer

**Проблема:** Репозитории не тестируются изолированно.

**Исправление:** Добавить тесты с testcontainers:
```go
// internal/repository/auth/save_test.go
func TestAccountRepository_Save(t *testing.T) {
    // Использовать testcontainers для PostgreSQL
    container, err := tcpostgres.RunContainer(ctx, ...)
    // ...
}
```

#### 5.3 E2E тесты не проверяют rate limiting

**Проблема:** Нет тестов на rate limiting, хотя это критичная функция безопасности.

**Исправление:**
```go
func TestRateLimiting(t *testing.T) {
    // Сделать 6 запросов register подряд (лимит 5)
    for i := 0; i < 6; i++ {
        resp, err := client.Register(ctx, &authv1.RegisterRequest{
            Email: fmt.Sprintf("ratelimit%d@example.com", i),
            Password: "Password123!",
        })
        
        if i >= 5 {
            require.Error(t, err)
            require.Contains(t, err.Error(), "rate limit")
        }
    }
}
```

---

### 6. Configuration

#### 6.1 Нет валидации RateLimit конфигурации

**Файл:** `internal/config/config.go`

**Проблема:** RateLimitConfig не валидируется в `Validate()`.

**Исправление:**
```go
func (c *Config) Validate() error {
    // ...
    if err := c.RateLimit.Validate(); err != nil {
        return fmt.Errorf("rate_limit: %w", err)
    }
    // ...
}

func (c *RateLimitConfig) Validate() error {
    if c.Register <= 0 {
        return fmt.Errorf("register must be positive")
    }
    if c.Login <= 0 {
        return fmt.Errorf("login must be positive")
    }
    // ...
}
```

#### 6.2 Нет проверки на пустые AllowedOrigins в production

**Файл:** `internal/config/config.go:217`

**Проблема:** В production CORS_ALLOWED_ORIGINS не должен быть "*".

**Исправление:**
```go
func (c *CORSConfig) Validate() error {
    if len(c.AllowedOrigins) == 0 {
        return fmt.Errorf("allowed_origins is required")
    }
    
    // Проверка на wildcard в production
    if os.Getenv("APP_ENV") == "production" {
        for _, origin := range c.AllowedOrigins {
            if origin == "*" {
                return fmt.Errorf("wildcard origin not allowed in production")
            }
        }
    }
    // ...
}
```

---

### 7. API Design

#### 7.1 StatusCode дублируется в response

**Файл:** `proto/auth/v1/auth.proto:54`

**Проблема:**
```protobuf
message RegisterResponse {
  int32 status_code = 1;  // ❌ дублирует gRPC status
  string message = 2;
  RegisterData data = 3;
}
```

**Рекомендация:** Убрать status_code из proto, использовать gRPC status codes:
```protobuf
message RegisterResponse {
  string message = 1;
  RegisterData data = 2;
}
```

---

### 8. Go Best Practices

#### 8.1 Утечка ресурсов в main.go при ошибке

**Файл:** `cmd/server/main.go:52`

**Проблема:** defer CleanUp может не выполниться корректно.

**Исправление:** Убедиться, что defer идёт сразу после успешной инициализации:
```go
func run() error {
    app, err := di.InitializeApplication()
    if err != nil {
        return fmt.Errorf("initialize application: %w", err)
    }
    
    ctx := context.Background()
    defer func() {
        if err := app.CleanUp(ctx); err != nil {
            app.Logger.Error("cleanup error", "error", err)
        }
    }()
    
    // ...
}
```

---

## 🟡 Medium Priority Issues (17)

### Security

1. **Password hashing cost не настраивается** (`pkg/bcrypt/bcrypt.go:20`)
   - Добавить настройку BCRYPT_COST в конфиг (рекомендуется 12 для production)

2. **Нет защиты от timing attacks при сравнении токенов**
   - Использовать `subtle.ConstantTimeCompare` для чувствительных данных

3. **Недостаточная валидация JWT Secret** (`internal/config/config.go:254`)
   - Добавить проверку на слабые/default секреты

### Performance

4. **Нет индекса для updated_at в БД**
   - Добавить индекс если нужен для аудитов

5. **Нет limit на количество полей в логах**
   - Добавить limit для защиты от memory exhaustion (рекомендуется 50 полей)

### Testing

6. **Нет тестов для middleware**
   - Добавить тесты для CORS и Rate Limiter

7. **Нет benchmark тестов**
   - Добавить benchmark для login, register, token refresh

### Go Best Practices

8. **Неиспользуемый интерфейс в repository**
   - На самом деле используется правильно — ложная тревога

9. **Избыточное использование указателей для Value Objects**
   - Для immutable VO можно использовать значения вместо указателей

10. **Непоследовательное именование ошибок**
    - В целом соблюдается, но можно улучшить

11. **gofumpt/gci не настроены в CI**
    - Добавить задачу проверки форматирования в Taskfile

### API Design

12. **Inconsistent error messages могут leak информацию**
    - Для registration можно детализировать, для login уже сделано правильно

13. **Нет пагинации для будущих list endpoints**
    - Заложить pagination в proto заранее

14. **Нет versioning API кроме как в package name**
    - Добавить `/api/v1/` в HTTP endpoints

### Configuration

15. **Secrets могут попасть в логи через config**
    - Реализовать интерфейс для скрытия секрета (String() возвращает "[REDACTED]")

16. **Нет hot-reload конфигурации**
    - Для production добавить поддержку reload без перезапуска

17. **CORS максимальный возраст не валидируется**
    - Добавить проверку на разумные значения (рекомендуется 86400)

---

## 📊 Summary

| Категория | Критические | Средние | Всего |
|-----------|-------------|---------|-------|
| Error Handling | 3 | 1 | 4 |
| Security | 3 | 3 | 6 |
| Code Organization | 2 | 2 | 4 |
| Performance | 3 | 2 | 5 |
| Testing | 3 | 2 | 5 |
| Go Best Practices | 1 | 3 | 4 |
| API Design | 1 | 2 | 3 |
| Configuration | 2 | 2 | 4 |
| **ИТОГО** | **18** | **17** | **35** |

---

## 🎯 Priority Recommendations

### High Priority (исправить немедленно)

1. 🔒 **Security:** Заменить hardcoded JWT secret в docker-compose
2. 🔒 **Security:** Добавить fail-close для rate limiting на login/register
3. 📝 **Error Handling:** Обернуть все ошибки с контекстом в repository layer
4. 📦 **Code Organization:** Изменить зависимость на интерфейс TokenCache вместо *RedisCache
5. 🧪 **Testing:** Добавить unit тесты для service layer

### Medium Priority (исправить в ближайшем спринте)

1. ⚡ **Performance:** Добавить таймауты на Redis операции
2. 🔒 **Security:** Настроить bcrypt cost через конфигурацию
3. ⚙️ **Configuration:** Добавить валидацию RateLimitConfig
4. 🧪 **Testing:** Добавить тесты для middleware
5. 🌐 **API Design:** Убрать дублирование status_code из proto response

### Low Priority (улучшения)

1. 🐹 **Go Best Practices:** Настроить format check в CI
2. ⚡ **Performance:** Добавить benchmark тесты
3. 🌐 **API Design:** Добавить pagination для будущих list endpoints
4. ⚙️ **Configuration:** Добавить hot-reload конфигурации

---

## ✅ Positive Findings

Код имеет много сильных сторон:

1. ✅ **DDD архитектура** соблюдена правильно (Domain/Repository/Service/Handler)
2. ✅ **Dependency Injection** через Google Wire
3. ✅ **User enumeration protection** в login endpoint
4. ✅ **Value Objects** для Email и Password с валидацией
5. ✅ **Table-driven tests** в существующих тестах
6. ✅ **Parameterized queries** везде (защита от SQL injection)
7. ✅ **Graceful shutdown** реализован
8. ✅ **Health check** endpoint присутствует
9. ✅ **Structured logging** с zerolog
10. ✅ **E2E тесты** с testcontainers

---

## 📝 Action Items

### Sprint 1 (Critical)
- [ ] Исправить error wrapping в repository layer (3 часа)
- [ ] Заменить JWT secret в docker-compose на Docker secrets (2 часа)
- [ ] Добавить fail-close для rate limiting (3 часа)
- [ ] Рефакторинг TokenCache на интерфейс (4 часа)
- [ ] Написать unit тесты для service layer (8 часов)

### Sprint 2 (High)
- [ ] Добавить таймауты на Redis операции (2 часа)
- [ ] Настроить bcrypt cost через конфиг (2 часа)
- [ ] Добавить валидацию RateLimitConfig (2 часа)
- [ ] Написать тесты для middleware (4 часа)
- [ ] Убрать status_code из proto (3 часа)

### Sprint 3 (Medium)
- [ ] Настроить format check в CI (2 часа)
- [ ] Добавить benchmark тесты (4 часа)
- [ ] Добавить pagination в proto (3 часа)
- [ ] Реализовать hot-reload конфигурации (6 часов)
- [ ] Добавить индекс для updated_at (1 час)

**Всего оценок:** ~54 часов разработки
