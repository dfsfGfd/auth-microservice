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
		"HTTP_PORT":                     "8080",
		"GRPC_PORT":                     "9090",
		"ENV":                           "development",
		"READ_TIMEOUT":                  "10",
		"WRITE_TIMEOUT":                 "10",
		"IDLE_TIMEOUT":                  "60",
		"DATABASE_URL":                  "postgres://postgres:postgres@localhost:5432/auth?sslmode=disable",
		"DATABASE_MAX_CONNECTIONS":      "25",
		"DATABASE_CONNECTION_TIMEOUT":   "10",
		"REDIS_URL":                     "redis://localhost:6379",
		"REDIS_DB":                      "0",
		"REDIS_CONNECTION_TIMEOUT":      "5",
		"JWT_SECRET":                    "super-secret-key-minimum-32-characters-long",
		"JWT_ACCESS_TTL":                "15m",
		"JWT_REFRESH_TTL":               "336h",
		"JWT_ISSUER":                    "auth-service",
		"LOG_LEVEL":                     "debug",
		"LOG_FORMAT":                    "console",
		"LOG_SERVICE_NAME":              "auth-service",
		"RATE_LIMIT_REGISTER":           "5",
		"RATE_LIMIT_LOGIN":              "10",
		"RATE_LIMIT_REFRESH":            "30",
		"RATE_LIMIT_LOGOUT":             "60",
		"HEALTH_PATH":                   "/health",
		"SHUTDOWN_TIMEOUT":              "30",
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
			Server:   config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090, Env: "development", ReadTimeout: 10, WriteTimeout: 10, IdleTimeout: 60},
			Database: config.DatabaseConfig{URL: "postgres://localhost/auth", MaxConnections: 25, MinConnections: 0, ConnectionTimeout: 10, MaxConnLifetime: 1800, MaxConnIdleTime: 300},
			Redis:    config.RedisConfig{URL: "redis://localhost:6379", DB: 0, PoolSize: 10, ConnectionTimeout: 5, ReadTimeout: 3, WriteTimeout: 3},
			JWT:      config.JWTConfig{Secret: "super-secret-key-minimum-32-characters-long", AccessTTL: "15m", RefreshTTL: "336h", Issuer: "auth-service"},
			Logging:  config.LoggingConfig{Level: "info", Format: "json", ServiceName: "auth-service"},
			RateLimit: config.RateLimitConfig{Register: 5, Login: 10, Refresh: 30, Logout: 60},
			Health:   config.HealthConfig{Path: "/health"},
			Shutdown: config.ShutdownConfig{Timeout: 30},
		}

		err := cfg.Validate()
		assert.NoError(t, err)
	})

	t.Run("отсутствует database url", func(t *testing.T) {
		cfg := &config.Config{
			Server:   config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090, Env: "development", ReadTimeout: 10, WriteTimeout: 10, IdleTimeout: 60},
			Database: config.DatabaseConfig{},
			Redis:    config.RedisConfig{URL: "redis://localhost:6379", DB: 0, PoolSize: 10, ConnectionTimeout: 5, ReadTimeout: 3, WriteTimeout: 3},
			JWT:      config.JWTConfig{Secret: "super-secret-key-minimum-32-characters-long", AccessTTL: "15m", RefreshTTL: "336h", Issuer: "auth-service"},
			Logging:  config.LoggingConfig{Level: "info", Format: "json", ServiceName: "auth-service"},
		}

		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database")
		assert.Contains(t, err.Error(), "url is required")
	})

	t.Run("отсутствует redis url", func(t *testing.T) {
		cfg := &config.Config{
			Server:   config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090, Env: "development", ReadTimeout: 10, WriteTimeout: 10, IdleTimeout: 60},
			Database: config.DatabaseConfig{URL: "postgres://localhost/auth", MaxConnections: 25, MinConnections: 0, ConnectionTimeout: 10, MaxConnLifetime: 1800, MaxConnIdleTime: 300},
			Redis:    config.RedisConfig{},
			JWT:      config.JWTConfig{Secret: "super-secret-key-minimum-32-characters-long", AccessTTL: "15m", RefreshTTL: "336h", Issuer: "auth-service"},
			Logging:  config.LoggingConfig{Level: "info", Format: "json", ServiceName: "auth-service"},
		}

		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "redis")
		assert.Contains(t, err.Error(), "url is required")
	})

	t.Run("короткий jwt secret", func(t *testing.T) {
		cfg := &config.Config{
			Server:   config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090, Env: "development", ReadTimeout: 10, WriteTimeout: 10, IdleTimeout: 60},
			Database: config.DatabaseConfig{URL: "postgres://localhost/auth", MaxConnections: 25, MinConnections: 0, ConnectionTimeout: 10, MaxConnLifetime: 1800, MaxConnIdleTime: 300},
			Redis:    config.RedisConfig{URL: "redis://localhost:6379", DB: 0, PoolSize: 10, ConnectionTimeout: 5, ReadTimeout: 3, WriteTimeout: 3},
			JWT:      config.JWTConfig{Secret: "short", AccessTTL: "15m", RefreshTTL: "336h", Issuer: "auth-service"},
			Logging:  config.LoggingConfig{Level: "info", Format: "json", ServiceName: "auth-service"},
		}

		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "jwt")
		assert.Contains(t, err.Error(), "at least 32 characters")
	})
}

