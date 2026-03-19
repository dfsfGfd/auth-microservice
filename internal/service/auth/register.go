package auth

import (
	"context"
	"fmt"

	"auth-microservice/internal/errors"
	"auth-microservice/internal/model"
)

// Register регистрирует нового пользователя.
func (s *AuthService) Register(ctx context.Context, email, password string) (*model.Account, error) {
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
		s.log.Error("check email exists", "email", email, "error", err)
		return nil, fmt.Errorf("check email exists: %w", err)
	}
	if exists {
		return nil, errors.ErrAccountExists
	}

	// Хеширование пароля
	hashedPassword, err := s.hasher.Hash(passwordVO.Value(), 0)
	if err != nil {
		s.log.Error("hash password", "error", err)
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Создание PasswordHash VO
	passwordHash, err := model.NewPasswordHash(hashedPassword)
	if err != nil {
		s.log.Error("create password hash", "error", err)
		return nil, fmt.Errorf("create password hash: %w", err)
	}

	// Создание аккаунта
	account, err := model.NewAccount(emailVO, passwordHash)
	if err != nil {
		s.log.Error("create account", "error", err)
		return nil, fmt.Errorf("create account: %w", err)
	}

	// Сохранение аккаунта
	if err := s.accountRepo.Save(ctx, account); err != nil {
		s.log.Error("save account", "email", email, "error", err)
		return nil, fmt.Errorf("save account: %w", err)
	}

	s.log.Info("account registered", "account_id", account.ID(), "email", email)
	return account, nil
}
