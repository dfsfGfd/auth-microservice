// Package converter предоставляет конвертеры между domain моделями и DB моделями.
package converter

import (
	"auth-microservice/internal/model"
	dbmodel "auth-microservice/internal/repository/model"
)

// AccountToDB конвертирует domain Account в DB модель.
func AccountToDB(account *model.Account) *dbmodel.Account {
	if account == nil {
		return nil
	}

	return &dbmodel.Account{
		ID:           account.ID(),
		Email:        account.Email().String(),
		PasswordHash: account.PasswordHash().Value(),
		CreatedAt:    account.CreatedAt(),
		UpdatedAt:    account.UpdatedAt(),
	}
}

// AccountToDomain конвертирует DB модель в domain Account.
func AccountToDomain(db *dbmodel.Account) (*model.Account, error) {
	if db == nil {
		return nil, nil
	}

	// Создаём Value Objects
	// Email уже валидирован при сохранении, используем доверенный конструктор
	email := model.NewEmailFromDB(db.Email)

	passwordHash := model.NewPasswordHashFromString(db.PasswordHash)

	// Создаём агрегат с установленным ID
	account, err := model.NewAccount(email, passwordHash)
	if err != nil {
		return nil, err
	}

	// Устанавливаем ID и временные метки из БД
	account.SetID(db.ID)
	account.SetCreatedAt(db.CreatedAt)
	account.SetUpdatedAt(db.UpdatedAt)

	return account, nil
}
