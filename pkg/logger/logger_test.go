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
	assert.Contains(t, output, "msg")
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "key")
	assert.Contains(t, output, "value")
	assert.Contains(t, output, "ts")
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

func TestLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer

	log, err := logger.New(logger.Config{
		Level:  "error",
		Format: "json",
		Output: &buf,
	})
	require.NoError(t, err)

	t.Run("warn не пишется при уровне error", func(t *testing.T) {
		buf.Reset()
		log.Warn("warn message")
		assert.Empty(t, buf.String())
	})

	t.Run("info не пишется при уровне error", func(t *testing.T) {
		buf.Reset()
		log.Info("info message")
		assert.Empty(t, buf.String())
	})

	t.Run("error пишется при уровне error", func(t *testing.T) {
		buf.Reset()
		log.Error("error message")
		assert.NotEmpty(t, buf.String())
	})
}
