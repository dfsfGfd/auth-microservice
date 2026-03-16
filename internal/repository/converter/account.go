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
	email, err := model.NewEmail(db.Email)
	if err != nil {
		return nil, err
	}

	passwordHash := model.NewPasswordHashFromString(db.PasswordHash)

	// Создаём агрегат
	return model.NewAccount(email, passwordHash)
}

// AccountListToDB конвертирует список domain Account в список DB моделей.
func AccountListToDB(accounts []*model.Account) []*dbmodel.Account {
	if accounts == nil {
		return nil
	}

	result := make([]*dbmodel.Account, 0, len(accounts))
	for _, account := range accounts {
		result = append(result, AccountToDB(account))
	}
	return result
}

// AccountListToDomain конвертирует список DB моделей в список domain Account.
func AccountListToDomain(dbAccounts []*dbmodel.Account) ([]*model.Account, error) {
	if dbAccounts == nil {
		return nil, nil
	}

	result := make([]*model.Account, 0, len(dbAccounts))
	for _, db := range dbAccounts {
		account, err := AccountToDomain(db)
		if err != nil {
			return nil, err
		}
		result = append(result, account)
	}
	return result, nil
}
