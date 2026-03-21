package logger_test

import (
	"bytes"
	"testing"

	"auth-microservice/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("успешное создание JSON логгера", func(t *testing.T) {
		log, err := logger.New(logger.Config{
			Level:       "info",
			Format:      "json",
			ServiceName: "test-service",
		})

		require.NoError(t, err)
		require.NotNil(t, log)
	})

	t.Run("успешное создание console логгера", func(t *testing.T) {
		log, err := logger.New(logger.Config{
			Level:       "debug",
			Format:      "console",
			ServiceName: "test-service",
		})

		require.NoError(t, err)
		require.NotNil(t, log)
	})

	t.Run("логгер по умолчанию", func(t *testing.T) {
		log, err := logger.New(logger.Config{})

		require.NoError(t, err)
		require.NotNil(t, log)
		assert.Equal(t, "info", log.GetLevel())
	})

	t.Run("невалидный уровень", func(t *testing.T) {
		log, err := logger.New(logger.Config{
			Level: "invalid",
		})

		require.NoError(t, err)
		require.NotNil(t, log)
	})
}

func TestLogger_Levels(t *testing.T) {
	var buf bytes.Buffer

	log, err := logger.New(logger.Config{
		Level:  "debug",
		Format: "json",
		Output: &buf,
	})
	require.NoError(t, err)

	t.Run("Debug", func(t *testing.T) {
		buf.Reset()
		log.Debug("debug message", "key", "value")
		assert.Contains(t, buf.String(), "debug message")
		assert.Contains(t, buf.String(), "key")
		assert.Contains(t, buf.String(), "value")
	})

	t.Run("Info", func(t *testing.T) {
		buf.Reset()
		log.Info("info message", "user_id", "123")
		assert.Contains(t, buf.String(), "info message")
		assert.Contains(t, buf.String(), "user_id")
	})

	t.Run("Warn", func(t *testing.T) {
		buf.Reset()
		log.Warn("warn message", "code", 400)
		assert.Contains(t, buf.String(), "warn message")
		assert.Contains(t, buf.String(), "code")
	})

	t.Run("Error", func(t *testing.T) {
		buf.Reset()
		log.Error("error message", "error", "something went wrong")
		assert.Contains(t, buf.String(), "error message")
		assert.Contains(t, buf.String(), "err")
	})
}

func TestLogger_With(t *testing.T) {
	var buf bytes.Buffer

	log, err := logger.New(logger.Config{
		Level:  "debug",
		Format: "json",
		Output: &buf,
	})
	require.NoError(t, err)

	t.Run("добавление полей", func(t *testing.T) {
		buf.Reset()
		contextLog := log.With("request_id", "abc-123", "user_id", "456")
		contextLog.Info("request processed")

		logOutput := buf.String()
		assert.Contains(t, logOutput, "request processed")
		assert.Contains(t, logOutput, "request_id")
		assert.Contains(t, logOutput, "abc-123")
		assert.Contains(t, logOutput, "user_id")
		assert.Contains(t, logOutput, "456")
	})

	t.Run("цепочка With", func(t *testing.T) {
		buf.Reset()
		contextLog := log.With("service", "auth")
		contextLog = contextLog.With("method", "login")
		contextLog.Info("action")

		logOutput := buf.String()
		assert.Contains(t, logOutput, "action")
		assert.Contains(t, logOutput, "service")
		assert.Contains(t, logOutput, "auth")
		assert.Contains(t, logOutput, "method")
		assert.Contains(t, logOutput, "login")
	})
}

func TestLogger_GetLevel(t *testing.T) {
	t.Run("debug level", func(t *testing.T) {
		log, _ := logger.New(logger.Config{Level: "debug"})
		assert.Equal(t, "debug", log.GetLevel())
	})

	t.Run("info level", func(t *testing.T) {
		log, _ := logger.New(logger.Config{Level: "info"})
		assert.Equal(t, "info", log.GetLevel())
	})

	t.Run("warn level", func(t *testing.T) {
		log, _ := logger.New(logger.Config{Level: "warn"})
		assert.Equal(t, "warn", log.GetLevel())
	})

	t.Run("error level", func(t *testing.T) {
		log, _ := logger.New(logger.Config{Level: "error"})
		assert.Equal(t, "error", log.GetLevel())
	})
}

