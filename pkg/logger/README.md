# Logger Package

Пакет для структурированного логирования. Может быть экспортирован и использован в других проектах.

## 📦 Установка

Пакет использует [zerolog](https://github.com/rs/zerolog) под капотом:

```bash
go get github.com/rs/zerolog
```

## 🚀 Быстрый старт

### Создание логгера

```go
import "auth-microservice/pkg/logger"

// Создание логгера с конфигурацией
log, err := logger.New(logger.Config{
    Level:       "info",
    Format:      "json",
    ServiceName: "auth-service",
})
if err != nil {
    panic(err)
}
```

### Базовое логирование

```go
// Логирование с полями (key-value)
log.Info("user logged in", "user_id", userID, "email", email)
log.Error("database error", "error", err, "query", query)
log.Warn("rate limit exceeded", "ip", clientIP, "attempts", attempts)
log.Debug("request received", "method", r.Method, "path", r.URL.Path)
```

### Контекстное логирование

```go
// Добавление контекста к логгеру
contextLog := log.With("request_id", requestID, "trace_id", traceID)

// Все логи этого контекста будут содержать указанные поля
contextLog.Info("processing request")
contextLog.Error("processing failed", "error", err)
```

### Цепочка контекстов

```go
// Можно создавать вложенные контексты
baseLog := log.With("service", "auth")
requestLog := baseLog.With("request_id", requestID)
userLog := requestLog.With("user_id", userID)

userLog.Info("user action") 
// Выведет: {"service":"auth","request_id":"...","user_id":"...","message":"user action"}
```

## 📖 API

### Типы

#### `Config`

Конфигурация логгера.

```go
type Config struct {
    // Level уровень логирования (debug, info, warn, error, fatal)
    Level string
    
    // Format формат вывода (json, console)
    Format string
    
    // Output поток вывода (по умолчанию os.Stderr)
    Output io.Writer
    
    // ServiceName название сервиса для добавления в логи
    ServiceName string
}
```

**Поля:**

| Поле | Тип | По умолчанию | Описание |
|------|-----|--------------|----------|
| `Level` | `string` | `"info"` | Уровень логирования |
| `Format` | `string` | `"json"` | Формат вывода (json/console) |
| `Output` | `io.Writer` | `os.Stderr` | Поток вывода |
| `ServiceName` | `string` | `""` | Название сервиса |

#### `Format`

Формат вывода логов.

```go
type Format string

const (
    JSONFormat    Format = "json"    // JSON формат (для продакшена)
    ConsoleFormat Format = "console" // Человекочитаемый (для разработки)
)
```

#### `Level`

Уровень логирования.

```go
type Level string

const (
    DebugLevel Level = "debug" // Отладочные сообщения
    InfoLevel  Level = "info"  // Информационные сообщения
    WarnLevel  Level = "warn"  // Предупреждения
    ErrorLevel Level = "error" // Ошибки
    FatalLevel Level = "fatal" // Фатальные ошибки (с завершением программы)
)
```

#### `Logger`

Структурированный логгер.

```go
type Logger struct {
    // внутренние поля
}
```

### Функции

#### `New(config Config) (*Logger, error)`

Создаёт новый логгер с указанной конфигурацией.

```go
log, err := logger.New(logger.Config{
    Level:       "info",
    Format:      "json",
    ServiceName: "auth-service",
})
if err != nil {
    log.Fatal("failed to create logger", "error", err)
}
```

**Ошибки:**
- Невалидный уровень логирования (возвращает уровень по умолчанию `info`)

#### `SetGlobalLevel(level string) error`

Устанавливает глобальный уровень логирования для всех логгеров.

```go
err := logger.SetGlobalLevel("error")
if err != nil {
    log.Fatal("failed to set global level", "error", err)
}
```

#### `DisableLogging()`

Полностью отключает логирование.

```go
// Для тестов или когда логи не нужны
logger.DisableLogging()
```

### Методы логгера

#### `Debug(msg string, keysAndValues ...interface{})`

Логирует отладочное сообщение.

```go
log.Debug("cache miss", "key", cacheKey, "ttl", ttl)
```

#### `Info(msg string, keysAndValues ...interface{})`

Логирует информационное сообщение.

```go
log.Info("server started", "port", port, "host", host)
```

#### `Warn(msg string, keysAndValues ...interface{})`

Логирует предупреждение.

```go
log.Warn("deprecated endpoint called", "path", r.URL.Path, "client", clientIP)
```

#### `Error(msg string, keysAndValues ...interface{})`

Логирует ошибку.

```go
log.Error("database connection failed", "error", err, "retry_count", retryCount)
```

#### `Fatal(msg string, keysAndValues ...interface{})`

Логирует фатальную ошибку и завершает программу.

```go
if err != nil {
    log.Fatal("critical error", "error", err)
}
```

#### `With(keysAndValues ...interface{}) *Logger`

Возвращает новый логгер с дополнительными полями контекста.

```go
requestLog := log.With("request_id", requestID, "method", r.Method)
requestLog.Info("handling request")
```

#### `GetLevel() string`

Возвращает текущий уровень логирования.

```go
level := log.GetLevel()
fmt.Println(level) // "info"
```

## 🎯 Уровни логирования

| Уровень | Когда использовать |
|---------|-------------------|
| `debug` | Отладочная информация (значения переменных, SQL запросы, детали выполнения) |
| `info` | Нормальная работа системы (старт сервера, успешные операции) |
| `warn` | Предупреждения (устаревшие вызовы, временные проблемы, rate limiting) |
| `error` | Ошибки (сбои БД, валидация, внешние сервисы недоступны) |
| `fatal` | Критические ошибки (невозможно запустить сервер, конфигурация невалидна) |

## 📝 Примеры использования

### JSON формат (продакшен)

```go
log, _ := logger.New(logger.Config{
    Level:       "info",
    Format:      "json",
    ServiceName: "auth-service",
})

log.Info("user registered", "user_id", userID, "email", email)
```

**Вывод:**
```json
{"level":"info","service":"auth-service","time":"2024-01-15T10:30:00Z","user_id":"550e8400-e29b-41d4-a716-446655440000","email":"user@example.com","message":"user registered"}
```

### Console формат (разработка)

```go
log, _ := logger.New(logger.Config{
    Level:  "debug",
    Format: "console",
})

log.Info("server started", "port", 8080)
```

**Вывод:**
```
10:30:00 INF server started port=8080
```

### Middleware для HTTP запросов

```go
func LoggingMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Создаём контекстный логгер для запроса
            reqLog := log.With(
                "method", r.Method,
                "path", r.URL.Path,
                "request_id", uuid.New().String(),
            )
            
            reqLog.Debug("request received", "headers", r.Header)
            
            // Вызываем следующий хендлер
            next.ServeHTTP(w, r)
            
            // Логируем время выполнения
            reqLog.Info("request completed", "duration", time.Since(start))
        })
    }
}
```

### Логирование в сервисном слое

```go
type AuthService struct {
    log    *logger.Logger
    repo   repository.UserRepository
    jwtSvc *jwt.Service
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
    log := s.log.With("method", "Login", "email", email)
    
    log.Debug("finding user by email")
    user, err := s.repo.FindByEmail(ctx, email)
    if err != nil {
        log.Error("user not found", "error", err)
        return nil, ErrInvalidCredentials
    }
    
    log.Debug("verifying password")
    if err := bcrypt.CompareHashAndPassword(user.PasswordHash(), []byte(password)); err != nil {
        log.Error("password verification failed", "error", err)
        return nil, ErrInvalidCredentials
    }
    
    log.Info("login successful")
    return s.jwtSvc.GenerateTokens(user.ID().String(), user.Email().String(), user.Username().String())
}
```

### Логирование ошибок с контекстом

```go
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    r.log.Debug("creating user", 
        "email", user.Email().String(),
        "username", user.Username().String(),
    )
    
    query := `INSERT INTO users (id, email, username, password_hash, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
    
    _, err := r.db.Exec(ctx, query,
        user.ID(),
        user.Email().String(),
        user.Username().String(),
        user.PasswordHash().String(),
        user.CreatedAt(),
        user.UpdatedAt(),
    )
    
    if err != nil {
        r.log.Error("failed to create user",
            "error", err,
            "email", user.Email().String(),
            "query", query,
        )
        return ErrFailedToCreateUser
    }
    
    r.log.Info("user created successfully", "user_id", user.ID())
    return nil
}
```

### Логирование с recover от паник

```go
func RecoveryMiddleware(log *logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    log.Error("panic recovered",
                        "error", err,
                        "stack", string(debug.Stack()),
                        "method", r.Method,
                        "path", r.URL.Path,
                    )
                    http.Error(w, "internal server error", http.StatusInternalServerError)
                }
            }()
            next.ServeHTTP(w, r)
        })
    }
}
```

## 🔧 Конфигурация

### Переменные окружения

Рекомендуется использовать переменные окружения для конфигурации:

```go
config := logger.Config{
    Level:       os.Getenv("LOG_LEVEL"),       // "info", "debug", "warn", "error"
    Format:      os.Getenv("LOG_FORMAT"),      // "json" или "console"
    ServiceName: os.Getenv("SERVICE_NAME"),    // "auth-service"
}