func TestServerConfig_Validate(t *testing.T) {
	t.Run("невалидный http_port", func(t *testing.T) {
		cfg := config.ServerConfig{HTTPPort: 0, GRPCPort: 9090, Env: "development", ReadTimeout: 10, WriteTimeout: 10, IdleTimeout: 60}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "http_port")
	})

	t.Run("невалидный grpc_port", func(t *testing.T) {
		cfg := config.ServerConfig{HTTPPort: 8080, GRPCPort: 70000, Env: "development", ReadTimeout: 10, WriteTimeout: 10, IdleTimeout: 60}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "grpc_port")
	})

	t.Run("отсутствует env", func(t *testing.T) {
		cfg := config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090, Env: "", ReadTimeout: 10, WriteTimeout: 10, IdleTimeout: 60}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "env is required")
	})

	t.Run("невалидный read_timeout", func(t *testing.T) {
		cfg := config.ServerConfig{HTTPPort: 8080, GRPCPort: 9090, Env: "development", ReadTimeout: 0, WriteTimeout: 10, IdleTimeout: 60}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "read_timeout must be positive")
	})
}

func TestDatabaseConfig_Validate(t *testing.T) {
	t.Run("отсутствует url", func(t *testing.T) {
		cfg := config.DatabaseConfig{}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "url is required")
	})

	t.Run("невалидный max_connections", func(t *testing.T) {
		cfg := config.DatabaseConfig{URL: "postgres://localhost/auth", MaxConnections: 0, MinConnections: 0, ConnectionTimeout: 10, MaxConnLifetime: 1800, MaxConnIdleTime: 300}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max_connections must be positive")
	})

	t.Run("невалидный min_connections", func(t *testing.T) {
		cfg := config.DatabaseConfig{URL: "postgres://localhost/auth", MaxConnections: 25, MinConnections: -1, ConnectionTimeout: 10, MaxConnLifetime: 1800, MaxConnIdleTime: 300}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "min_connections must be non-negative")
	})

	t.Run("min_connections > max_connections", func(t *testing.T) {
		cfg := config.DatabaseConfig{URL: "postgres://localhost/auth", MaxConnections: 5, MinConnections: 10, ConnectionTimeout: 10, MaxConnLifetime: 1800, MaxConnIdleTime: 300}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "min_connections cannot exceed max_connections")
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
		cfg := config.RedisConfig{URL: "redis://localhost", DB: 20, PoolSize: 10, ConnectionTimeout: 5, ReadTimeout: 3, WriteTimeout: 3}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db must be between 0 and 15")
	})

	t.Run("невалидный pool_size", func(t *testing.T) {
		cfg := config.RedisConfig{URL: "redis://localhost", DB: 0, PoolSize: 0, ConnectionTimeout: 5, ReadTimeout: 3, WriteTimeout: 3}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "pool_size must be positive")
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
		cfg := config.JWTConfig{Secret: "short", AccessTTL: "15m", RefreshTTL: "336h", Issuer: "auth-service"}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least 32 characters")
	})

	t.Run("отсутствует access_ttl", func(t *testing.T) {
		cfg := config.JWTConfig{Secret: "super-secret-key-minimum-32-characters-long", AccessTTL: "", RefreshTTL: "336h", Issuer: "auth-service"}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "access_ttl is required")
	})

	t.Run("валидный jwt config", func(t *testing.T) {
		cfg := config.JWTConfig{Secret: "super-secret-key-minimum-32-characters-long", AccessTTL: "15m", RefreshTTL: "336h", Issuer: "auth-service"}
		err := cfg.Validate()
		assert.NoError(t, err)
	})
}

func TestLoggingConfig_Validate(t *testing.T) {
	t.Run("отсутствует level", func(t *testing.T) {
		cfg := config.LoggingConfig{}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "level is required")
	})

	t.Run("невалидный level", func(t *testing.T) {
		cfg := config.LoggingConfig{Level: "invalid", Format: "json", ServiceName: "auth-service"}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "level must be")
	})

	t.Run("невалидный format", func(t *testing.T) {
		cfg := config.LoggingConfig{Level: "info", Format: "xml", ServiceName: "auth-service"}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "format must be")
	})

	t.Run("валидный logging config", func(t *testing.T) {
		cfg := config.LoggingConfig{Level: "info", Format: "json", ServiceName: "auth-service"}
		err := cfg.Validate()
		assert.NoError(t, err)
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
