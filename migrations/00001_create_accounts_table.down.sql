-- +migrate down
-- Откат миграции: удаление таблицы accounts

DROP INDEX IF EXISTS idx_accounts_email;
DROP TABLE IF EXISTS accounts;