log, err := logger.New(config)
```

### Рекомендации по настройке

| Окружение | Level | Format |
|-----------|-------|--------|
| **Development** | `debug` | `console` |
| **Testing** | `warn` или `error` | `console` |
| **Staging** | `info` | `json` |
| **Production** | `warn` или `error` | `json` |

## 🧪 Тесты

```bash
go test ./pkg/logger/... -v
```

### Пример тестирования с логгером

```go
func TestService_WithLogger(t *testing.T) {
    var buf bytes.Buffer
    
    // Создаём тестовый логгер
    log, _ := logger.New(logger.Config{
        Level:  "debug",
        Format: "json",
        Output: &buf,
    })
    
    // Используем в сервисе
    service := NewService(log, repo)
    err := service.DoSomething()
    
    // Проверяем что лог записался
    assert.Contains(t, buf.String(), "DoSomething")
    assert.NoError(t, err)
}
```

## 📦 Экспорт в другие проекты

Пакет может быть использован в других проектах:

```go
import "github.com/dfsfGfd/auth-microservice/pkg/logger"

// Использование в другом проекте
log, _ := logger.New(logger.Config{
    Level:       os.Getenv("LOG_LEVEL"),
    Format:      os.Getenv("LOG_FORMAT"),
    ServiceName: "my-service",
})

log.Info("service started")
```

## 🎯 Best Practices

### ✅ Делайте

- Используйте структурированные поля (key-value)
- Добавляйте контекст через `With()` для связанных операций
- Используйте соответствующие уровни логирования
- Логируйте ошибки с достаточным контекстом для отладки

### ❌ Не делайте

- Не логируйте чувствительные данные (пароли, токены, PII)
- Не злоупотребляйте `Debug` в продакшене
- Не используйте логирование вместо метрик
- Не логируйте в циклах без необходимости

## 📚 Ссылки

- [zerolog документация](https://github.com/rs/zerolog)
- [Основное README](../../README.md)
- [API документация](../../docs/api.md)
