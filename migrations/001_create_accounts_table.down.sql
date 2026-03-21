-- +migrate down
-- Откат миграции: удаление таблицы accounts

-- Удаляем триггер и функцию
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Удаляем индексы
DROP INDEX IF EXISTS idx_accounts_created_at;
DROP INDEX IF EXISTS idx_accounts_email;

-- Удаляем таблицу
DROP TABLE IF EXISTS accounts;
