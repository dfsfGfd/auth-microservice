//go:build wireinject
// +build wireinject

package di

import "github.com/google/wire"

// InitializeApplication инициализирует приложение с внедрением зависимостей.
// Файл генерируется автоматически через: go generate ./...
func InitializeApplication() (*Application, error) {
	wire.Build(ProviderSet)
	return &Application{}, nil
}
