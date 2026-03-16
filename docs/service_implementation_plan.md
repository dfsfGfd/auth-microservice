# План реализации Service Layer

## Обзор

Service layer содержит бизнес-логику приложения и связывает все слои:
- **Repository** — доступ к данным (PostgreSQL)
- **Cache** — кэш токенов (Redis)
- **JWT** — генерация/валидация токенов
- **Bcrypt** — хеширование паролей

---

## Структура пакета

```
internal/service/
├── service.go              # Интерфейс AuthService
├── auth_service.go         # Реализация AuthService
├── auth_service_test.go    # Unit-тесты
└── errors.go               # Специфичные ошибки сервиса
```

---

## 1. Интерфейс AuthService (`service.go`)

**Методы:**

| Метод | Сигнатура | Описание |
|-------|-----------|----------|
| `Register` | `Register(ctx, email, password) (*Account, error)` | Регистрация нового аккаунта |
| `Login` | `Login(ctx, email, password) (*TokenPair, error)` | Вход (получение токенов) |
| `Logout` | `Logout(ctx, refreshToken) error` | Выход (отзыв refresh токена) |
| `Refresh` | `Refresh(ctx, refreshToken) (*TokenPair, error)` | Обновление пары токенов |

---

## 2. Реализация (`auth_service.go`)

### Зависимости сервиса

```go
type AuthService struct {
    accountRepo repository.AccountRepository
    tokenCache  *token.RedisCache
    jwtService  *jwt.Service
    hasher      *bcrypt.Service
    logger      *logger.Logger
}
```

### Метод 1: Register

**Логика:**
1. Валидация email (через `model.NewEmail`)
2. Валидация пароля (через `model.NewPlainPassword`)
3. Проверка: не существует ли аккаунт с таким email (`ExistsByEmail`)
4. Хеширование пароля (`hasher.Hash`)
5. Создание PasswordHash VO (`model.NewPasswordHash`)
6. Создание Account (`model.NewAccount`)
7. Сохранение (`accountRepo.Save`)
8. Возврат Account

**Ошибки:**
- `ErrEmailInvalid` — невалидный email
- `ErrPasswordInvalid` — невалидный пароль
- `ErrAccountExists` — аккаунт уже существует
- `ErrRepository` — ошибка сохранения

---

### Метод 2: Login

**Логика:**
1. Валидация email
2. Поиск аккаунта по email (`GetByEmail`)
3. Сравнение пароля (`hasher.Compare`)
4. Генерация JWT токенов (`jwtService.GenerateTokens`)
5. Сохранение refresh токена в кэш (`tokenCache.Set`)
6. Возврат TokenPair

**Ошибки:**
- `ErrEmailInvalid` — невалидный email
- `ErrAccountNotFound` — аккаунт не найден
- `ErrInvalidCredentials` — неверный пароль
- `ErrInternal` — ошибка генерации токенов

---

### Метод 3: Logout

**Логика:**
1. Валидация refresh токена (`jwtService.ValidateRefreshToken`)
2. Удаление токена из кэша (`tokenCache.Delete`)
3. Возврат nil

**Ошибки:**
- `ErrTokenInvalid` — невалидный токен
- `ErrTokenExpired` — токен истёк

---

### Метод 4: Refresh

**Логика:**
1. Валидация refresh токена (`jwtService.ValidateRefreshToken`)
2. Проверка токена в кэше (`tokenCache.Get`)
3. Если не найден — ошибка (токен отозван)
4. Получение account_id из токена
5. Генерация новой пары токенов
6. Обновление токена в кэше (TTL сбрасывается)
7. Возврат новой TokenPair

**Ошибки:**
- `ErrTokenInvalid` — невалидный токен
- `ErrTokenExpired` — токен истёк
- `ErrRefreshTokenNotFound` — токен не найден в кэше (отозван)
- `ErrAccountNotFound` — аккаунт не найден

---

## 3. Ошибки сервиса (`errors.go`)

```go
var (
    ErrRegisterFailed      = errors.New("registration failed")
    ErrLoginFailed         = errors.New("login failed")
    ErrLogoutFailed        = errors.New("logout failed")
    ErrRefreshFailed       = errors.New("refresh failed")
    ErrPasswordMismatch    = errors.New("password does not match")
    ErrTokenRevoked        = errors.New("token has been revoked")
)
```

---

## 4. DI интеграция (`internal/di/provider.go`)

Добавить провайдеры:

```go
// В ProviderSet:
service.NewAuthService,

// Provide функции:
func ProvideAuthService(
    accountRepo repository.AccountRepository,
    tokenCache *token.RedisCache,
    jwtService *jwt.Service,
    hasher *bcrypt.Service,
    log *logger.Logger,
) *service.AuthService {
    return service.NewAuthService(accountRepo, tokenCache, jwtService, hasher, log)
}
```

Обновить `Application`:
```go
type Application struct {
    ...
    AuthService *service.AuthService
}
```

---

## 5. Тесты (`auth_service_test.go`)

### Тесты Register:
- ✅ Успешная регистрация
- ❌ Невалидный email
- ❌ Невалидный пароль
- ❌ Аккаунт уже существует

### Тесты Login:
- ✅ Успешный вход
- ❌ Аккаунт не найден
- ❌ Неверный пароль

### Тесты Logout:
- ✅ Успешный logout
- ❌ Невалидный токен

### Тесты Refresh:
- ✅ Успешный refresh
- ❌ Токен истёк
- ❌ Токен отозван (нет в кэше)

---

## Порядок реализации

### Этап 1: Базовая структура
1. Создать `internal/service/service.go` (интерфейс)
2. Создать `internal/service/errors.go` (ошибки)

### Этап 2: Реализация
3. Создать `internal/service/auth_service.go`
   - Конструктор `NewAuthService`
   - Метод `Register`
   - Метод `Login`
   - Метод `Logout`
   - Метод `Refresh`

### Этап 3: DI интеграция
4. Добавить `ProvideAuthService` в `internal/di/provider.go`
5. Добавить `AuthService` в `Application`
6. Перегенерировать `wire_gen.go`

### Этап 4: Тесты
7. Создать `internal/service/auth_service_test.go`
8. Написать тесты для всех методов

### Этап 5: Коммит
9. Закоммитить изменения

---

## Связи с другими слоями

```
┌─────────────────────────────────────────────────────────┐
│                    Handler (gRPC)                        │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                   AuthService                            │
├──────────────┬──────────────┬──────────────┬────────────┤
│              │              │              │            │
▼              ▼              ▼              ▼            ▼
AccountRepo  TokenCache   JWT Service   Hasher      Logger
(PostgreSQL)  (Redis)      (pkg/jwt)   (pkg/bcrypt)
```

---

## Критерии готовности

- [ ] Все 4 метода реализованы
- [ ] DI настроен через Wire
- [ ] Тесты покрывают основные сценарии
- [ ] Сборка проходит без ошибок
- [ ] Логирование добавлено в ключевых точках
