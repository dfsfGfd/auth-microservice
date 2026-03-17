# ==================== BUILD STAGE ====================
FROM golang:1.26-alpine3.23 AS builder

# Устанавливаем необходимые зависимости
RUN apk add --no-cache git ca-certificates tzdata

# Создаём пользователя для запуска приложения
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -D appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости (кэшируется слоем)
RUN go mod download && go mod verify

# Копируем исходный код
COPY . .

# Собираем бинарный файл с оптимизациями
# -ldflags "-s -w" удаляет символы отладки и таблицу символов
# CGO_ENABLED=0 создаёт статический бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X 'main.Version=${VERSION:-dev}'" \
    -o /app/bin/server \
    ./cmd/server/main.go

# ==================== RUNTIME STAGE ====================
FROM alpine:3.23 AS runtime

# Устанавливаем необходимые пакеты
RUN apk add --no-cache ca-certificates tzdata

# Создаём пользователя для запуска приложения (без root)
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -D appuser

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из builder
COPY --from=builder /app/bin/server /app/server

# Копируем .env.example как шаблон
COPY --from=builder /app/.env.example /app/.env.example

# Устанавливаем владельца на appuser
RUN chown -R appuser:appgroup /app

# Переключаемся на не-root пользователя
USER appuser

# Открываем порты
EXPOSE 8080 9090

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Запускаем сервер
ENTRYPOINT ["/app/server"]
