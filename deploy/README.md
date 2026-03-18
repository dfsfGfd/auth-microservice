# Deploy — Auth Microservice

Запуск микросервиса в Docker.

---

## 🚀 Быстрый старт

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
| `docker-entrypoint.sh` | Миграции + запуск сервера |
| `.env` | Переменные окружения |
| `.dockerignore` | Исключения для Docker |

---

## 🔧 Команды

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

**docker-entrypoint.sh:**
1. Применение миграций
2. Запуск сервера

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

## 📚 Ссылки

- [README](../docs/README.md)
- [Config](../docs/config.md)
