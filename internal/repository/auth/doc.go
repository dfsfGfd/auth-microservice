// Package auth предоставляет PostgreSQL реализацию репозиториев.
//
// Реализация интерфейсов:
//   - repository.UserRepository
//
// Структура пакета:
//
//	auth/
//	├── user_repository.go      # UserRepository реализация
//	├── transaction.go          # TransactionManager реализация
//	└── mapper.go               # Вспомогательные функции маппинга
package auth
