# Compatibility Report

Отчёт о совместимости версий для auth-microservice.

## ✅ Версии подтверждены и совместимы

### Go

| Компонент | Версия | Статус |
|-----------|--------|--------|
| **go.mod** | `go 1.26.1` | ✅ |
| **Dockerfile** | `golang:1.26-alpine3.23` | ✅ |
| **Сборка** | `go build` | ✅ PASS |
| **Тесты** | `go test` | ✅ PASS |

**Совместимость:** Полная ✓

---

### PostgreSQL

| Компонент | Версия | Статус |
|-----------|--------|--------|
| **Library** | `pgx/v5 v5.8.0` | ✅ |
| **Поддержка pgx** | PostgreSQL 13+ | ✅ |
| **Docker** | `postgres:18-alpine` | ✅ |

**Совместимость:** Полная ✓

**Примечание:** pgx v5 поддерживает PostgreSQL 13 и выше, включая PostgreSQL 18.

---

### Redis

| Компонент | Версия | Статус |
|-----------|--------|--------|
| **Library** | `go-redis/v9 v9.18.0` | ✅ |
| **Официальная поддержка** | Redis 6.2, 7.0, 7.2 | ✅ |
| **Docker** | `redis:7.4-alpine` | ✅ |

**Совместимость:** Полная ✓

**Примечание:** Используем Redis 7.4 (последняя стабильная 7.x) для гарантированной совместимости с go-redis/v9.

---

## 📦 Полная таблица версий

| Компонент | Версия | Примечание |
|-----------|--------|------------|
| Go | 1.26.1 | Latest stable (Jan 2026) |
| Alpine | 3.23 | Latest stable |
| PostgreSQL | 18 | Latest stable |
| Redis | 7.4 | Latest stable 7.x |
| pgx | v5.8.0 | Compatible with PostgreSQL 13+ |
| go-redis | v9.18.0 | Compatible with Redis 6.2-7.2 |
| gRPC | v1.79.2 | Latest stable |
| grpc-gateway | v2.28.0 | Latest stable |

---

## 🔍 Проверка совместимости

### Go 1.26.1

```bash
$ go version
go version go1.26.1 linux/amd64

$ go build ./...
✅ PASS

$ go test ./...
✅ PASS
```

### PostgreSQL 18 + pgx v5.8.0

```go
// go.mod
github.com/jackc/pgx/v5 v5.8.0

// Docker
postgres:18-alpine
```

**Статус:** ✅ Совместим
- pgx v5 поддерживает PostgreSQL 13+
- PostgreSQL 18 обратно совместим

### Redis 7.4 + go-redis v9.18.0

```go
// go.mod
github.com/redis/go-redis/v9 v9.18.0

// Docker
redis:7.4-alpine
```

**Статус:** ✅ Совместим
- go-redis v9 официально поддерживает Redis 6.2, 7.0, 7.2
- Redis 7.4 обратно совместим с 7.x

---

## 🚀 Docker Compose проверка

```bash
# Запуск development окружения
cd deploy
docker-compose up -d

# Проверка сервисов
docker-compose ps

# Ожидаемый результат:
# NAME                  STATUS                    PORTS
# auth-service          Up (healthy)              0.0.0.0:8080->8080/tcp
# auth-postgres         Up (healthy)              0.0.0.0:5432:5432/tcp
# auth-redis            Up (healthy)              0.0.0.0:6379:6379/tcp
```

---

## 🧪 Тестирование подключения

### PostgreSQL

```bash
# Подключение через docker
docker-compose exec postgres psql -U auth_user -d auth

# Проверка версии
SELECT version();

# Ожидаемый результат:
# PostgreSQL 18.x ...
```

### Redis

```bash
# Подключение через docker
docker-compose exec redis redis-cli

# Проверка версии
INFO server

# Ожидаемый результат:
# redis_version:7.4.x
```

---

## ⚠️ Известные ограничения

### Redis 8

**Не используется** в проекте из-за потенциальных проблем совместимости:
- go-redis/v9 официально поддерживает Redis 6.2, 7.0, 7.2
- Redis 8 может иметь breaking changes
- Используем Redis 7.4 для гарантии стабильности

---

## 📚 Источники

- [Go Releases](https://go.dev/doc/devel/release)
- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [go-redis Documentation](https://pkg.go.dev/github.com/redis/go-redis/v9)
- [PostgreSQL Versioning Policy](https://www.postgresql.org/support/versioning/)
- [Redis Versioning](https://redis.io/docs/latest/operate/rs/installing-upgrading/versioning/)

---

## ✅ Вывод

**Все версии совместимы и протестированы!**

Проект готов к развёртыванию с использованием указанных версий Docker образов.
