-- Migration: 001_create_accounts_table
-- Description: Создание таблицы accounts для хранения аккаунтов пользователей
-- Version: 1.0.0
-- Created: 2024-01-15

-- +goose Up
-- +goose StatementBegin

-- Создаём таблицу accounts
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(254) NOT NULL,
    password VARCHAR(72) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индексы для производительности
-- idx_accounts_email - для быстрого поиска по email (login, registration check)
CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);

-- idx_accounts_created_at - для аналитики и пагинации
CREATE INDEX IF NOT EXISTS idx_accounts_created_at ON accounts(created_at DESC);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_accounts_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Комментарии для документации
COMMENT ON TABLE accounts IS 'Таблица хранения аккаунтов пользователей';
COMMENT ON COLUMN accounts.id IS 'Уникальный идентификатор аккаунта (UUID)';
COMMENT ON COLUMN accounts.email IS 'Email пользователя (уникальный, до 254 символов)';
COMMENT ON COLUMN accounts.password IS 'Хеш пароля (bcrypt, 72 символа)';
COMMENT ON COLUMN accounts.created_at IS 'Время создания аккаунта';
COMMENT ON COLUMN accounts.updated_at IS 'Время последнего обновления аккаунта';

-- +goose StatementEnd

-- +goose Down
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
