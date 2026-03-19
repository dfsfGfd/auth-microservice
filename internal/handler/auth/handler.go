// Package auth предоставляет реализацию gRPC хендлера для сервиса аутентификации.
package auth

import (
	svc "auth-microservice/internal/service/auth"
	"auth-microservice/pkg/cookies"
	"auth-microservice/pkg/logger"
	"auth-microservice/pkg/proto/auth/v1"
)

// Handler реализация handler.AuthHandler.
type Handler struct {
	authv1.UnimplementedAuthServiceServer

	authService   *svc.AuthService
	cookieService *cookies.Service
	log           *logger.Logger
}

// NewHandler создаёт новый gRPC хендлер.
func NewHandler(
	authService *svc.AuthService,
	cookieService *cookies.Service,
	log *logger.Logger,
) *Handler {
	return &Handler{
		authService:   authService,
		cookieService: cookieService,
		log:           log,
	}
}
