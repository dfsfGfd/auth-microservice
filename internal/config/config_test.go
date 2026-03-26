package config_test

import (
	"testing"
	"time"

	"auth-microservice/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setTestEnv(t *testing.T) {
	t.Helper()

	envs := map[string]string{
		"HTTP_PORT":                "8080",
		"GRPC_PORT":                "9090",
		"APP_ENV":                  "development",
		"DATABASE_URL":             "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable",
		"DATABASE_MAX_CONNECTIONS": "25",
		"REDIS_URL":                "redis://localhost:6379",
		"REDIS_DB":                 "0",
		"JWT_SECRET":               "super-secret-key-minimum-32-characters-long",
		"JWT_ACCESS_TTL":           "15m",
		"JWT_REFRESH_TTL":          "336h",
		"JWT_ISSUER":               "auth-service",
		"LOG_LEVEL":                "debug",
		"LOG_FORMAT":               "console",
		"LOG_SERVICE_NAME":         "auth-service",
		"CORS_ALLOWED_ORIGINS":     "http://localhost:3000",
		"CORS_ALLOW_CREDENTIALS":   "false",
		"HEALTH_PATH":              "/health",
		"SHUTDOWN_TIMEOUT":         "30",
	}

	for k, v := range envs {
		t.Setenv(k, v)
	}
}

func TestLoad(t *testing.T) {
	setTestEnv(t)

	cfg, err := config.Load()

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

func TestLoad_Defaults(t *testing.T) {
	// Устанавливаем только обязательные переменные
	t.Setenv("DATABASE_URL", "postgres://localhost/auth")
	t.Setenv("REDIS_URL", "redis://localhost:6379")
	t.Setenv("JWT_SECRET", "super-secret-key-minimum-32-characters-long")

	cfg, err := config.Load()

	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Проверяем дефолты
	assert.Equal(t, 8080, cfg.Server.HTTPPort)
	assert.Equal(t, 9090, cfg.Server.GRPCPort)
	assert.Equal(t, "development", cfg.Server.Env)
	assert.Equal(t, 25, cfg.Database.MaxConnections)
	assert.Equal(t, "info", cfg.Logging.Level)
	assert.Equal(t, "json", cfg.Logging.Format)
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

	t.Run("wildcard origin с credentials", func(t *testing.T) {
		cfg := &config.Config{
			Server:   config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090},
			Database: config.DatabaseConfig{URL: "postgres://localhost/auth"},
			Redis:    config.RedisConfig{URL: "redis://localhost:6379"},
			JWT:      config.JWTConfig{Secret: "super-secret-key-minimum-32-characters-long"},
			Logging:  config.LoggingConfig{Level: "info", Format: "json"},
			CORS: config.CORSConfig{
				AllowedOrigins:   []string{"*"},
				AllowCredentials: true,
			},
		}

		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wildcard")
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
