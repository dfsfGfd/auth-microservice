// Package converter предоставляет конвертеры между domain моделями и DB моделями.
package converter

import (
	"auth-microservice/internal/model"
	dbmodel "auth-microservice/internal/repository/model"
)

// UserToDB конвертирует domain User в DB модель.
func UserToDB(user *model.User) *dbmodel.User {
	if user == nil {
		return nil
	}

	return &dbmodel.User{
		ID:           user.ID(),
		Email:        user.Email().String(),
		Username:     user.Username().String(),
		PasswordHash: user.PasswordHash().Value(),
		CreatedAt:    user.CreatedAt(),
		UpdatedAt:    user.UpdatedAt(),
	}
}

// UserToDomain конвертирует DB модель в domain User.
func UserToDomain(db *dbmodel.User) (*model.User, error) {
	if db == nil {
		return nil, nil
	}

	// Создаём Value Objects
	email, err := model.NewEmail(db.Email)
	if err != nil {
		return nil, err
	}

	username, err := model.NewUsername(db.Username)
	if err != nil {
		return nil, err
	}

	passwordHash := model.NewPasswordHashFromString(db.PasswordHash)

	// Создаём агрегат
	return model.NewUser(email, username, passwordHash)
}

// UserListToDB конвертирует список domain User в список DB моделей.
func UserListToDB(users []*model.User) []*dbmodel.User {
	if users == nil {
		return nil
	}

	result := make([]*dbmodel.User, 0, len(users))
	for _, user := range users {
		result = append(result, UserToDB(user))
	}
	return result
}

// UserListToDomain конвертирует список DB моделей в список domain User.
func UserListToDomain(dbUsers []*dbmodel.User) ([]*model.User, error) {
	if dbUsers == nil {
		return nil, nil
	}

	result := make([]*model.User, 0, len(dbUsers))
	for _, db := range dbUsers {
		user, err := UserToDomain(db)
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	return result, nil
}
