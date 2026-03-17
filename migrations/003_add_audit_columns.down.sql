-- Migration: 003_add_audit_columns (Down)
-- Description: Удаление audit полей
-- Version: 1.2.0

-- +goose Up
-- +goose StatementBegin

-- Удаляем индексы
DROP INDEX IF EXISTS idx_accounts_failed_attempts;
DROP INDEX IF EXISTS idx_accounts_locked_until;

-- Удаляем колонки
ALTER TABLE accounts 
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS last_login_ip,
    DROP COLUMN IF EXISTS failed_login_attempts,
    DROP COLUMN IF EXISTS locked_until;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Восстановление полей (не рекомендуется в production)
ALTER TABLE accounts 
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_login_ip INET,
    ADD COLUMN IF NOT EXISTS failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_accounts_locked_until 
ON accounts(locked_until) WHERE locked_until IS NOT NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_accounts_failed_attempts 
ON accounts(failed_login_attempts) WHERE failed_login_attempts > 0;

-- +goose StatementEnd
