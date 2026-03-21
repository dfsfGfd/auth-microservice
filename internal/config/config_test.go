package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"auth-microservice/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTempConfig(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")

	err := os.WriteFile(tmpFile, []byte(content), 0o644)
	require.NoError(t, err)

	return tmpFile
}

func TestLoad_Success(t *testing.T) {
	content := `
server:
  http_port: 8080
  grpc_port: 9090
  env: development

database:
  url: postgres://postgres:postgres@localhost:5432/auth?sslmode=disable
  max_connections: 25

redis:
  url: redis://localhost:6379
  db: 0

jwt:
  secret: super-secret-key-minimum-32-characters-long
  access_ttl: 15m
  refresh_ttl: 336h
  issuer: auth-service

logging:
  level: debug
  format: console
  service_name: auth-service

cors:
  allowed_origins:
    - http://localhost:3000
  allowed_methods:
    - GET
    - POST
  allowed_headers:
    - Authorization

shutdown:
  timeout: 30
`

	tmpFile := createTempConfig(t, content)

	cfg, err := config.Load(tmpFile)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Server
	assert.Equal(t, 8080, cfg.Server.HTTPPort)
	assert.Equal(t, 9090, cfg.Server.GRPCPort)
	assert.Equal(t, "development", cfg.Server.Env)

	// Database
	assert.Equal(t, "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable", cfg.Database.URL)
	assert.Equal(t, 25, cfg.Database.MaxConnections)

	// Redis
	assert.Equal(t, "redis://localhost:6379", cfg.Redis.URL)
	assert.Equal(t, 0, cfg.Redis.DB)

	// JWT
	assert.Equal(t, "super-secret-key-minimum-32-characters-long", cfg.JWT.Secret)
	assert.Equal(t, "15m", cfg.JWT.AccessTTL)
	assert.Equal(t, "336h", cfg.JWT.RefreshTTL)
	assert.Equal(t, "auth-service", cfg.JWT.Issuer)

	// Logging
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.Equal(t, "console", cfg.Logging.Format)
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("nonexistent.yaml")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpFile := createTempConfig(t, "invalid: yaml: content: [")

	_, err := config.Load(tmpFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse config file")
}

func TestConfig_Validate(t *testing.T) {
	t.Run("валидная конфигурация", func(t *testing.T) {
		cfg := &config.Config{
			Server:   config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090},
			Database: config.DatabaseConfig{URL: "postgres://localhost/auth"},
			Redis:    config.RedisConfig{URL: "redis://localhost:6379"},
			JWT:      config.JWTConfig{Secret: "super-secret-key-minimum-32-characters-long"},
			Logging:  config.LoggingConfig{Level: "info", Format: "json"},
		}

		err := cfg.Validate()
		assert.NoError(t, err)
	})

	t.Run("отсутствует database url", func(t *testing.T) {
		cfg := &config.Config{
			Server:   config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090},
			Database: config.DatabaseConfig{},
			Redis:    config.RedisConfig{URL: "redis://localhost:6379"},
			JWT:      config.JWTConfig{Secret: "super-secret-key-minimum-32-characters-long"},
			Logging:  config.LoggingConfig{Level: "info", Format: "json"},
		}

		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database")
		assert.Contains(t, err.Error(), "url is required")
	})

	t.Run("отсутствует redis url", func(t *testing.T) {
		cfg := &config.Config{
			Server:   config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090},
			Database: config.DatabaseConfig{URL: "postgres://localhost/auth"},
			Redis:    config.RedisConfig{},
			JWT:      config.JWTConfig{Secret: "super-secret-key-minimum-32-characters-long"},
			Logging:  config.LoggingConfig{Level: "info", Format: "json"},
		}

		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis")
		assert.Contains(t, err.Error(), "url is required")
	})

	t.Run("короткий jwt secret", func(t *testing.T) {
		cfg := &config.Config{
			Server:   config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090},
			Database: config.DatabaseConfig{URL: "postgres://localhost/auth"},
			Redis:    config.RedisConfig{URL: "redis://localhost:6379"},
			JWT:      config.JWTConfig{Secret: "short"},
			Logging:  config.LoggingConfig{Level: "info", Format: "json"},
		}

		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "jwt")
		assert.Contains(t, err.Error(), "at least 32 characters")
	})
}

func TestServerConfig_Validate(t *testing.T) {
	t.Run("невалидный http_port", func(t *testing.T) {
		cfg := config.ServerConfig{HTTPPort: 0, GRPCPort: 9090}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "http_port")
	})

	t.Run("невалидный grpc_port", func(t *testing.T) {
		cfg := config.ServerConfig{HTTPPort: 8080, GRPCPort: 70000}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "grpc_port")
	})

	t.Run("env по умолчанию", func(t *testing.T) {
		cfg := config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090, Env: ""}
		err := cfg.Validate()
		assert.NoError(t, err)
		assert.Equal(t, "development", cfg.Env)
	})

	t.Run("таймауты по умолчанию", func(t *testing.T) {
		cfg := config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090}
		err := cfg.Validate()
		assert.NoError(t, err)
		assert.Equal(t, 10, cfg.ReadTimeout)
		assert.Equal(t, 10, cfg.WriteTimeout)
		assert.Equal(t, 60, cfg.IdleTimeout)
	})
}

