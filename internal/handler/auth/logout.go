package auth

import (
	"context"

	stderrors "errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"auth-microservice/internal/errors"
	"auth-microservice/pkg/proto/auth/v1"
)

// Logout выполняет выход (отзыв refresh токена).
func (h *Handler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	// Вызов сервиса
	err := h.authService.Logout(ctx, req.RefreshToken)
	if err != nil {
		return h.logoutError(err)
	}

	// Конвертация domain → proto
	return &authv1.LogoutResponse{
		StatusCode: 200,
		Message:    "Logout successful",
		Data: &authv1.LogoutData{
			Success: true,
		},
	}, nil
}

func (h *Handler) logoutError(err error) (*authv1.LogoutResponse, error) {
	switch {
	case stderrors.Is(err, errors.ErrTokenExpired):
		return nil, status.Error(codes.Unauthenticated, "token expired")
	case stderrors.Is(err, errors.ErrTokenInvalid):
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	default:
		h.log.Error("logout failed", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
}
