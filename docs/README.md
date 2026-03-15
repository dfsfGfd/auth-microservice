# Auth Microservice

> Микросервис аутентификации с поддержкой gRPC и REST API, использующий JWT access/refresh токены.

---

## 📋 Оглавление

- [Возможности](#-возможности)
- [Архитектура](#-архитектура)
- [Структура проекта](#-структура-проекта)
- [Быстрый старт](#-быстрый-старт)
- [API Endpoints](#-api-endpoints)
- [Разработка](#-разработка)
- [Безопасность](#-безопасность)
- [Зависимости](#-зависимости)

---

## 📋 Возможности

| Функция | Описание |
|---------|----------|
| ✅ **Регистрация** | Создание нового пользователя |
| ✅ **Вход/Выход** | Аутентификация и завершение сессии |
| ✅ **Обновление токенов** | Ротация JWT access/refresh токенов |
| ✅ **gRPC + REST** | Единый сервис для обоих протоколов (grpc-gateway) |
| ✅ **Swagger/OpenAPI** | Автогенерируемая документация API |

---

## 🏗 Архитектура

### Стек технологий

| Компонент | Технология |
|-----------|------------|
| **RPC** | gRPC + REST (grpc-gateway) |
| **Токены** | JWT (access + refresh) |
| **Кэш** | Redis |
| **БД** | PostgreSQL (pgx) |
| **Протокол** | Protocol Buffers v3 |
| **UUID** | Генерация на бэкенде (Go) |

### Время жизни токенов

| Токен | TTL | Хранение |
|-------|-----|----------|
| **Access Token** | 15 минут | Клиент (Authorization header) |
| **Refresh Token** | 2 недели | Redis (`refresh:{token}` → `user_id`) |

---

## 📁 Структура проекта

```
.
├── cmd/
│   └── server/                 # Точка входа приложения
│
├── internal/
│   ├── model/                  # Доменные модели (агрегаты, VO)
│   │   ├── user.go             # User агрегат
│   │   ├── email.go            # Email VO
│   │   ├── username.go         # Username VO
│   │   ├── password.go         # PlainPassword VO
│   │   └── password_hash.go    # PasswordHash VO
│   │
│   ├── service/                # Бизнес-логика (сервисный слой)
│   ├── handler/                # gRPC хендлеры
│   ├── repository/             # Репозитории (PG, Redis)
│   └── errors/                 # Доменные ошибки
│
├── pkg/
│   ├── proto/                  # Сгенерированный Proto код
│   ├── bcrypt/                 # Хеширование паролей
│   ├── jwt/                    # JWT утилиты
│   └── cookies/                # Cookie утилиты
│
├── api/                        # Swagger/OpenAPI спецификации
├── proto/
│   └── auth/v1/
│       └── auth.proto          # Proto контракты
│
├── docs/
│   ├── api.md                  # API документация
│   └── README.md               # Этот файл
│
├── Taskfile.yml                # Taskfile команды
└── go.mod
```

---

## 🚀 Быстрый старт

### Требования

| Зависимость | Версия |
|-------------|--------|
| Go | 1.25+ |
| PostgreSQL | 15+ |
| Redis | 7+ |
| Task | latest |

### Установка зависимостей

```bash
# Установка инструментов разработки
task proto:install-plugins
task install-buf
task install-formatters
task install-golangci-lint
```

### Генерация Proto

```bash
task proto:gen
```

### Запуск

```bash
# Переменные окружения
export DATABASE_URL="postgres://user:pass@localhost:5432/auth?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="your-secret-key"

# Запуск сервера
go run cmd/server/main.go
```

---

## 📡 API Endpoints

### Основные методы

| Метод | gRPC | HTTP | Описание |
|-------|------|------|----------|
| `Register` | `Register` | `POST /api/auth/register` | Регистрация |
| `Login` | `Login` | `POST /api/auth/login` | Вход |
| `Logout` | `Logout` | `POST /api/auth/logout` | Выход |
| `Refresh` | `Refresh` | `POST /api/auth/refresh` | Обновление токенов |

### Формат ответов

Все ответы API имеют единую структуру:

```json
{
  "status_code": 200,
  "message": "Success",
  "data": { ... }
}
```

| Поле | Тип | Описание |
|------|-----|----------|
| `status_code` | `int` | HTTP статус код |
| `message` | `string` | Сообщение статуса |
| `data` | `object` | Тело ответа (может быть `null`) |

📖 **Подробное описание API:** [docs/api.md](api.md)

---

## 🛠 Разработка

### Форматирование

```bash
task format
```

### Линтинг

```bash
# Go код
task lint

# Proto файлы
task proto:lint
```

### Тесты

```bash
go test ./... -v
```

### Taskfile команды

| Команда | Описание |
|---------|----------|
| `task proto:gen` | Генерация Proto (gRPC + REST + Swagger) |
| `task proto:lint` | Линтинг Proto файлов |
| `task proto:deps` | Обновление Proto зависимостей |
| `task format` | Форматирование Go кода |
| `task lint` | Линтинг Go кода |
| `task tidy` | Очистка зависимостей |

---

## 🔒 Безопасность

### JWT Claims

Токены содержат название сервиса в поле `iss` (issuer):

```json
{
  "iss": "auth-service",
  "sub": "{user_id}",
  "email": "{email}",
  "username": "{username}",
  "iat": 1705312200,
  "exp": 1705313100,
  "type": "access"
}
```

### Требования к паролю

| Требование | Значение |
|------------|----------|
| Минимальная длина | 8 символов |
| Заглавные буквы | Минимум 1 (A-Z) |
| Строчные буквы | Минимум 1 (a-z) |
| Цифры | Минимум 1 (0-9) |

### Требования к username

| Требование | Значение |
|------------|----------|
| Длина | 3-30 символов |
| Допустимые символы | Буквы, цифры, `_` |
| Ограничения | Не может начинаться/заканчиваться на `_` |

### Требования к email

| Требование | Значение |
|------------|----------|
| Формат | RFC 5321 |
| Максимальная длина | 254 символа |

---

## 📦 Зависимости

### Внешние библиотеки

```go
// Авторские библиотеки
github.com/dfsfGfd/redis-connect        // Redis клиент
github.com/dfsfGfd/postgresql-connect   // PostgreSQL клиент (pgx)

// Стандартные библиотеки
github.com/google/uuid                  // UUID генерация
golang.org/x/crypto                     // bcrypt
google.golang.org/grpc                  // gRPC
google.golang.org/protobuf              // Protocol Buffers
```

### Proto зависимости

```yaml
deps:
  - buf.build/googleapis/googleapis
  - buf.build/grpc-ecosystem/grpc-gateway
```

---

## 📚 Документация

| Документ | Описание |
|----------|----------|
| [API Documentation](api.md) | Полное описание API endpoints |
| [Swagger UI](http://localhost:8080/swagger/) | Интерактивная документация (после запуска) |

---

## 📝 License

MIT