func TestDatabaseConfig_Validate(t *testing.T) {
	t.Run("отсутствует url", func(t *testing.T) {
		cfg := config.DatabaseConfig{}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "url is required")
	})

	t.Run("max_connections по умолчанию", func(t *testing.T) {
		cfg := config.DatabaseConfig{URL: "postgres://localhost/auth"}
		err := cfg.Validate()
		assert.NoError(t, err)
		assert.Equal(t, 25, cfg.MaxConnections)
	})
}

func TestRedisConfig_Validate(t *testing.T) {
	t.Run("отсутствует url", func(t *testing.T) {
		cfg := config.RedisConfig{}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "url is required")
	})

	t.Run("невалидный db", func(t *testing.T) {
		cfg := config.RedisConfig{URL: "redis://localhost", DB: 20}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db must be between 0 and 15")
	})
}

func TestJWTConfig_Validate(t *testing.T) {
	t.Run("отсутствует secret", func(t *testing.T) {
		cfg := config.JWTConfig{}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "secret is required")
	})

	t.Run("короткий secret", func(t *testing.T) {
		cfg := config.JWTConfig{Secret: "short"}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least 32 characters")
	})

	t.Run("значения по умолчанию", func(t *testing.T) {
		cfg := config.JWTConfig{Secret: "super-secret-key-minimum-32-characters-long"}
		err := cfg.Validate()
		assert.NoError(t, err)
		assert.Equal(t, "15m", cfg.AccessTTL)
		assert.Equal(t, "336h", cfg.RefreshTTL)
		assert.Equal(t, "auth-service", cfg.Issuer)
	})
}

func TestLoggingConfig_Validate(t *testing.T) {
	t.Run("level по умолчанию", func(t *testing.T) {
		cfg := config.LoggingConfig{}
		err := cfg.Validate()
		assert.NoError(t, err)
		assert.Equal(t, "info", cfg.Level)
	})

	t.Run("format по умолчанию", func(t *testing.T) {
		cfg := config.LoggingConfig{}
		err := cfg.Validate()
		assert.NoError(t, err)
		assert.Equal(t, "json", cfg.Format)
	})

	t.Run("service_name по умолчанию", func(t *testing.T) {
		cfg := config.LoggingConfig{}
		err := cfg.Validate()
		assert.NoError(t, err)
		assert.Equal(t, "auth-service", cfg.ServiceName)
	})

	t.Run("невалидный level", func(t *testing.T) {
		cfg := config.LoggingConfig{Level: "invalid"}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "level must be")
	})

	t.Run("невалидный format", func(t *testing.T) {
		cfg := config.LoggingConfig{Level: "info", Format: "xml"}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "format must be")
	})
}

func TestJWTConfig_Duration(t *testing.T) {
	t.Run("AccessTTLDuration", func(t *testing.T) {
		cfg := config.JWTConfig{AccessTTL: "15m"}
		duration, err := cfg.AccessTTLDuration()
		require.NoError(t, err)
		assert.Equal(t, 15*time.Minute, duration)
	})

	t.Run("RefreshTTLDuration", func(t *testing.T) {
		cfg := config.JWTConfig{RefreshTTL: "336h"}
		duration, err := cfg.RefreshTTLDuration()
		require.NoError(t, err)
		assert.Equal(t, 336*time.Hour, duration)
	})
}

func TestServerConfig_Duration(t *testing.T) {
	cfg := config.ServerConfig{
		HTTPPort:     8080,
		GRPCPort:     9090,
		ReadTimeout:  10,
		WriteTimeout: 20,
		IdleTimeout:  60,
	}

	t.Run("ReadTimeoutDuration", func(t *testing.T) {
		assert.Equal(t, 10*time.Second, cfg.ReadTimeoutDuration())
	})

	t.Run("WriteTimeoutDuration", func(t *testing.T) {
		assert.Equal(t, 20*time.Second, cfg.WriteTimeoutDuration())
	})

	t.Run("IdleTimeoutDuration", func(t *testing.T) {
		assert.Equal(t, 60*time.Second, cfg.IdleTimeoutDuration())
	})
}

func TestShutdownConfig_Duration(t *testing.T) {
	cfg := config.ShutdownConfig{Timeout: 30}
	assert.Equal(t, 30*time.Second, cfg.TimeoutDuration())
}
