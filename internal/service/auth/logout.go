package auth

import (
	"context"
)

// Logout выполняет выход (отзыв refresh токена).
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	// Валидация refresh токена
	_, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		s.log.Warn("invalid refresh token", "error", err)
		return err
	}

	// Удаление токена из кэша
	if err := s.tokenCache.Delete(ctx, refreshToken); err != nil {
		s.log.Error("delete refresh token", "error", err)
		return err
	}

	s.log.Info("account logged out")
	return nil
}
