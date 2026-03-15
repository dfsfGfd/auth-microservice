// Package auth предоставляет PostgreSQL реализацию репозиториев.
//
// Реализация интерфейсов:
//   - repository.UserRepository
//   - repository.TokenRepository (Redis)
//
// Структура пакета:
//
//	auth/
//	├── user_repository.go      # UserRepository реализация
//	├── token_repository.go     # TokenRepository реализация (Redis)
//	├── transaction.go          # TransactionManager реализация
//	└── mapper.go               # Вспомогательные функции маппинга
package auth
