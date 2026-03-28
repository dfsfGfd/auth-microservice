-- +migrate up
-- Создание таблицы accounts для хранения аккаунтов пользователей

CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    email VARCHAR(254) NOT NULL,
    password VARCHAR(72) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_accounts_email ON accounts(email);
CREATE INDEX idx_accounts_created_at ON accounts(created_at DESC);

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

COMMENT ON TABLE accounts IS 'Таблица хранения аккаунтов пользователей';
COMMENT ON COLUMN accounts.id IS 'Уникальный идентификатор аккаунта (UUID)';
COMMENT ON COLUMN accounts.email IS 'Email пользователя (уникальный, до 254 символов)';
COMMENT ON COLUMN accounts.password IS 'Хеш пароля (bcrypt, 72 символа)';
COMMENT ON COLUMN accounts.created_at IS 'Время создания аккаунта';
COMMENT ON COLUMN accounts.updated_at IS 'Время последнего обновления аккаунта';
