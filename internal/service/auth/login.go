package auth

import (
	"context"

	"auth-microservice/internal/errors"
	"auth-microservice/internal/model"
	"auth-microservice/pkg/jwt"
)

// Login выполняет вход и возвращает пару токенов.
func (s *AuthService) Login(ctx context.Context, email, password string) (*jwt.TokenPair, error) {
	// Валидация email
	emailVO, err := model.NewEmail(email)
	if err != nil {
		return nil, err
	}

	// Поиск аккаунта по email
	account, err := s.accountRepo.GetByEmail(ctx, emailVO.Value())
	if err != nil {
		s.log.Error("get account by email", "email", email, "error", err)
		return nil, err
	}

	// Сравнение пароля
	if err := s.hasher.Compare(account.PasswordHash().Value(), password); err != nil {
		s.log.Warn("invalid password", "email", email)
		return nil, errors.ErrInvalidCredentials
	}

	// Генерация JWT токенов
	tokens, err := s.jwtService.GenerateTokens(account.ID().String(), account.Email().Value())
	if err != nil {
		s.log.Error("generate tokens", "error", err)
		return nil, err
	}

	// Сохранение refresh токена в кэш
	refreshTTL, err := s.jwtService.RefreshTTLDuration()
	if err != nil {
		s.log.Error("get refresh ttl", "error", err)
		return nil, err
	}

	if err := s.tokenCache.Set(ctx, tokens.RefreshToken, account.ID().String(), refreshTTL); err != nil {
		s.log.Error("cache refresh token", "error", err)
		return nil, err
	}

	s.log.Info("account logged in", "account_id", account.ID(), "email", email)
	return tokens, nil
}
