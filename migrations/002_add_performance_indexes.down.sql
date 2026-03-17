-- Migration: 002_add_performance_indexes (Down)
-- Description: Удаление дополнительных индексов
-- Version: 1.1.0

-- +goose Up
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_accounts_email_login;
DROP INDEX IF EXISTS idx_accounts_created_at_date;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Восстановление индексов (не рекомендуется в production)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_accounts_email_login 
ON accounts(email) INCLUDE (id, password, created_at);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_accounts_created_at_date 
ON accounts((created_at::date));

-- +goose StatementEnd
