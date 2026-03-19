// Package config предоставляет загрузку и валидацию конфигурации приложения.
//
// Пример использования:
//
//	cfg, err := config.Load("config.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Доступ к конфигурации
//	fmt.Println(cfg.Server.HTTPPort)
//	fmt.Println(cfg.Database.URL)
//
// Конфигурация также может быть загружена из переменных окружения (.env файл):
//
//	export JWT_SECRET="your-secret"
//	export DATABASE_URL="postgres://..."
//	cfg, err := config.LoadFromEnv()
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config полная конфигурация приложения
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Logging  LoggingConfig  `yaml:"logging"`
	CORS     CORSConfig     `yaml:"cors"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Health   HealthConfig   `yaml:"health"`
	Shutdown ShutdownConfig `yaml:"shutdown"`
}

// ServerConfig конфигурация сервера
type ServerConfig struct {
	HTTPPort     int    `yaml:"http_port"`
	GRPCPort     int    `yaml:"grpc_port"`
	Env          string `yaml:"env"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
	IdleTimeout  int    `yaml:"idle_timeout"`
}

// DatabaseConfig конфигурация PostgreSQL
type DatabaseConfig struct {
	URL              string `yaml:"url"`
	MaxConnections   int    `yaml:"max_connections"`
	ConnectionTimeout int   `yaml:"connection_timeout"`
}

// RedisConfig конфигурация Redis
type RedisConfig struct {
	URL              string `yaml:"url"`
	DB               int    `yaml:"db"`
	ConnectionTimeout int   `yaml:"connection_timeout"`
}

// JWTConfig конфигурация JWT
type JWTConfig struct {
	Secret       string `yaml:"secret"`
	AccessTTL    string `yaml:"access_ttl"`
	RefreshTTL   string `yaml:"refresh_ttl"`
	Issuer       string `yaml:"issuer"`
}

// LoggingConfig конфигурация логирования
type LoggingConfig struct {
	Level        string `yaml:"level"`
	Format       string `yaml:"format"`
	ServiceName  string `yaml:"service_name"`
}

// CORSConfig конфигурация CORS
type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// RateLimitConfig конфигурация rate limiting
type RateLimitConfig struct {
	Register int `yaml:"register"`
	Login    int `yaml:"login"`
	Refresh  int `yaml:"refresh"`
	Logout   int `yaml:"logout"`
}

// HealthConfig конфигурация health check
type HealthConfig struct {
	Path string `yaml:"path"`
}

// ShutdownConfig конфигурация graceful shutdown
type ShutdownConfig struct {
	Timeout int `yaml:"timeout"`
}

