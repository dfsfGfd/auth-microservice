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

	// Создание аккаунта
	account, err := model.NewAccount(emailVO, passwordHash)
	if err != nil {
		s.log.Error("create_account", "err", err)
		return nil, fmt.Errorf("create account: %w", err)
	}

	// Генерация Snowflake ID
	id, err := s.idGen.Next()
	if err != nil {
		s.log.Error("generate_id", "err", err)
		return nil, fmt.Errorf("generate snowflake id: %w", err)
	}
	account.SetID(id)

	// Сохранение аккаунта
	if err := s.accountRepo.Save(ctx, account); err != nil {
		s.log.Error("save_account", "err", err)
		return nil, fmt.Errorf("save account: %w", err)
	}

	s.log.Info("register", "user_id", id)
	return account, nil
}
