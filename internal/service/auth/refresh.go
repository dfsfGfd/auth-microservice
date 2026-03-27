package auth

import (
	"context"
	stderrors "errors"
	"fmt"

	"auth-microservice/internal/errors"
	"auth-microservice/pkg/jwt"
	"github.com/google/uuid"
)

// Refresh обновляет пару токенов.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	// Валидация refresh токена
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		s.log.Warn("refresh_failed", "reason", "invalid_token", "err", err)
		// Конвертируем ошибки jwt в доменные
		if stderrors.Is(err, jwt.ErrExpiredToken) {
			return nil, errors.ErrTokenExpired
		}
		if stderrors.Is(err, jwt.ErrInvalidToken) {
			return nil, errors.ErrTokenInvalid
		}
		return nil, err
	}

	// Проверка токена в кэше
	accountID, err := s.tokenCache.Get(ctx, refreshToken)
	if err != nil {
		s.log.Warn("refresh_failed", "reason", "token_not_in_cache", "err", err)
		return nil, errors.ErrRefreshTokenNotFound
	}

	// Парсинг UUID
	id, err := uuid.Parse(accountID)
	if err != nil {
		s.log.Error("parse_account_id", "err", err)
		return nil, fmt.Errorf("parse account id: %w", err)
	}

	// Проверка: существует ли аккаунт
	exists, err := s.accountRepo.ExistsByID(ctx, id)
	if err != nil {
		s.log.Error("check_account_exists", "err", err)
		return nil, fmt.Errorf("check account exists: %w", err)
	}
	if !exists {
		s.log.Warn("refresh_failed", "reason", "account_not_found")
		return nil, errors.ErrAccountNotFound
	}

	// Генерация новой пары токенов
	tokens, err := s.jwtService.GenerateTokens(accountID, claims.Subject)
	if err != nil {
		s.log.Error("generate_tokens", "err", err)
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	// Обновление токена в кэше (сброс TTL)
	if err := s.tokenCache.Set(ctx, tokens.RefreshToken, accountID, s.jwtService.RefreshTTLDuration()); err != nil {
		s.log.Error("cache_refresh_token", "err", err)
		return nil, fmt.Errorf("cache refresh token: %w", err)
	}

	// Удаление старого токена из кэша (ошибка не критична)
	if err := s.tokenCache.Delete(ctx, refreshToken); err != nil {
		s.log.Warn("delete_old_token", "err", err)
	}

	s.log.Info("refresh_token", "user_id", accountID)
	return tokens, nil
}
