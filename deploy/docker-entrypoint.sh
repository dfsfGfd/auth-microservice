#!/bin/sh
set -e

# JSON логирование для консистентности
log_json() {
    echo "{\"ts\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",\"lvl\":\"$1\",\"msg\":\"$2\",\"srv\":\"auth-service\"}"
}

log_json "info" "migrations_start"
/app/migrate -dsn "$DATABASE_URL" up
log_json "info" "migrations_complete"

log_json "info" "server_start"
exec /app/server
