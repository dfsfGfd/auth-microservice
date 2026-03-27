// Package config предоставляет загрузку конфигурации из переменных окружения.
//
// Пример использования:
//
//	cfg, err := config.Load()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Доступ к конфигурации
//	fmt.Println(cfg.Server.HTTPPort)
//	fmt.Println(cfg.Database.URL)
package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config полная конфигурация приложения
type Config struct {
	Server    ServerConfig    `envPrefix:"SERVER_"`
	Database  DatabaseConfig  `envPrefix:"DATABASE_"`
	Redis     RedisConfig     `envPrefix:"REDIS_"`
	JWT       JWTConfig       `envPrefix:"JWT_"`
	Logging   LoggingConfig   `envPrefix:"LOG_"`
	RateLimit RateLimitConfig `envPrefix:"RATE_LIMIT_"`
	Health    HealthConfig    `envPrefix:"HEALTH_"`
	Shutdown  ShutdownConfig  `envPrefix:"SHUTDOWN_"`
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	HTTPPort     int    `env:"HTTP_PORT" envDefault:"8080"`
	GRPCPort     int    `env:"GRPC_PORT" envDefault:"9090"`
	Env          string `env:"ENV" envDefault:"development"`
	ReadTimeout  int    `env:"READ_TIMEOUT" envDefault:"10"`
	WriteTimeout int    `env:"WRITE_TIMEOUT" envDefault:"10"`
	IdleTimeout  int    `env:"IDLE_TIMEOUT" envDefault:"60"`
}

// DatabaseConfig конфигурация PostgreSQL
type DatabaseConfig struct {
	URL               string `env:"URL,required"`
	MaxConnections    int    `env:"MAX_CONNECTIONS" envDefault:"25"`
	ConnectionTimeout int    `env:"CONNECTION_TIMEOUT" envDefault:"10"`
}

// RedisConfig конфигурация Redis
type RedisConfig struct {
	URL               string `env:"URL,required"`
	DB                int    `env:"DB" envDefault:"0"`
	ConnectionTimeout int    `env:"CONNECTION_TIMEOUT" envDefault:"5"`
}

// JWTConfig конфигурация JWT
type JWTConfig struct {
	Secret     string `env:"SECRET,required"`
	AccessTTL  string `env:"ACCESS_TTL" envDefault:"15m"`
	RefreshTTL string `env:"REFRESH_TTL" envDefault:"336h"`
	Issuer     string `env:"ISSUER" envDefault:"auth-service"`
}

// LoggingConfig конфигурация логирования
type LoggingConfig struct {
	Level       string `env:"LEVEL" envDefault:"info"`
	Format      string `env:"FORMAT" envDefault:"json"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"auth-service"`
}

// RateLimitConfig конфигурация rate limiting
type RateLimitConfig struct {
	Register int `env:"REGISTER" envDefault:"5"`
	Login    int `env:"LOGIN" envDefault:"10"`
	Refresh  int `env:"REFRESH" envDefault:"30"`
	Logout   int `env:"LOGOUT" envDefault:"60"`
}

// HealthConfig конфигурация health check
type HealthConfig struct {
	Path string `env:"PATH" envDefault:"/health"`
}

// ShutdownConfig конфигурация graceful shutdown
type ShutdownConfig struct {
	Timeout int `env:"TIMEOUT" envDefault:"30"`
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	return &cfg, nil
}

// Validate валидирует конфигурацию
func (c *Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server: %w", err)
	}
	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database: %w", err)
	}
	if err := c.Redis.Validate(); err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	if err := c.JWT.Validate(); err != nil {
		return fmt.Errorf("jwt: %w", err)
	}
	if err := c.Logging.Validate(); err != nil {
		return fmt.Errorf("logging: %w", err)
	}
	return nil
}

// Validate валидирует конфигурацию сервера
func (c *ServerConfig) Validate() error {
	if c.HTTPPort <= 0 || c.HTTPPort > 65535 {
		return fmt.Errorf("http_port must be between 1 and 65535")
	}
	if c.GRPCPort <= 0 || c.GRPCPort > 65535 {
		return fmt.Errorf("grpc_port must be between 1 and 65535")
	}
	if c.Env == "" {
		return fmt.Errorf("env is required")
	}
	if c.ReadTimeout <= 0 {
		return fmt.Errorf("read_timeout must be positive")
	}
	if c.WriteTimeout <= 0 {
		return fmt.Errorf("write_timeout must be positive")
	}
	if c.IdleTimeout <= 0 {
		return fmt.Errorf("idle_timeout must be positive")
	}
	return nil
}

// Validate валидирует конфигурацию базы данных
func (c *DatabaseConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("url is required")
	}
	if c.MaxConnections <= 0 {
		return fmt.Errorf("max_connections must be positive")
	}
	if c.ConnectionTimeout <= 0 {
		return fmt.Errorf("connection_timeout must be positive")
	}
	return nil
}

// Validate валидирует конфигурацию Redis
func (c *RedisConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("url is required")
	}
	if c.DB < 0 || c.DB > 15 {
		return fmt.Errorf("db must be between 0 and 15")
	}
	if c.ConnectionTimeout <= 0 {
		return fmt.Errorf("connection_timeout must be positive")
	}
	return nil
}

// Validate валидирует конфигурацию JWT
func (c *JWTConfig) Validate() error {
	if c.Secret == "" {
		return fmt.Errorf("secret is required")
	}
	if len(c.Secret) < 32 {
		return fmt.Errorf("secret must be at least 32 characters long")
	}
	if c.AccessTTL == "" {
		return fmt.Errorf("access_ttl is required")
	}
	if c.RefreshTTL == "" {
		return fmt.Errorf("refresh_ttl is required")
	}
	if c.Issuer == "" {
		return fmt.Errorf("issuer is required")
	}
	return nil
}

// Validate валидирует конфигурацию логирования
func (c *LoggingConfig) Validate() error {
	if c.Level == "" {
		return fmt.Errorf("level is required")
	}
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLevels[c.Level] {
		return fmt.Errorf("level must be debug, info, warn, error, or fatal")
	}
	if c.Format == "" {
		return fmt.Errorf("format is required")
	}
	validFormats := map[string]bool{
		"json":    true,
		"console": true,
	}
	if !validFormats[c.Format] {
		return fmt.Errorf("format must be json or console")
	}
	if c.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}
	return nil
}

// AccessTTLDuration возвращает время жизни access токена как time.Duration
func (c *JWTConfig) AccessTTLDuration() (time.Duration, error) {
	return time.ParseDuration(c.AccessTTL)
}

// RefreshTTLDuration возвращает время жизни refresh токена как time.Duration
func (c *JWTConfig) RefreshTTLDuration() (time.Duration, error) {
	return time.ParseDuration(c.RefreshTTL)
}

// ReadTimeoutDuration возвращает таймаут чтения как time.Duration
func (c *ServerConfig) ReadTimeoutDuration() time.Duration {
	return time.Duration(c.ReadTimeout) * time.Second
}

// WriteTimeoutDuration возвращает таймаут записи как time.Duration
func (c *ServerConfig) WriteTimeoutDuration() time.Duration {
	return time.Duration(c.WriteTimeout) * time.Second
}

// IdleTimeoutDuration возвращает таймаут простоя как time.Duration
func (c *ServerConfig) IdleTimeoutDuration() time.Duration {
	return time.Duration(c.IdleTimeout) * time.Second
}

// TimeoutDuration возвращает таймаут shutdown как time.Duration
func (c *ShutdownConfig) TimeoutDuration() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}
