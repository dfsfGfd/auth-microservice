// Package handler предоставляет gRPC хендлеры для сервиса аутентификации.
package handler

import (
	"context"

	"auth-microservice/pkg/proto/auth/v1"
)

// AuthHandler интерфейс для обработки gRPC запросов.
type AuthHandler interface {
	authv1.AuthServiceServer

	// Register регистрирует нового пользователя.
	Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error)

	// Login выполняет вход и возвращает пару токенов.
	Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error)

	// Logout выполняет выход (отзыв refresh токена).
	Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error)

	// Refresh обновляет пару токенов.
	Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error)
}
