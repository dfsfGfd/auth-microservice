// Package converter предоставляет конвертеры между proto и domain моделями.
package converter

import (
	"auth-microservice/internal/model"
	"auth-microservice/pkg/proto/auth/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AccountToProto конвертирует domain Account в proto RegisterData.
func AccountToProto(account *model.Account) *authv1.RegisterData {
	if account == nil {
		return nil
	}

	return &authv1.RegisterData{
		AccountId: account.ID().String(),
		Email:     account.Email().String(),
		CreatedAt: timestamppb.New(account.CreatedAt()),
	}
}
