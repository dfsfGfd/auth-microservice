package auth

import (
	"context"

	stderrors "errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"auth-microservice/internal/errors"
	"auth-microservice/pkg/proto/auth/v1"
)

// Login выполняет вход и возвращает пару токенов.
func (h *Handler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	// Вызов сервиса
	tokens, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return h.loginError(err)
	}

	// Конвертация domain → proto
	return &authv1.LoginResponse{
		StatusCode: 200,
		Message:    "Login successful",
		Data: &authv1.TokenData{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresIn:    int32(tokens.ExpiresIn),
			TokenType:    tokens.TokenType,
		},
	}, nil
}

func (h *Handler) loginError(err error) (*authv1.LoginResponse, error) {
	// Логируем ошибку для внутреннего использования
	if stderrors.Is(err, errors.ErrAccountNotFound) ||
		stderrors.Is(err, errors.ErrInvalidCredentials) ||
		stderrors.Is(err, errors.ErrEmailInvalid) {
		h.log.Warn("login failed", "error", err)
	} else {
		h.log.Error("login failed", "error", err)
	}

	// Всегда возвращаем одинаковую ошибку для предотвращения user enumeration
	// Злоумышленник не должен узнавать, существует ли аккаунт с таким email
	return nil, status.Error(codes.Unauthenticated, "invalid credentials")
}
