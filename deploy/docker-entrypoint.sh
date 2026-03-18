#!/bin/sh
set -e

echo "🔄 Running database migrations..."

# Запускаем миграции
/app/migrate -dsn "$DATABASE_URL" up

echo "✅ Migrations completed"
echo "🚀 Starting auth service..."

# Запускаем сервер
exec /app/server
