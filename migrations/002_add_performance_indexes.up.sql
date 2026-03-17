-- Migration: 002_add_performance_indexes
-- Description: Добавление дополнительных индексов для высоконагруженного сервиса
-- Version: 1.1.0
-- Created: 2024-01-15
-- Requires: 001_create_accounts_table

-- +goose Up
-- +goose StatementBegin

-- CONCURRENTLY не блокирует таблицу при создании индекса (важно для production)
-- Добавляем covering index для оптимизации login запросов
-- Этот индекс позволяет выполнять поиск по email без обращения к основной таблице
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_accounts_email_login 
ON accounts(email) INCLUDE (id, password, created_at);

-- Индекс для аналитики по дате регистрации
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_accounts_created_at_date 
ON accounts((created_at::date));

-- Статистика для оптимизатора запросов (важно для больших таблиц)
ANALYZE accounts;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_accounts_email_login;
DROP INDEX IF EXISTS idx_accounts_created_at_date;

-- +goose StatementEnd
