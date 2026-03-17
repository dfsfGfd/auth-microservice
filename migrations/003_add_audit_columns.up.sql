-- Migration: 003_add_audit_columns
-- Description: Добавление audit полей для отслеживания изменений
-- Version: 1.2.0
-- Created: 2024-01-15
-- Requires: 001_create_accounts_table

-- +goose Up
-- +goose StatementBegin

-- Добавляем поля для аудита (если их нет)
ALTER TABLE accounts 
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_login_ip INET,
    ADD COLUMN IF NOT EXISTS failed_login_attempts INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS locked_until TIMESTAMPTZ;

-- Индексы для новых полей
-- Для поиска заблокированных аккаунтов
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_accounts_locked_until 
ON accounts(locked_until) WHERE locked_until IS NOT NULL;

-- Для мониторинга failed login attempts
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_accounts_failed_attempts 
ON accounts(failed_login_attempts) WHERE failed_login_attempts > 0;

-- Комментарии
COMMENT ON COLUMN accounts.last_login_at IS 'Время последнего успешного входа';
COMMENT ON COLUMN accounts.last_login_ip IS 'IP адрес последнего входа (для аудита)';
COMMENT ON COLUMN accounts.failed_login_attempts IS 'Количество неудачных попыток входа';
COMMENT ON COLUMN accounts.locked_until IS 'Время блокировки аккаунта (защита от брутфорса)';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Удаляем индексы
DROP INDEX IF EXISTS idx_accounts_failed_attempts;
DROP INDEX IF EXISTS idx_accounts_locked_until;

-- Удаляем колонки (осторожно в production!)
ALTER TABLE accounts 
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS last_login_ip,
    DROP COLUMN IF EXISTS failed_login_attempts,
    DROP COLUMN IF EXISTS locked_until;

-- +goose StatementEnd