// Load загружает конфигурацию из YAML файла
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// LoadFromEnv загружает конфигурацию из переменных окружения
// Если .env файл существует, он будет загружен автоматически
func LoadFromEnv() (*Config, error) {
	// Пытаемся загрузить .env файл (не критично если не найден)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			HTTPPort:     getEnvInt("HTTP_PORT", 8080),
			GRPCPort:     getEnvInt("GRPC_PORT", 9090),
			Env:          getEnv("APP_ENV", "development"),
			ReadTimeout:  getEnvInt("READ_TIMEOUT", 10),
			WriteTimeout: getEnvInt("WRITE_TIMEOUT", 10),
			IdleTimeout:  getEnvInt("IDLE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			URL:              getEnv("DATABASE_URL", ""),
			MaxConnections:   getEnvInt("DATABASE_MAX_CONNECTIONS", 25),
			ConnectionTimeout: getEnvInt("DATABASE_CONNECTION_TIMEOUT", 10),
		},
		Redis: RedisConfig{
			URL:              getEnv("REDIS_URL", ""),
			DB:               getEnvInt("REDIS_DB", 0),
			ConnectionTimeout: getEnvInt("REDIS_CONNECTION_TIMEOUT", 5),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			AccessTTL:  getEnv("JWT_ACCESS_TTL", "15m"),
			RefreshTTL: getEnv("JWT_REFRESH_TTL", "336h"),
			Issuer:     getEnv("JWT_ISSUER", "auth-service"),
		},
		Logging: LoggingConfig{
			Level:       getEnv("LOG_LEVEL", "info"),
			Format:      getEnv("LOG_FORMAT", "json"),
			ServiceName: getEnv("LOG_SERVICE_NAME", "auth-service"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvSlice("CORS_ALLOWED_ORIGINS", ","),
			AllowedMethods: getEnvSlice("CORS_ALLOWED_METHODS", ","),
			AllowedHeaders: getEnvSlice("CORS_ALLOWED_HEADERS", ","),
			MaxAge:         getEnvInt("CORS_MAX_AGE", 86400),
		},
		RateLimit: RateLimitConfig{
			Register: getEnvInt("RATE_LIMIT_REGISTER", 5),
			Login:    getEnvInt("RATE_LIMIT_LOGIN", 10),
			Refresh:  getEnvInt("RATE_LIMIT_REFRESH", 30),
			Logout:   getEnvInt("RATE_LIMIT_LOGOUT", 60),
		},
		Health: HealthConfig{
			Path: getEnv("HEALTH_PATH", "/health"),
		},
		Shutdown: ShutdownConfig{
			Timeout: getEnvInt("SHUTDOWN_TIMEOUT", 30),
		},
	}

	// Парсим CORS allow credentials
	cfg.CORS.AllowCredentials = getEnvBool("CORS_ALLOW_CREDENTIALS", true)

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// getEnv получает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt получает целочисленное значение переменной окружения
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool получает булево значение переменной окружения
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		switch strings.ToLower(value) {
		case "true", "1", "yes":
			return true
		case "false", "0", "no":
			return false
		}
	}
	return defaultValue
}

// getEnvSlice получает срез строк из переменной окружения
func getEnvSlice(key, separator string) []string {
	if value := os.Getenv(key); value != "" {
		parts := strings.Split(value, separator)
		result := make([]string, len(parts))
		for i, part := range parts {
			result[i] = strings.TrimSpace(part)
		}
		return result
	}
	return []string{}
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
		c.Env = "development"
	}
	if c.ReadTimeout <= 0 {
		c.ReadTimeout = 10
	}
	if c.WriteTimeout <= 0 {
		c.WriteTimeout = 10
	}
	if c.IdleTimeout <= 0 {
		c.IdleTimeout = 60
	}
	return nil
}

// Validate валидирует конфигурацию базы данных
func (c *DatabaseConfig) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("url is required")
	}
	if c.MaxConnections <= 0 {
		c.MaxConnections = 25
	}
	if c.ConnectionTimeout <= 0 {
		c.ConnectionTimeout = 10
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
		c.ConnectionTimeout = 5
	}
	return nil
}

// Validate валидирует конфигурацию JWT
func (c *JWTConfig) Validate() error {
	// Сначала пробуем получить секрет из переменной окружения
	if c.Secret == "" {
		c.Secret = os.Getenv("JWT_SECRET")
	}
	
	if c.Secret == "" {
		return fmt.Errorf("secret is required (set JWT_SECRET env var or config file)")
	}
	
	if len(c.Secret) < 32 {
		return fmt.Errorf("secret must be at least 32 characters long")
	}
	if c.AccessTTL == "" {
		c.AccessTTL = "15m"
	}
	if c.RefreshTTL == "" {
		c.RefreshTTL = "336h"
	}
	if c.Issuer == "" {
		c.Issuer = "auth-service"
	}
	return nil
}

// Validate валидирует конфигурацию логирования
func (c *LoggingConfig) Validate() error {
	if c.Level == "" {
		c.Level = "info"
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
		c.Format = "json"
	}
	validFormats := map[string]bool{
		"json":    true,
		"console": true,
	}
	if !validFormats[c.Format] {
		return fmt.Errorf("format must be json or console")
	}
	if c.ServiceName == "" {
		c.ServiceName = "auth-service"
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

// ShutdownTimeoutDuration возвращает таймаут shutdown как time.Duration
func (c *ShutdownConfig) TimeoutDuration() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}
