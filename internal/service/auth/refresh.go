package auth

import (
	"context"

	"github.com/google/uuid"

	"auth-microservice/internal/errors"
	"auth-microservice/pkg/jwt"
)

// Refresh обновляет пару токенов.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	// Валидация refresh токена
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		s.log.Warn("invalid refresh token", "error", err)
		return nil, err
	}

	// Проверка токена в кэше
	accountID, err := s.tokenCache.Get(ctx, refreshToken)
	if err != nil {
		s.log.Warn("refresh token not found in cache", "error", err)
		return nil, errors.ErrRefreshTokenNotFound
	}

	// Парсинг UUID
	id, err := uuid.Parse(accountID)
	if err != nil {
		s.log.Error("parse account id", "account_id", accountID, "error", err)
		return nil, err
	}

	// Проверка: существует ли аккаунт
	exists, err := s.accountRepo.ExistsByID(ctx, id)
	if err != nil {
		s.log.Error("check account exists", "account_id", accountID, "error", err)
		return nil, err
	}
	if !exists {
		s.log.Warn("account not found", "account_id", accountID)
		return nil, errors.ErrAccountNotFound
	}

	// Генерация новой пары токенов
	tokens, err := s.jwtService.GenerateTokens(accountID, claims.Email)
	if err != nil {
		s.log.Error("generate tokens", "error", err)
		return nil, err
	}

	// Обновление токена в кэше (сброс TTL)
	refreshTTL, err := s.jwtService.RefreshTTLDuration()
	if err != nil {
		s.log.Error("get refresh ttl", "error", err)
		return nil, err
	}

	if err := s.tokenCache.Set(ctx, tokens.RefreshToken, accountID, refreshTTL); err != nil {
		s.log.Error("cache refresh token", "error", err)
		return nil, err
	}

	// Удаление старого токена из кэша
	_ = s.tokenCache.Delete(ctx, refreshToken)

	s.log.Info("tokens refreshed", "account_id", accountID)
	return tokens, nil
}
