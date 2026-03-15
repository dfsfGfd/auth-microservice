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
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config полная конфигурация приложения
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Cookie   CookieConfig   `yaml:"cookie"`
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

// CookieConfig конфигурация cookie
type CookieConfig struct {
	Secure    bool   `yaml:"secure"`
	HTTPOnly  bool   `yaml:"http_only"`
	SameSite  string `yaml:"same_site"`
	Domain    string `yaml:"domain"`
	Path      string `yaml:"path"`
	MaxAge    int    `yaml:"max_age"`
}

// LoggingConfig конфигурация логирования
type LoggingConfig struct {
	Level        string `yaml:"level"`
	Format       string `yaml:"format"`
	ServiceName  string `yaml:"service_name"`
}

// CORSConfig конфигурация CORS
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
	MaxAge         int      `yaml:"max_age"`
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
	if err := c.Cookie.Validate(); err != nil {
		return fmt.Errorf("cookie: %w", err)
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
	if c.Secret == "" {
		return fmt.Errorf("secret is required")
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

// Validate валидирует конфигурацию cookie
func (c *CookieConfig) Validate() error {
	if c.Path == "" {
		c.Path = "/"
	}
	if c.MaxAge <= 0 {
		c.MaxAge = int((14 * 24 * time.Hour).Seconds())
	}
	if c.SameSite == "" {
		c.SameSite = "Lax"
	}
	validSameSite := map[string]bool{
		"Strict": true,
		"Lax":    true,
		"None":   true,
	}
	if !validSameSite[c.SameSite] {
		return fmt.Errorf("same_site must be Strict, Lax, or None")
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
