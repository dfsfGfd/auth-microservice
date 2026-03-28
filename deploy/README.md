# Deploy — Auth Microservice

Запуск микросервиса в Docker.

---

## 🚀 Быстрый старт

### Через Taskfile (рекомендуется)

```bash
task server:docker:up
```

### Напрямую

```bash
cd deploy
docker compose up -d --build
```

Сервисы:
- **auth-service:** `http://localhost:8080`
- **postgres:** `localhost:5432`
- **redis:** `localhost:6379`

---

## 📁 Файлы

| Файл | Описание |
|------|----------|
| `docker-compose.yml` | Development окружение |
| `Dockerfile` | Multi-stage сборка |
| `.env` | Переменные окружения |
| `.dockerignore` | Исключения для Docker |

---

## 🔧 Команды

### Taskfile

```bash
task server:docker:up       # Запуск
task server:docker:down     # Остановка
task server:docker:logs     # Логи
task server:docker:restart  # Перезапуск
```

### Docker Compose

```bash
# Запуск
docker compose up -d

# Логи
docker compose logs -f auth-service

# Остановка
docker compose down

# Пересборка
docker compose up -d --build --force-recreate
```

---

## 🏗 Архитектура

**Dockerfile:**
- Build stage: golang:1.26-alpine
- Runtime stage: alpine:3.23
- Размер: ~20MB
- Non-root пользователь: appuser (1000)
- CMD: /app/server

**Миграции:**
- Применяются автоматически при старте сервера (встроены в сервер)
- golang-migrate отслеживает применённые миграции в таблице `schema_migrations`
- При повторном запуске миграции не применяются (данные сохраняются в volume)

---

## 📊 Логирование

**Формат:** JSON (production)

**Пример:**
```json
{"ts":"2026-03-21T15:00:00Z","lvl":"info","msg":"server_start","srv":"auth-service"}
```

---

## 🐛 Troubleshooting

```bash
# Проверка статуса
docker compose ps

# Логи сервиса
docker compose logs auth-service

# Вход в контейнер
docker compose exec auth-service sh

# Перезапуск
docker compose restart auth-service
```

---

## 🔒 Безопасность

**Production checklist:**

- [ ] JWT_SECRET ≥ 32 символов
- [ ] APP_ENV=production
- [ ] LOG_LEVEL=info (не debug)
- [ ] Docker secrets для чувствительных данных

**Пример:**
```bash
# Генерация JWT_SECRET
openssl rand -base64 32 > .jwt_secret

# В docker-compose.yml
- JWT_SECRET=${JWT_SECRET}
```

---

## 📚 Ссылки

- [README](../docs/README.md)
- [Config](../docs/config.md)
- [Taskfile](../Taskfile.yml)