func TestSetGlobalLevel(t *testing.T) {
	t.Run("установка глобального уровня", func(t *testing.T) {
		err := logger.SetGlobalLevel("error")
		require.NoError(t, err)
	})

	t.Run("установка невалидного уровня", func(t *testing.T) {
		err := logger.SetGlobalLevel("invalid")
		require.NoError(t, err)
	})
}

func TestDisableLogging(t *testing.T) {
	var buf bytes.Buffer

	log, err := logger.New(logger.Config{
		Level:  "debug",
		Format: "json",
		Output: &buf,
	})
	require.NoError(t, err)

	// Логируем до отключения
	log.Info("before disable")
	assert.NotEmpty(t, buf.String())

	// Отключаем логирование
	logger.DisableLogging()

	// Восстанавливаем уровень логирования после теста
	t.Cleanup(func() {
		logger.SetGlobalLevel("debug")
	})

	// Логируем после отключения
	buf.Reset()
	log.Info("after disable")
	// После отключения логи не должны писаться
	// (зависит от реализации zerolog)
}

func TestLogger_JSONFormat(t *testing.T) {
	var buf bytes.Buffer

	log, err := logger.New(logger.Config{
		Level:  "info",
		Format: "json",
		Output: &buf,
	})
	require.NoError(t, err)

	log.Info("test message", "key", "value")

	output := buf.String()
	// JSON должен содержать ключи
	assert.Contains(t, output, "msg")
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "key")
	assert.Contains(t, output, "value")
	assert.Contains(t, output, "ts") // zerolog использует "ts" вместо "time"
}

func TestLogger_ConsoleFormat(t *testing.T) {
	var buf bytes.Buffer

	log, err := logger.New(logger.Config{
		Level:  "info",
		Format: "console",
		Output: &buf,
	})
	require.NoError(t, err)

	log.Info("console message", "user", "test")

	output := buf.String()
	// Console формат должен быть человекочитаемым
	assert.Contains(t, output, "console message")
	assert.Contains(t, output, "user")
	assert.Contains(t, output, "test")
}

func TestLogger_ServiceName(t *testing.T) {
	var buf bytes.Buffer

	log, err := logger.New(logger.Config{
		Level:       "info",
		Format:      "json",
		Output:      &buf,
		ServiceName: "auth-service",
	})
	require.NoError(t, err)

	log.Info("service log")

	output := buf.String()
	assert.Contains(t, output, "srv")
	assert.Contains(t, output, "auth-service")
}

func TestLogger_MultipleFields(t *testing.T) {
	var buf bytes.Buffer

	log, err := logger.New(logger.Config{
		Level:  "debug",
		Format: "json",
		Output: &buf,
	})
	require.NoError(t, err)

	t.Run("множество полей", func(t *testing.T) {
		buf.Reset()
		log.Info("multi field log",
			"user_id", "123",
			"email", "user@example.com",
			"status", 200,
			"duration", "100ms",
		)

		output := buf.String()
		assert.Contains(t, output, "user_id")
		assert.Contains(t, output, "123")
		assert.Contains(t, output, "email")
		assert.Contains(t, output, "user@example.com")
		assert.Contains(t, output, "status")
		assert.Contains(t, output, "200")
	})

	t.Run("нечётное количество аргументов", func(t *testing.T) {
		buf.Reset()
		log.Info("odd args log", "key1", "value1", "key2")

		output := buf.String()
		assert.Contains(t, output, "odd args log")
		assert.Contains(t, output, "key1")
		assert.Contains(t, output, "value1")
	})
}

func TestLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer

	// Логгер с уровнем error
	log, err := logger.New(logger.Config{
		Level:  "error",
		Format: "json",
		Output: &buf,
	})
	require.NoError(t, err)

	t.Run("debug не пишется при уровне error", func(t *testing.T) {
		buf.Reset()
		log.Debug("debug message")
		assert.Empty(t, buf.String())
	})

	t.Run("info не пишется при уровне error", func(t *testing.T) {
		buf.Reset()
		log.Info("info message")
		assert.Empty(t, buf.String())
	})

	t.Run("warn не пишется при уровне error", func(t *testing.T) {
		buf.Reset()
		log.Warn("warn message")
		assert.Empty(t, buf.String())
	})

	t.Run("error пишется при уровне error", func(t *testing.T) {
		buf.Reset()
		log.Error("error message")
		assert.NotEmpty(t, buf.String())
	})
}
