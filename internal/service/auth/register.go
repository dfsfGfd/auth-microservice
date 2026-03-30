package auth

import (
	"context"
	"fmt"
	"strconv"

	"auth-microservice/internal/errors"
	"auth-microservice/internal/model"
	"auth-microservice/pkg/jwt"
)

// RegisterResult содержит аккаунт и пару токенов.
type RegisterResult struct {
	Account *model.Account
	Tokens  *jwt.TokenPair
}

// Register регистрирует нового пользователя и возвращает аккаунт + токены (автовход).
func (s *AuthService) Register(ctx context.Context, email, password string) (*RegisterResult, error) {
	// Валидация email
	emailVO, err := model.NewEmail(email)
	if err != nil {
		return nil, err
	}

	// Валидация пароля
	passwordVO, err := model.NewPlainPassword(password)
	if err != nil {
		return nil, err
	}

	// Проверка: существует ли аккаунт с таким email
	exists, err := s.accountRepo.ExistsByEmail(ctx, emailVO.Value())
	if err != nil {
		s.log.Error("check_email_exists", "err", err)
		return nil, fmt.Errorf("check email exists: %w", err)
	}
	if exists {
		s.log.Warn("register_failed", "reason", "account_exists")
		return nil, errors.ErrAccountExists
	}

	// Хеширование пароля
	hashedPassword, err := s.hasher.Hash(passwordVO.Value(), 0)
	if err != nil {
		s.log.Error("hash_password", "err", err)
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Создание PasswordHash VO
	passwordHash, err := model.NewPasswordHash(hashedPassword)
	if err != nil {
		s.log.Error("create_password_hash", "err", err)
		return nil, fmt.Errorf("create password hash: %w", err)
	}

	// Генерация Snowflake ID
	id, err := s.idGen.Next()
	if err != nil {
		s.log.Error("generate_id", "err", err)
		return nil, fmt.Errorf("generate snowflake id: %w", err)
	}

	// Создание аккаунта с ID
	account, err := model.NewAccount(id, emailVO, passwordHash)
	if err != nil {
		s.log.Error("create_account", "err", err)
		return nil, fmt.Errorf("create account: %w", err)
	}

	// Сохранение аккаунта
	if err := s.accountRepo.Save(ctx, account); err != nil {
		s.log.Error("save_account", "err", err)
		return nil, fmt.Errorf("save account: %w", err)
	}

	// === АВТОХОД: генерация токенов ===
	tokens, err := s.generateTokens(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	s.log.Info("register_with_auto_login", "user_id", id)
	return &RegisterResult{
		Account: account,
		Tokens:  tokens,
	}, nil
}

// generateTokens создаёт пару access + refresh токенов и сохраняет refresh в Redis.
func (s *AuthService) generateTokens(ctx context.Context, account *model.Account) (*jwt.TokenPair, error) {
	accountID := strconv.FormatInt(account.ID(), 10)

	tokens, err := s.jwtService.GenerateTokens(accountID, account.Email().Value())
	if err != nil {
		return nil, fmt.Errorf("generate jwt tokens: %w", err)
	}

	// Сохраняем refresh токен в Redis (для возможности отзыва)
	if err := s.tokenCache.Set(ctx, tokens.RefreshToken, accountID, s.jwtService.RefreshTTLDuration()); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return tokens, nil
}
