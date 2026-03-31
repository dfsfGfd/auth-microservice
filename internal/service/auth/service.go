// Package auth предоставляет реализацию сервиса аутентификации.
package auth

import (
	"auth-microservice/internal/cache/token"
	"auth-microservice/internal/repository"
	"auth-microservice/pkg/bcrypt"
	"auth-microservice/pkg/jwt"
	"auth-microservice/pkg/logger"
	"auth-microservice/pkg/snowflake"
)

// AuthService реализация service.AuthService.
type AuthService struct {
	accountRepo repository.AccountRepository
	tokenCache  token.TokenCache
	jwtService  *jwt.Service
	hasher      *bcrypt.Service
	log         *logger.Logger
	idGen       *snowflake.Generator
}

// NewAuthService создаёт новый AuthService.
func NewAuthService(
	accountRepo repository.AccountRepository,
	tokenCache token.TokenCache,
	jwtService *jwt.Service,
	hasher *bcrypt.Service,
	log *logger.Logger,
	idGen *snowflake.Generator,
) *AuthService {
	return &AuthService{
		accountRepo: accountRepo,
		tokenCache:  tokenCache,
		jwtService:  jwtService,
		hasher:      hasher,
		log:         log,
		idGen:       idGen,
	}
}
