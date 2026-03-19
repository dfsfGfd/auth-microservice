// Package auth предоставляет реализацию gRPC хендлера для сервиса аутентификации.
package auth

import (
	svc "auth-microservice/internal/service/auth"
	"auth-microservice/pkg/logger"
	"auth-microservice/pkg/proto/auth/v1"
)

// Handler реализация handler.AuthHandler.
type Handler struct {
	authv1.UnimplementedAuthServiceServer

	authService *svc.AuthService
	log         *logger.Logger
}

// NewHandler создаёт новый gRPC хендлер.
func NewHandler(
	authService *svc.AuthService,
	log *logger.Logger,
) *Handler {
	return &Handler{
		authService: authService,
		log:         log,
	}
}
