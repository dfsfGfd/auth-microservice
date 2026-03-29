package auth

import (
	"context"
	stderrors "errors"
	"strconv"

	"auth-microservice/internal/errors"
	"auth-microservice/pkg/proto/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Register регистрирует нового пользователя и возвращает токены (автовход).
func (h *Handler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	// Вызов сервиса — теперь возвращает и аккаунт, и токены
	account, tokens, err := h.authService.Register(ctx, req.Email, req.Password)
	if err != nil {
		return h.registerError(err)
	}

	// Конвертация domain → proto
	return &authv1.RegisterResponse{
		StatusCode: 200,
		Message:    "Account registered successfully",
		Data: &authv1.RegisterData{
			AccountId:    strconv.FormatInt(account.ID(), 10),
			Email:        account.Email().String(),
			CreatedAt:    timestamppb.New(account.CreatedAt()),
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresIn:    int32(tokens.ExpiresIn),
			TokenType:    tokens.TokenType,
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
