// Package di предоставляет dependency injection для приложения.
//
// Использование:
//   1. Запустить wiregen: go generate ./...
//   2. Скомпилировать: go build ./cmd/server
//
//go:generate go run github.com/google/wire/cmd/wiregen
package di

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/wire"
	goredis "github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5/pgxpool"

	"auth-microservice/internal/config"
	"auth-microservice/pkg/db/postgres"
	dbredis "auth-microservice/pkg/db/redis"
	"auth-microservice/pkg/cookies"
	"auth-microservice/pkg/jwt"
	"auth-microservice/pkg/logger"
)

// Application содержит все зависимости приложения
type Application struct {
	Config        *config.Config
	Logger        *logger.Logger
	JWTService    *jwt.Service
	CookieService *cookies.Service
	DB            *pgxpool.Pool
	Redis         *goredis.Client
	// TODO: добавить сервисы и репозитории
	// AuthService  *service.AuthService
	// UserRepo     repository.UserRepository
	// TokenRepo    repository.TokenRepository
}

// CleanUp очищает ресурсы приложения
func (a *Application) CleanUp(ctx context.Context) error {
	// Закрываем подключения
	if a.DB != nil {
		a.DB.Close()
	}
	if a.Redis != nil {
		a.Redis.Close()
	}
	return nil
}

// ProviderSet набор провайдеров для DI
var ProviderSet = wire.NewSet(
	// Конфигурация
	loadConfig,

	// Подключения
	ProvidePostgresConfig,
	ProvideRedisConfig,
	postgres.NewPool,
	dbredis.NewClient,

	// Логгер
	NewLogger,

	// JWT сервис
	NewJWTService,

	// Cookie сервис
	NewCookieService,

	// TODO: добавить сервисы и репозитории
	// service.NewAuthService,
	// repository.NewUserRepository,
	// repository.NewTokenRepository,

	// Application
	NewApplication,
)

// loadConfig загружает конфигурацию из файла
func loadConfig() (*config.Config, error) {
	return config.Load("config.yaml")
}

// ProvidePostgresConfig предоставляет конфигурацию PostgreSQL
func ProvidePostgresConfig(cfg *config.Config) postgres.Config {
	return postgres.Config{
		DSN:         cfg.Database.URL,
		MaxConns:    int32(cfg.Database.MaxConnections),
		ConnTimeout: time.Duration(cfg.Database.ConnectionTimeout) * time.Second,
	}
}

// ProvideRedisConfig предоставляет конфигурацию Redis
func ProvideRedisConfig(cfg *config.Config) dbredis.Config {
	return dbredis.Config{
		Addr:         cfg.Redis.URL,
		DB:           cfg.Redis.DB,
		PoolSize:     10,
		ConnTimeout:  time.Duration(cfg.Redis.ConnectionTimeout) * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
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

// NewCookieService создаёт cookie сервис из конфигурации
func NewCookieService(cfg *config.Config) *cookies.Service {
	return cookies.NewService(cookies.Config{
		Secure:   cfg.Cookie.Secure,
		HTTPOnly: cfg.Cookie.HTTPOnly,
		SameSite: parseSameSite(cfg.Cookie.SameSite),
		Domain:   cfg.Cookie.Domain,
		Path:     cfg.Cookie.Path,
		MaxAge:   cfg.Cookie.MaxAge,
	})
}

// parseSameSite парсит строку в http.SameSite
func parseSameSite(s string) http.SameSite {
	switch s {
	case "Strict":
		return http.SameSiteStrictMode
	case "None":
		return http.SameSiteNoneMode
	case "Lax", "":
		return http.SameSiteLaxMode
	default:
		return http.SameSiteLaxMode
	}
}

// NewApplication создаёт приложение из зависимостей
func NewApplication(
	cfg *config.Config,
	log *logger.Logger,
	jwtSvc *jwt.Service,
	cookieSvc *cookies.Service,
	db *pgxpool.Pool,
	redisClient *goredis.Client,
) (*Application, error) {
	return &Application{
		Config:        cfg,
		Logger:        log,
		JWTService:    jwtSvc,
		CookieService: cookieSvc,
		DB:            db,
		Redis:         redisClient,
	}, nil
}
