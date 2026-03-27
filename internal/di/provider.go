// Package di предоставляет dependency injection для приложения.
//
// Использование:
//  1. Запустить wire: wire gen ./internal/di
//  2. Скомпилировать: go build ./cmd/server
package di

import (
	"context"
	"fmt"
	"strings"
	"time"

	"auth-microservice/internal/cache/token"
	"auth-microservice/internal/config"
	"auth-microservice/internal/handler/auth"
	"auth-microservice/internal/middleware"
	"auth-microservice/internal/repository"
	repositoryAuth "auth-microservice/internal/repository/auth"
	serviceAuth "auth-microservice/internal/service/auth"
	"auth-microservice/pkg/bcrypt"
	"auth-microservice/pkg/db/postgresql"
	"auth-microservice/pkg/db/redisdb"
	"auth-microservice/pkg/jwt"
	"auth-microservice/pkg/logger"
	"auth-microservice/pkg/proto/auth/v1"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"
)

// Application содержит все зависимости приложения
type Application struct {
	Config      *config.Config
	Logger      *logger.Logger
	JWTService  *jwt.Service
	DB          *pgxpool.Pool
	Redis       *goredis.Client
	AccountRepo repository.AccountRepository
	TokenCache  *token.RedisCache
	AuthService *serviceAuth.AuthService
	AuthHandler authv1.AuthServiceServer
	RateLimiter *middleware.RateLimiter
}

// CleanUp очищает ресурсы приложения
func (a *Application) CleanUp() error {
	// Закрываем подключения
	if a.DB != nil {
		a.DB.Close()
	}
	if a.Redis != nil {
		a.Redis.Close()
	}
	return nil
}

// ProvideContext предоставляет контекст для инициализации
func ProvideContext() context.Context {
	return context.Background()
}

// ProviderSet набор провайдеров для DI
var ProviderSet = wire.NewSet(
	// Конфигурация
	loadConfig,

	// Подключения
	ProvidePostgresConfig,
	ProvideRedisConfig,
	ProvidePostgresPool,
	ProvideRedisClient,
	ProvideContext,

	// Репозитории
	repositoryAuth.NewAccountRepository,

	// Кэш токенов
	ProvideTokenCachePrefix,
	token.NewRedisCache,

	// Логгер
	NewLogger,

	// JWT сервис
	NewJWTService,

	// Bcrypt hasher
	bcrypt.NewService,

	// Rate Limiter
	ProvideRateLimitConfigs,
	middleware.NewRateLimiter,

	// AuthService
	serviceAuth.NewAuthService,

	// AuthHandler
	auth.NewHandler,

	// Application
	NewApplication,
)

// loadConfig загружает конфигурацию из .env файла
func loadConfig() (*config.Config, error) {
	return config.Load()
}

// ProvidePostgresConfig предоставляет конфигурацию PostgreSQL
func ProvidePostgresConfig(cfg *config.Config) postgresql.Config {
	return postgresql.Config{
		DSN:         cfg.Database.URL,
		MaxConns:    int32(cfg.Database.MaxConnections),
		ConnTimeout: time.Duration(cfg.Database.ConnectionTimeout) * time.Second,
	}
}

// ProvideRedisConfig предоставляет конфигурацию Redis
func ProvideRedisConfig(cfg *config.Config) redisdb.Config {
	// Удаляем протокол redis:// из URL если есть
	addr := strings.TrimPrefix(cfg.Redis.URL, "redis://")

	return redisdb.Config{
		Addr:         addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		ConnTimeout:  time.Duration(cfg.Redis.ConnectionTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Redis.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Redis.WriteTimeout) * time.Second,
	}
}

// ProvidePostgresPool создаёт пул подключений к PostgreSQL
func ProvidePostgresPool(ctx context.Context, cfg postgresql.Config) (*pgxpool.Pool, error) {
	return postgresql.NewPool(ctx, cfg)
}

// ProvideRedisClient создаёт Redis клиент
func ProvideRedisClient(ctx context.Context, cfg redisdb.Config) (*goredis.Client, error) {
	return redisdb.NewClient(ctx, cfg)
}

// ProvideTokenCachePrefix предоставляет префикс для ключей токенов
func ProvideTokenCachePrefix() string {
	return "refresh:"
}

// NewLogger создаёт логгер из конфигурации
func NewLogger(cfg *config.Config) (*logger.Logger, error) {
	return logger.New(logger.Config{
		Level:       cfg.Logging.Level,
		Format:      cfg.Logging.Format,
		ServiceName: cfg.Logging.ServiceName,
	})
}

// NewJWTService создаёт JWT сервис из конфигурации
func NewJWTService(cfg *config.Config) (*jwt.Service, error) {
	accessTTL, err := cfg.JWT.AccessTTLDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid access_ttl: %w", err)
	}

	refreshTTL, err := cfg.JWT.RefreshTTLDuration()
	if err != nil {
		return nil, fmt.Errorf("invalid refresh_ttl: %w", err)
	}

	return jwt.NewService(jwt.Config{
		SecretKey:       cfg.JWT.Secret,
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
		Issuer:          cfg.JWT.Issuer,
	})
}

// ProvideRateLimitConfigs предоставляет конфигурации rate limiter для всех endpoint'ов
func ProvideRateLimitConfigs(cfg *config.Config) map[string]middleware.RateLimiterConfig {
	return map[string]middleware.RateLimiterConfig{
		"register": {
			Window: time.Minute,
			Limit:  cfg.RateLimit.Register,
			Prefix: "ratelimit:register:",
		},
		"login": {
			Window: time.Minute,
			Limit:  cfg.RateLimit.Login,
			Prefix: "ratelimit:login:",
		},
		"refresh": {
			Window: time.Minute,
			Limit:  cfg.RateLimit.Refresh,
			Prefix: "ratelimit:refresh:",
		},
		"logout": {
			Window: time.Minute,
			Limit:  cfg.RateLimit.Logout,
			Prefix: "ratelimit:logout:",
		},
	}
}

// NewApplication создаёт приложение из зависимостей
func NewApplication(
	cfg *config.Config,
	log *logger.Logger,
	jwtSvc *jwt.Service,
	db *pgxpool.Pool,
	redisClient *goredis.Client,
	accountRepo repository.AccountRepository,
	tokenCache *token.RedisCache,
	rateLimiter *middleware.RateLimiter,
	authService *serviceAuth.AuthService,
	authHandler *auth.Handler,
) (*Application, error) {
	return &Application{
		Config:      cfg,
		Logger:      log,
		JWTService:  jwtSvc,
		DB:          db,
		Redis:       redisClient,
		AccountRepo: accountRepo,
		TokenCache:  tokenCache,
		RateLimiter: rateLimiter,
		AuthService: authService,
		AuthHandler: authHandler,
	}, nil
}
