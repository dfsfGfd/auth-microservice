// Package logger предоставляет утилиты для структурированного логирования.
//
// Пример использования:
//
//	log := logger.New(logger.Config{
//	    Level:  "info",
//	    Format: "json",
//	})
//	log.Info("user logged in", "user_id", userID)
//	log.Error("database error", "error", err)
package logger

import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
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

// Logger структурированный логгер
type Logger struct {
	log zerolog.Logger
}

// New создаёт новый логгер
func New(cfg Config) (*Logger, error) {
	if cfg.Level == "" {
		cfg.Level = "info"
	}
	if cfg.Format == "" {
		cfg.Format = "json"
	}
	if cfg.Output == nil {
		cfg.Output = os.Stderr
	}

	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	var output io.Writer = cfg.Output

	if strings.ToLower(cfg.Format) == "console" {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: "2006-01-02 15:04:05",
		}
	}

	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "lvl"
	zerolog.MessageFieldName = "msg"

	log := zerolog.New(output).
		Level(level).
		With().
		Timestamp()

	if cfg.ServiceName != "" {
		log = log.Str("srv", cfg.ServiceName)
	}

	return &Logger{
		log: log.Logger(),
	}, nil
}

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
	default:
		return zerolog.InfoLevel, nil
	}
}

// Debug логирует отладочное сообщение
func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.log.Debug(), msg, keysAndValues)
}

// Info логирует информационное сообщение
func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.log.Info(), msg, keysAndValues)
}

// Warn логирует предупреждение
func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.log.Warn(), msg, keysAndValues)
}

// Error логирует ошибку
func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.logEvent(l.log.Error(), msg, keysAndValues)
}

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
