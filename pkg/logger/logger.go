// Package logger предоставляет утилиты для структурированного логирования.
//
// Пример использования:
//
//	// Создание логгера
//	log := logger.New(logger.Config{
//	    Level:  "info",
//	    Format: "json",
//	})
//
//	// Логирование
//	log.Info("user logged in", "user_id", userID)
//	log.Error("database error", "error", err)
//
// Формат JSON лога (оптимизированный):
// {"ts":"2026-03-21T15:00:00Z","lvl":"info","msg":"user logged in","srv":"auth-service","user_id":"550e8400","dur_ms":45}
package logger

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Format определяет формат вывода логов
type Format string

const (
	// JSONFormat JSON формат (для продакшена)
	JSONFormat Format = "json"
	// ConsoleFormat Человекочитаемый формат (для разработки)
	ConsoleFormat Format = "console"
)

// Level определяет уровень логирования
type Level string

const (
	// DebugLevel отладочные сообщения
	DebugLevel Level = "debug"
	// InfoLevel информационные сообщения
	InfoLevel Level = "info"
	// WarnLevel предупреждения
	WarnLevel Level = "warn"
	// ErrorLevel ошибки
	ErrorLevel Level = "error"
	// FatalLevel фатальные ошибки (с завершением программы)
	FatalLevel Level = "fatal"
)

// contextKey тип для ключей контекста
type contextKey string

const (
	// RequestIDKey ключ для request_id в контексте
	RequestIDKey contextKey = "request_id"
	// TraceIDKey ключ для trace_id в контексте
	TraceIDKey contextKey = "trace_id"
)

// Config конфигурация логгера
type Config struct {
	// Level уровень логирования (debug, info, warn, error)
	Level string
	// Format формат вывода (json, console)
	Format string
	// Output поток вывода (по умолчанию os.Stderr)
	Output io.Writer
	// ServiceName название сервиса для добавления в логи
	ServiceName string
}

// Validate валидирует конфигурацию
func (c *Config) Validate() error {
	if c.Level == "" {
		c.Level = "info"
	}
	if c.Format == "" {
		c.Format = "json"
	}
	if c.Output == nil {
		c.Output = os.Stderr
	}
	return nil
}

// Logger структурированный логгер
type Logger struct {
	zlog zerolog.Logger
}

// New создаёт новый логгер с оптимизированными именами полей
func New(config Config) (*Logger, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Парсим уровень логирования
	level, err := parseLevel(config.Level)
	if err != nil {
		return nil, err
	}

	// Настраиваем вывод
	var output io.Writer = config.Output

	// Для console формата используем PrettyWriter
	if strings.ToLower(config.Format) == string(ConsoleFormat) {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: "2006-01-02 15:04:05",
		}
	}

	// Создаём zerolog логгер с оптимизированными именами полей
	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"

	zlog := zerolog.New(output).
		Level(level).
		With().
		Timestamp()

	// Добавляем название сервиса если указано
	if config.ServiceName != "" {
		zlog = zlog.Str("srv", config.ServiceName)
	}

	return &Logger{
		zlog: zlog.Logger(),
	}, nil
}

// parseLevel парсит строковый уровень в zerolog.Level
func parseLevel(level string) (zerolog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel, nil
	case "info":
		return zerolog.InfoLevel, nil
	case "warn", "warning":
		return zerolog.WarnLevel, nil
	case "error":
		return zerolog.ErrorLevel, nil
	case "fatal":
		return zerolog.FatalLevel, nil
	default:
		return zerolog.InfoLevel, nil
	}
}

// Debug логирует отладочное сообщение
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.zlog.Debug(), msg, keysAndValues)
}

// Info логирует информационное сообщение
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.zlog.Info(), msg, keysAndValues)
}

// Warn логирует предупреждение
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.zlog.Warn(), msg, keysAndValues)
}

// Error логирует ошибку
func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.zlog.Error(), msg, keysAndValues)
}

// Fatal логирует фатальную ошибку и завершает программу
func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.zlog.Fatal(), msg, keysAndValues)
}

// logEvent записывает событие с полями (оптимизированные имена)
func (l *Logger) logEvent(event *zerolog.Event, msg string, keysAndValues []interface{}) {
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			if key, ok := keysAndValues[i].(string); ok {
				// Специальная обработка ошибок
				if key == "error" && keysAndValues[i+1] != nil {
					if err, ok := keysAndValues[i+1].(error); ok {
						event.Err(err)
						continue
					}
				}
				// Оптимизированные имена полей
				optKey := optimizeFieldName(key)
				event.Any(optKey, keysAndValues[i+1])
			}
		}
	}
	event.Msg(msg)
}

// optimizeFieldName сокращает имена полей для компактности
func optimizeFieldName(key string) string {
	switch key {
	case "user_id":
		return "user_id" // Оставляем user_id для читаемости
	case "request_id":
		return "rid"
	case "trace_id":
		return "trace"
	case "span_id":
		return "span"
	case "duration_ms":
		return "dur_ms"
	case "email":
		return "email"
	case "method":
		return "method"
	case "path":
		return "path"
	case "status":
		return "status"
	case "error":
		return "err"
	default:
		return key
	}
}

// With возвращает новый логгер с дополнительными полями
func (l *Logger) With(keysAndValues ...interface{}) *Logger {
	// Создаём контекст с полями
	ctx := l.zlog.With()

	// Добавляем каждое поле в контекст
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			if key, ok := keysAndValues[i].(string); ok {
				ctx = ctx.Any(key, keysAndValues[i+1])
			}
		}
	}

	newLogger := ctx.Logger()
	return &Logger{
		zlog: newLogger,
	}
}

// WithContext возвращает логгер с полями из контекста
func (l *Logger) WithContext(ctx context.Context) *Logger {
	event := l.zlog.With()

	// Добавляем trace_id если есть (оптимизированное имя)
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		event = event.Str("trace", traceID.(string))
	}

	// Добавляем request_id если есть (оптимизированное имя)
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		event = event.Str("rid", requestID.(string))
	}

	return &Logger{
		zlog: event.Logger(),
	}
}

// WithRequestID возвращает логгер с request_id
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		zlog: l.zlog.With().Str("rid", requestID).Logger(),
	}
}

// WithError возвращает логгер с ошибкой
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		zlog: l.zlog.With().Err(err).Logger(),
	}
}

// SetGlobalLevel устанавливает глобальный уровень логирования
func SetGlobalLevel(level string) error {
	zLevel, err := parseLevel(level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(zLevel)
	return nil
}

// DisableLogging отключает логирование
func DisableLogging() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

// GetLevel возвращает текущий уровень логирования
func (l *Logger) GetLevel() string {
	return l.zlog.GetLevel().String()
}

// WithContext добавляет request_id в контекст
func WithContext(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// FromContext получает request_id из контекста
func FromContext(ctx context.Context) string {
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		return requestID.(string)
	}
	return ""
}
