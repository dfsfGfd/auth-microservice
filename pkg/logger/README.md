# Logger Package

Структурированное логирование на основе zerolog.

---

## 🚀 Использование

```go
import "auth-microservice/pkg/logger"

// Создание
log, err := logger.New(logger.Config{
    Level:       "info",
    Format:      "json",
    ServiceName: "auth-service",
})

// Логирование
log.Info("user logged in", "user_id", userID, "email", email)
log.Error("database error", "error", err)
log.Warn("rate limit exceeded", "ip", clientIP)
log.Debug("request received", "method", r.Method)

// Контекстный логгер
reqLog := log.With("request_id", requestID)
reqLog.Info("processing request")
```

---

## API

### Config

```go
type Config struct {
    Level       string      // debug, info, warn, error, fatal
    Format      string      // json, console
    Output      io.Writer   // os.Stderr (по умолчанию)
    ServiceName string      // название сервиса
}
```

### Методы

```go
// Создание
func New(config Config) (*Logger, error)

// Уровни логирования
func (l *Logger) Debug(msg string, keysAndValues ...interface{})
func (l *Logger) Info(msg string, keysAndValues ...interface{})
func (l *Logger) Warn(msg string, keysAndValues ...interface{})
func (l *Logger) Error(msg string, keysAndValues ...interface{})
func (l *Logger) Fatal(msg string, keysAndValues ...interface{})

// Контекст
func (l *Logger) With(keysAndValues ...interface{}) *Logger
```

---

## Уровни

| Уровень | Когда |
|---------|-------|
| `debug` | Отладка (SQL, детали) |
| `info` | Нормальная работа |
| `warn` | Предупреждения |
| `error` | Ошибки |
| `fatal` | Критические (с exit) |

---

## Примеры

**JSON (production):**
```json
{"level":"info","service":"auth-service","time":"2024-01-15T10:30:00Z","user_id":"123","message":"user logged in"}
```

**Console (development):**
```
10:30:00 INF user logged in user_id=123
```

---

## 🧪 Тесты

```bash
go test ./pkg/logger/... -v
```
