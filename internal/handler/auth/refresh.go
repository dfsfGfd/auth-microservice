package auth

import (
	"context"

	stderrors "errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"auth-microservice/internal/errors"
	"auth-microservice/pkg/proto/auth/v1"
)

// Refresh обновляет пару токенов.
func (h *Handler) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	// Вызов сервиса
	tokens, err := h.authService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return h.refreshError(err)
	}

	// Конвертация domain → proto
	return &authv1.RefreshResponse{
		StatusCode: 200,
		Message:    "Token refreshed successfully",
		Data: &authv1.TokenData{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresIn:    int32(tokens.ExpiresIn),
			TokenType:    tokens.TokenType,
		},
	}, nil
}

func (h *Handler) refreshError(err error) (*authv1.RefreshResponse, error) {
	switch {
	case stderrors.Is(err, errors.ErrRefreshTokenNotFound):
		return nil, status.Error(codes.Unauthenticated, "refresh token not found")
	case stderrors.Is(err, errors.ErrTokenExpired):
		return nil, status.Error(codes.Unauthenticated, "token expired")
	case stderrors.Is(err, errors.ErrTokenInvalid):
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	case stderrors.Is(err, errors.ErrAccountNotFound):
		return nil, status.Error(codes.NotFound, "account not found")
	default:
		h.log.Error("refresh failed", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
}
