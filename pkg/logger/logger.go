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
package logger

import (
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

// New создаёт новый логгер
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

	// Создаём zerolog логгер
	zlog := zerolog.New(output).
		Level(level).
		With().
		Timestamp()

	// Добавляем название сервиса если указано
	if config.ServiceName != "" {
		zlog = zlog.Str("service", config.ServiceName)
	}

	return &Logger{
		zlog: zlog.Caller().Logger(),
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

// logEvent записывает событие с полями
func (l *Logger) logEvent(event *zerolog.Event, msg string, keysAndValues []interface{}) {
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			if key, ok := keysAndValues[i].(string); ok {
				event.Any(key, keysAndValues[i+1])
			}
		}
	}
	event.Msg(msg)
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
