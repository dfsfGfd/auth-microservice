-- +migrate up
-- Создание таблицы accounts для хранения аккаунтов пользователей

CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    email VARCHAR(254) NOT NULL,
    password VARCHAR(72) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX idx_accounts_email ON accounts(email);

COMMENT ON TABLE accounts IS 'Таблица хранения аккаунтов пользователей';
COMMENT ON COLUMN accounts.id IS 'Уникальный идентификатор аккаунта (UUID)';
COMMENT ON COLUMN accounts.email IS 'Email пользователя (уникальный, до 254 символов)';
COMMENT ON COLUMN accounts.password IS 'Хеш пароля (bcrypt, 72 символа)';
COMMENT ON COLUMN accounts.created_at IS 'Время создания аккаунта';
COMMENT ON COLUMN accounts.updated_at IS 'Время последнего обновления аккаунта';
