// Package auth предоставляет реализацию сервиса аутентификации.
package auth

import (
	"auth-microservice/internal/cache/token"
	"auth-microservice/internal/repository"
	"auth-microservice/pkg/bcrypt"
	"auth-microservice/pkg/jwt"
	"auth-microservice/pkg/logger"
)

// AuthService реализация service.AuthService.
type AuthService struct {
	accountRepo repository.AccountRepository
	tokenCache  *token.RedisCache
	jwtService  *jwt.Service
	hasher      *bcrypt.Service
	log         *logger.Logger
}

// NewAuthService создаёт новый AuthService.
func NewAuthService(
	accountRepo repository.AccountRepository,
	tokenCache *token.RedisCache,
	jwtService *jwt.Service,
	hasher *bcrypt.Service,
	log *logger.Logger,
) *AuthService {
	return &AuthService{
		accountRepo: accountRepo,
		tokenCache:  tokenCache,
		jwtService:  jwtService,
		hasher:      hasher,
		log:         log,
	}
}
