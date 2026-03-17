-- Migration: 001_create_accounts_table (Down)
-- Description: Удаление таблицы accounts
-- Version: 1.0.0

-- +goose Up
-- +goose StatementBegin

-- Удаляем триггер и функцию
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаляем индексы
DROP INDEX IF EXISTS idx_accounts_created_at;
DROP INDEX IF EXISTS idx_accounts_email;

-- Удаляем таблицу
DROP TABLE IF EXISTS accounts;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Эта миграция не должна применяться в production!
-- Восстановление таблицы из up миграции
SELECT 'This migration should not be reversed in production!' AS warning;

-- +goose StatementEnd
