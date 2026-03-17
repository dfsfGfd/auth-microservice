# Docker Guide

Руководство по развёртыванию auth-microservice в Docker.

## 📋 Быстрый старт

### 1. Разработка (локально)

```bash
# Запуск всех сервисов (auth + postgres + redis)
docker-compose up -d

# Просмотр логов
docker-compose logs -f auth-service

# Остановка
docker-compose down
```

### 2. Production

```bash
# 1. Скопируйте production конфиг
cp .env.production.example .env.production

# 2. Настройте переменные в .env.production
# Обязательно: JWT_SECRET, DB_PASSWORD, REDIS_PASSWORD

# 3. Соберите и запустите
docker-compose -f docker-compose.production.yml up -d --build

# 4. Проверьте статус
docker-compose -f docker-compose.production.yml ps
```

## 🏗 Архитектура

```
┌─────────────────────────────────────────────────┐
│              Docker Network                     │
│                                                 │
│  ┌──────────────┐    ┌──────────────┐          │
│  │   auth-      │───▶│   postgres   │          │
│  │   service    │    │   :5432      │          │
│  │   :8080      │    └──────────────┘          │
│  │   :9090      │                               │
│  └──────┬───────┘    ┌──────────────┐          │
│         │           │    redis      │          │
│         └──────────▶│   :6379      │          │
│                     └──────────────┘          │
└─────────────────────────────────────────────────┘
           │
           ▼
    Host Ports
    8080 (HTTP)
    9090 (gRPC)
    5432 (PostgreSQL)
    6379 (Redis)
```

## 🔧 Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `APP_ENV` | Окружение (development/production) | `development` |
| `HTTP_PORT` | Порт HTTP сервера | `8080` |
| `GRPC_PORT` | Порт gRPC сервера | `9090` |
| `DATABASE_URL` | PostgreSQL connection string | - |
| `REDIS_URL` | Redis connection string | - |
| `JWT_SECRET` | Секретный ключ JWT (мин 32 символа) | - |
| `LOG_LEVEL` | Уровень логирования | `info` |
| `LOG_FORMAT` | Формат логов (json/console) | `json` |

### Порты

| Сервис | Порт | Описание |
|--------|------|----------|
| auth-service | 8080 | HTTP (REST + grpc-gateway) |
| auth-service | 9090 | gRPC (прямой доступ) |
| postgres | 5432 | PostgreSQL |
| redis | 6379 | Redis |

## 🚀 Команды

### Разработка

```bash
# Запуск
docker-compose up -d

# Пересборка
docker-compose up -d --build

# Логи
docker-compose logs -f auth-service
docker-compose logs -f postgres
docker-compose logs -f redis

# Остановка
docker-compose down

# Остановка с удалением volumes
docker-compose down -v
```

### Production

```bash
# Запуск
docker-compose -f docker-compose.production.yml up -d

# Масштабирование (если нужно несколько реплик)
docker-compose -f docker-compose.production.yml up -d --scale auth-service=3

# Обновление
docker-compose -f docker-compose.production.yml pull
docker-compose -f docker-compose.production.yml up -d --build

# Логи
docker-compose -f docker-compose.production.yml logs -f auth-service
```

## 🔒 Безопасность

### Production Checklist

- [ ] Измените все пароли в `.env.production`
- [ ] Сгенерируйте новый `JWT_SECRET` (минимум 32 символа)
- [ ] Используйте SSL для PostgreSQL (`sslmode=require`)
- [ ] Включите аутентификацию Redis (`REDIS_PASSWORD`)
- [ ] Не коммитьте `.env.production` в git
- [ ] Ограничьте доступ к портам (firewall)

### Генерация секретов

```bash
# JWT_SECRET
openssl rand -base64 32

# DB_PASSWORD
openssl rand -base64 32

# REDIS_PASSWORD
openssl rand -base64 32
```

## 📊 Мониторинг

### Health Check

```bash
# Проверка здоровья сервиса
curl http://localhost:8080/health

# Статус контейнеров
docker-compose ps

# Детальная информация
docker inspect auth-service
```

### Логи

```bash
# Последние 100 строк
docker-compose logs --tail=100 auth-service

# Логи за последнее время
docker-compose logs --since=1h auth-service

# JSON формат (production)
docker-compose logs auth-service | jq .
```

### Метрики

```bash
# Использование ресурсов
docker stats auth-service postgres redis

# Детальная статистика
docker stats --no-stream
```

## 🛠 Отладка

### Вход в контейнер

```bash
# Auth service
docker-compose exec auth-service sh

# PostgreSQL
docker-compose exec postgres psql -U auth_user -d auth

# Redis
docker-compose exec redis redis-cli
```

### Перезапуск сервисов

```bash
# Перезапуск auth-service
docker-compose restart auth-service

# Перезапуск всех
docker-compose restart
```

### Очистка

```bash
# Удалить все контейнеры и volumes
docker-compose down -v

# Удалить образы
docker-compose down -v --rmi all

# Полная очистка (production)
docker-compose -f docker-compose.production.yml down -v --rmi all
```

## 📦 Dockerfile

### Multi-stage сборка

```dockerfile
# Stage 1: Build
FROM golang:1.25-alpine AS builder
# ... сборка бинарного файла

# Stage 2: Runtime
FROM alpine:3.21
# ... минимальный runtime образ
```

### Оптимизации

- **Multi-stage**: Уменьшает размер образа с ~1GB до ~20MB
- **CGO_ENABLED=0**: Статический бинарник без зависимостей
- **ldflags "-s -w"**: Удаляет символы отладки
- **Alpine**: Минимальный базовый образ
- **Non-root user**: Запуск от имени пользователя appuser (UID 1000)

## 🐛 Troubleshooting

### Сервис не запускается

```bash
# Проверьте логи
docker-compose logs auth-service

# Проверьте переменные окружения
docker-compose exec auth-service env | grep -E 'DATABASE|REDIS|JWT'

# Проверьте подключение к БД
docker-compose exec auth-service wget --spider postgres://auth_user:auth_password@postgres:5432/auth
```

### Проблемы с подключением к PostgreSQL

```bash
# Проверьте статус PostgreSQL
docker-compose exec postgres pg_isready -U auth_user -d auth

# Проверьте логи PostgreSQL
docker-compose logs postgres
```

### Проблемы с Redis

```bash
# Проверьте подключение к Redis
docker-compose exec redis redis-cli ping

# Проверьте логи Redis
docker-compose logs redis
```

## 📚 Ссылки

- [Docker Compose documentation](https://docs.docker.com/compose/)
- [Dockerfile best practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [PostgreSQL Docker Hub](https://hub.docker.com/_/postgres)
- [Redis Docker Hub](https://hub.docker.com/_/redis)
