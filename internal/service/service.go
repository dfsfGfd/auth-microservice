// Package service предоставляет интерфейсы для бизнес-логики.
package service

import (
	"context"

	"auth-microservice/internal/model"
	"auth-microservice/pkg/jwt"
)

// AuthService интерфейс сервиса аутентификации.
type AuthService interface {
	// Register регистрирует нового пользователя.
	Register(ctx context.Context, email, password string) (*model.Account, error)

	// Login выполняет вход и возвращает пару токенов.
	Login(ctx context.Context, email, password string) (*jwt.TokenPair, error)

	// Logout выполняет выход (отзыв refresh токена).
	Logout(ctx context.Context, refreshToken string) error

	// Refresh обновляет пару токенов.
	Refresh(ctx context.Context, refreshToken string) (*jwt.TokenPair, error)
}
