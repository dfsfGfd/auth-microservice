## Qwen Added Memories
- Пользователь работает над Go-микросервисом аутентификации (auth-microservice) с использованием gRPC + REST (grpc-gateway), JWT access/refresh токенов, PostgreSQL (pgx) и Redis. Проект использует DDD с Value Objects (Email, Username, Password, PasswordHash) и агрегатом User. GitHub username пользователя: dfsfGfd.
- Пользователь работает над Go-микросервисом аутентификации (auth-microservice). Текущее состояние проекта:

**Реализовано:**
1. **Domain layer** (`internal/model/`): User агрегат, Value Objects (Email, Username, Password, PasswordHash)
2. **Repository layer** (`internal/repository/`):
   - Интерфейс: `repository.go` — 6 методов (Save, DeleteByID, GetByID, GetByEmail, GetByUsername, GetAll)
   - Реализация (`auth/`): repository.go (конструктор), save.go, delete_by_id.go, get_by_id.go, get_by_email.go, get_by_username.go, get_all.go
   - DB модель: `model/user.go`
   - Конвертер: `converter/user.go` (UserToDB, UserToDomain)
   - Ошибки: `internal/errors/` (ErrUserNotFound, ErrRepository, ErrDBConnection, ErrDBQuery)
3. **DI** (`internal/di/`): Google Wire с UserRepository
4. **Config** (`internal/config/`): YAML конфигурация
5. **DB подключения** (`pkg/db/`): postgresql/, redisdb/
6. **Пакеты** (`pkg/`): logger, jwt, cookies, proto
7. **Proto**: gRPC + REST (grpc-gateway), Swagger
8. **Документация** (`docs/`): README.md, api.md, config.md, repository_methods.md

**Следующие шаги:**
- Реализация service layer (`internal/service/`)
- Реализация handler layer (`internal/handler/`)
- Создание cmd/server/main.go
- Интеграция всех слоёв

**Важные правила проекта:**
- Один файл = один метод в repository/auth/
- Кастомные ошибки в `internal/errors/`
- Конвертеры domain ↔ DB в `internal/repository/converter/`
- DB модели в `internal/repository/model/`
- Интерфейсы в `internal/repository/repository.go`
- Реализация в `internal/repository/auth/`
- Проект: Go auth-microservice (gRPC + REST). Реализованы ВСЕ слои: Domain (model), Repository (PostgreSQL), Service (бизнес-логика), Handler (gRPC), Cache (Redis для токенов), DI (Google Wire), cmd/server/main.go (запуск серверов). Ошибки: конвертация jwt ошибок в доменные через errors.Is(). Структура: internal/{model,repository,service,handler,cache,di,errors,config}, pkg/{proto,jwt,bcrypt,logger,cookies,db}, cmd/server/main.go. Коммиты: последние - fix: handle errors properly with errors.Is() (a0f3c65). Следующие шаги: написание интеграционных тестов, Dockerfile, docker-compose.yml.
