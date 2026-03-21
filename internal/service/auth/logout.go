package auth

import (
	"context"
	"fmt"

	stderrors "errors"

	"auth-microservice/internal/errors"
	"auth-microservice/pkg/jwt"
)

// Logout выполняет выход (отзыв refresh токена).
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	// Валидация refresh токена
	_, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		s.log.Warn("logout_failed", "reason", "invalid_token", "err", err)
		// Конвертируем ошибки jwt в доменные
		if stderrors.Is(err, jwt.ErrExpiredToken) {
			return errors.ErrTokenExpired
		}
		if stderrors.Is(err, jwt.ErrInvalidToken) {
			return errors.ErrTokenInvalid
		}
		return err
	}

	// Удаление токена из кэша
	if err := s.tokenCache.Delete(ctx, refreshToken); err != nil {
		s.log.Error("delete_refresh_token", "err", err)
		return fmt.Errorf("delete refresh token: %w", err)
	}

	s.log.Info("logout")
	return nil
}
