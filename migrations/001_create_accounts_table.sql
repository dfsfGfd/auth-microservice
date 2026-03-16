-- Migration: 001_create_accounts_table.sql
-- Description: Создание таблицы accounts (новая сущность вместо users)
-- Up: Создать таблицу accounts
-- Down: Удалить таблицу accounts

-- +goose Up
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(254) NOT NULL UNIQUE,
    password VARCHAR(72) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Индекс для быстрого поиска по email
CREATE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);

-- +goose Down
DROP TABLE IF EXISTS accounts;
