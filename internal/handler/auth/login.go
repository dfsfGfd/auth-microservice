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
	switch {
	case stderrors.Is(err, errors.ErrEmailInvalid):
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	case stderrors.Is(err, errors.ErrAccountNotFound):
		return nil, status.Error(codes.NotFound, "account not found")
	case stderrors.Is(err, errors.ErrInvalidCredentials):
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	default:
		h.log.Error("login failed", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
}
