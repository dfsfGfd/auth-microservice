package auth

import (
	"context"
	"fmt"
	"strconv"

	"auth-microservice/internal/errors"
	"auth-microservice/internal/model"
	"auth-microservice/pkg/jwt"
)

// Login выполняет вход и возвращает пару токенов.
func (s *AuthService) Login(ctx context.Context, email, password string) (*jwt.TokenPair, error) {
	// Валидация email
	emailVO, err := model.NewEmail(email)
	if err != nil {
		// Возвращаем ErrInvalidCredentials вместо ErrEmailInvalid для предотвращения user enumeration
		s.log.Warn("login_failed", "reason", "invalid_email")
		return nil, errors.ErrInvalidCredentials
	}

	// Поиск аккаунта по email
	account, err := s.accountRepo.GetByEmail(ctx, emailVO.Value())
	if err != nil {
		// Если аккаунт не найден, возвращаем ErrInvalidCredentials вместо ErrAccountNotFound
		// Это предотвращает user enumeration - злоумышленник не узнает, существует ли email
		s.log.Debug("login_failed", "reason", "account_not_found")
		return nil, errors.ErrInvalidCredentials
	}

	// Сравнение пароля
	if cmpErr := s.hasher.Compare(account.PasswordHash().Value(), password); cmpErr != nil {
		s.log.Warn("login_failed", "reason", "invalid_password")
		return nil, errors.ErrInvalidCredentials
	}

	// Генерация JWT токенов
	accountID := strconv.FormatInt(account.ID(), 10)
	tokens, err := s.jwtService.GenerateTokens(accountID, account.Email().Value())
	if err != nil {
		s.log.Error("generate_tokens", "err", err)
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	// Сохранение refresh токена в кэш
	if err := s.tokenCache.Set(ctx, tokens.RefreshToken, accountID, s.jwtService.RefreshTTLDuration()); err != nil {
		s.log.Error("cache_refresh_token", "err", err)
		return nil, fmt.Errorf("cache refresh token: %w", err)
	}

	s.log.Info("login", "user_id", account.ID())
	return tokens, nil
}
