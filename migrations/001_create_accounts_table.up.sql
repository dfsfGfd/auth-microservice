-- +goose up
-- Создание таблицы accounts для хранения аккаунтов пользователей
-- Миграция идемпотентна: можно безопасно запускать многократно

-- Таблица аккаунтов
-- Примечание: UUID генерируется в доменном слое (Go), а не в БД
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    email VARCHAR(254) NOT NULL,
    password VARCHAR(72) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Уникальный индекс на email для быстрого поиска (O(1))
CREATE UNIQUE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);

-- Индекс для сортировки по дате создания
CREATE INDEX IF NOT EXISTS idx_accounts_created_at ON accounts(created_at DESC);

-- Функция для обновления updated_at (CREATE OR REPLACE — перезаписывается при каждом запуске)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер создаётся только если не существует (защита от повторного применения)
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'update_accounts_updated_at'
    ) THEN
        CREATE TRIGGER update_accounts_updated_at
            BEFORE UPDATE ON accounts
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- Комментарии (применяются только если таблица существует)
DO $$ BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'accounts') THEN
        COMMENT ON TABLE accounts IS 'Таблица хранения аккаунтов пользователей';
        COMMENT ON COLUMN accounts.id IS 'Уникальный идентификатор аккаунта (UUID)';
        COMMENT ON COLUMN accounts.email IS 'Email пользователя (уникальный, до 254 символов)';
        COMMENT ON COLUMN accounts.password IS 'Хеш пароля (bcrypt, 72 символа)';
        COMMENT ON COLUMN accounts.created_at IS 'Время создания аккаунта';
        COMMENT ON COLUMN accounts.updated_at IS 'Время последнего обновления аккаунта';
    END IF;
END $$;
