package auth

import (
	"context"
	"strconv"

	stderrors "errors"

	"auth-microservice/internal/errors"
	"auth-microservice/pkg/proto/auth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Register регистрирует нового пользователя и возвращает токены (автовход).
func (h *Handler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	// Вызов сервиса — возвращает аккаунт + токены
	result, err := h.authService.Register(ctx, req.Email, req.Password)
	if err != nil {
		return h.registerError(err)
	}

	// Конвертация domain → proto
	return &authv1.RegisterResponse{
		StatusCode: 200,
		Message:    "Account registered successfully",
		Data: &authv1.RegisterData{
			AccountId:    strconv.FormatInt(result.Account.ID(), 10),
			Email:        result.Account.Email().Value(),
			CreatedAt:    timestamppb.New(result.Account.CreatedAt()),
			AccessToken:  result.Tokens.AccessToken,
			RefreshToken: result.Tokens.RefreshToken,
			ExpiresIn:    int32(result.Tokens.ExpiresIn),
			TokenType:    result.Tokens.TokenType,
		},
	}, nil
}

func (h *Handler) registerError(err error) (*authv1.RegisterResponse, error) {
	switch {
	case stderrors.Is(err, errors.ErrEmailInvalid):
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	case stderrors.Is(err, errors.ErrEmailTooLong):
		return nil, status.Error(codes.InvalidArgument, "email too long")
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
