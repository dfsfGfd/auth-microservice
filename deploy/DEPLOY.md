# Deploy — Auth Microservice

Папка содержит все необходимые файлы для развёртывания микросервиса аутентификации.

## 📁 Структура

```
deploy/
├── README.md                      # Этот файл
├── Dockerfile                     # Multi-stage Dockerfile
├── .dockerignore                  # Исключения для Docker build
├── docker-compose.yml             # Development окружение
├── docker-compose.production.yml  # Production окружение
└── .env.production.example        # Шаблон production переменных
```

## 🚀 Быстрый старт

### 1. Разработка (локально)

```bash
# Из корня проекта
cd /home/user/Documents/auth

# Запуск всех сервисов
docker-compose -f deploy/docker-compose.yml up -d

# Проверка статуса
docker-compose -f deploy/docker-compose.yml ps

# Логи
docker-compose -f deploy/docker-compose.yml logs -f auth-service

# Остановка
docker-compose -f deploy/docker-compose.yml down
```

### 2. Production

```bash
# 1. Скопируйте шаблон
cp deploy/.env.production.example deploy/.env.production

# 2. Отредактируйте переменные (обязательно измените пароли!)
nano deploy/.env.production

# 3. Запуск
docker-compose -f deploy/docker-compose.production.yml up -d --build

# 4. Проверка
docker-compose -f deploy/docker-compose.production.yml ps
curl http://localhost:8080/health
```

## 📚 Документация

- **[README.md](README.md)** — Полное руководство по Docker развёртыванию

## 🔧 Файлы

### Dockerfile

Multi-stage Dockerfile для минимального размера образа:

- **Build stage**: golang:1.25-alpine
- **Runtime stage**: alpine:3.21
- **Размер**: ~20MB
- **Безопасность**: non-root user, health checks

### docker-compose.yml

Development окружение:

- auth-service (Go)
- postgres:15-alpine
- redis:7-alpine
- Pre-configured network и volumes

### docker-compose.production.yml

Production окружение:

- Resource limits (CPU, memory)
- Security hardening
- JSON logging с ротацией
- Health checks
- Scalability support

### .env.production.example

Шаблон для production переменных:

```bash
# Обязательно измените!
JWT_SECRET=<generate-new-secret>
DB_PASSWORD=<generate-strong-password>
REDIS_PASSWORD=<generate-strong-password>
```

## 🔒 Безопасность

### Checklist для production

- [ ] Сгенерировать новый `JWT_SECRET` (минимум 32 символа)
- [ ] Изменить все пароли в `.env.production`
- [ ] Включить SSL для PostgreSQL
- [ ] Ограничить доступ к портам (firewall)
- [ ] Не коммитить `.env.production` в git
- [ ] Включить HTTPS (reverse proxy)

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
# Проверка сервиса
curl http://localhost:8080/health

# Статус контейнеров
docker-compose -f deploy/docker-compose.production.yml ps
```

### Логи

```bash
# Последние логи
docker-compose -f deploy/docker-compose.production.yml logs --tail=100

# Логи в реальном времени
docker-compose -f deploy/docker-compose.production.yml logs -f auth-service
```

### Метрики

```bash
# Использование ресурсов
docker stats auth-service postgres redis
```

## 🛠 Отладка

### Вход в контейнер

```bash
# Auth service
docker-compose -f deploy/docker-compose.yml exec auth-service sh

# PostgreSQL
docker-compose -f deploy/docker-compose.yml exec postgres psql -U auth_user -d auth

# Redis
docker-compose -f deploy/docker-compose.yml exec redis redis-cli
```

### Перезапуск

```bash
# Перезапуск сервиса
docker-compose -f deploy/docker-compose.yml restart auth-service

# Полная пересборка
docker-compose -f deploy/docker-compose.yml up -d --build --force-recreate
```

## 📦 Docker Compose команды

| Команда | Описание |
|---------|----------|
| `up -d` | Запуск в фоновом режиме |
| `down` | Остановка |
| `down -v` | Остановка с удалением volumes |
| `logs -f` | Просмотр логов |
| `ps` | Статус контейнеров |
| `restart <service>` | Перезапуск сервиса |
| `exec <service> sh` | Вход в контейнер |
| `build` | Пересборка образов |

## 🐛 Troubleshooting

### Сервис не запускается

```bash
# Проверьте логи
docker-compose -f deploy/docker-compose.yml logs auth-service

# Проверьте переменные окружения
docker-compose -f deploy/docker-compose.yml exec auth-service env

# Проверьте подключение к БД
docker-compose -f deploy/docker-compose.yml exec auth-service \
  wget --spider postgres://auth_user:auth_password@postgres:5432/auth
```

### Проблемы с PostgreSQL

```bash
# Статус PostgreSQL
docker-compose -f deploy/docker-compose.yml exec postgres pg_isready -U auth_user -d auth

# Логи PostgreSQL
docker-compose -f deploy/docker-compose.yml logs postgres
```

### Проблемы с Redis

```bash
# Проверка подключения
docker-compose -f deploy/docker-compose.yml exec redis redis-cli ping

# Логи Redis
docker-compose -f deploy/docker-compose.yml logs redis
```

## 📚 Ссылки

- [Docker Compose documentation](https://docs.docker.com/compose/)
- [Dockerfile best practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [PostgreSQL Docker Hub](https://hub.docker.com/_/postgres)
- [Redis Docker Hub](https://hub.docker.com/_/redis)
