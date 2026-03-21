package auth

import (
	"context"
	stderrors "errors"

	"auth-microservice/internal/errors"
	"auth-microservice/internal/handler/converter"
	"auth-microservice/pkg/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Register регистрирует нового пользователя.
func (h *Handler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	// Вызов сервиса
	account, err := h.authService.Register(ctx, req.Email, req.Password)
	if err != nil {
		return h.registerError(err)
	}

	// Конвертация domain → proto
	return &authv1.RegisterResponse{
		StatusCode: 200,
		Message:    "Account registered successfully",
		Data:       converter.AccountToProto(account),
	}, nil
}

func (h *Handler) registerError(err error) (*authv1.RegisterResponse, error) {
	switch {
	case stderrors.Is(err, errors.ErrEmailInvalid):
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	case stderrors.Is(err, errors.ErrPasswordInvalid):
		return nil, status.Error(codes.InvalidArgument, "invalid password")
	case stderrors.Is(err, errors.ErrPasswordTooShort):
		return nil, status.Error(codes.InvalidArgument, "password too short")
	case stderrors.Is(err, errors.ErrAccountExists):
		return nil, status.Error(codes.AlreadyExists, "account already exists")
	default:
		h.log.Error("register failed", "error", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
}
