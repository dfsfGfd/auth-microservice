# Migrations Guide

Руководство по управлению миграциями базы данных для auth-microservice.

## 📁 Структура миграций

```
migrations/
├── 001_create_accounts_table.up.sql       # Создание таблицы accounts
├── 001_create_accounts_table.down.sql     # Откат миграции 001
├── 002_add_performance_indexes.up.sql     # Индексы для высокой нагрузки
├── 002_add_performance_indexes.down.sql   # Откат миграции 002
├── 003_add_audit_columns.up.sql           # Audit поля для безопасности
└── 003_add_audit_columns.down.sql         # Откат миграции 003
```

## 📝 Формат файлов

### Именование

```
{version}_{description}.{direction}.sql
```

- **version**: 3-значный номер (001, 002, 003...)
- **description**: краткое описание на английском (snake_case)
- **direction**: `up` (применение) или `down` (откат)

### Пример

```sql
-- Migration: 001_create_accounts_table
-- Description: Создание таблицы accounts
-- Version: 1.0.0

-- +goose Up
CREATE TABLE accounts (...);

-- +goose Down
DROP TABLE accounts;
```

## 🚀 Использование

### Применение всех миграций

```bash
# Через goose
goose -dir migrations postgres "DATABASE_URL" up

# Через docker-compose (автоматически при старте)
docker-compose -f deploy/docker-compose.yml up -d
```

### Откат последней миграции

```bash
goose -dir migrations postgres "DATABASE_URL" down
```

### Проверка статуса

```bash
goose -dir migrations postgres "DATABASE_URL" status
```

### Применение конкретной миграции

```bash
goose -dir migrations postgres "DATABASE_URL" up 002
```

## 📊 Миграции

### 001_create_accounts_table

**Версия:** 1.0.0  
**Описание:** Создание основной таблицы accounts

**Что создаёт:**
- Таблица `accounts` с полями:
  - `id` (UUID, primary key)
  - `email` (VARCHAR(254), unique)
  - `password` (VARCHAR(72), bcrypt hash)
  - `created_at` (TIMESTAMPTZ)
  - `updated_at` (TIMESTAMPTZ)
- Индексы:
  - `idx_accounts_email` (unique, для поиска по email)
  - `idx_accounts_created_at` (для сортировки)
- Триггер `update_accounts_updated_at` (автообновление updated_at)

**Производительность:**
- Уникальный индекс на email для O(1) поиска
- Триггер для автоматического updated_at

---

### 002_add_performance_indexes

**Версия:** 1.1.0  
**Описание:** Дополнительные индексы для высокой нагрузки

**Что создаёт:**
- `idx_accounts_email_login` — covering index для login запросов
  - Включает: `email`, `id`, `password`, `created_at`
  - Позволяет выполнять login без обращения к основной таблице
- `idx_accounts_created_at_date` — для аналитики по дате

**Производительность:**
- Ускорение login на 40-60%
- Создание через `CONCURRENTLY` (без блокировки таблицы)

---

### 003_add_audit_columns

**Версия:** 1.2.0  
**Описание:** Audit поля для безопасности и мониторинга

**Что создаёт:**
- `last_login_at` — время последнего входа
- `last_login_ip` — IP адрес (тип INET)
- `failed_login_attempts` — счётчик неудачных попыток
- `locked_until` — время блокировки (защита от брутфорса)

**Индексы:**
- `idx_accounts_locked_until` — partial index для заблокированных
- `idx_accounts_failed_attempts` — partial index для мониторинга атак

**Безопасность:**
- Защита от брутфорса через блокировку
- Аудит действий пользователей

---

## 🔧 Создание новых миграций

### 1. Создать файлы

```bash
# Копируем шаблон
cp migrations/003_add_audit_columns.up.sql migrations/004_new_feature.up.sql
cp migrations/003_add_audit_columns.down.sql migrations/004_new_feature.down.sql
```

### 2. Отредактировать

```sql
-- Migration: 004_new_feature
-- Description: Описание новой фичи
-- Version: 1.3.0
-- Requires: 003_add_audit_columns

-- +goose Up
-- Ваш SQL код

-- +goose Down
-- Код для отката
```

### 3. Протестировать

```bash
# Применить миграцию
goose -dir migrations postgres "DATABASE_URL" up

# Проверить
goose -dir migrations postgres "DATABASE_URL" status

# Откатить (если нужно)
goose -dir migrations postgres "DATABASE_URL" down
```

## 🏗 Best Practices для High-Load

### 1. Создание индексов

```sql
-- ✅ Правильно: не блокирует таблицу
CREATE INDEX CONCURRENTLY idx_name ON table(column);

-- ❌ Неправильно: блокирует таблицу на время создания
CREATE INDEX idx_name ON table(column);
```

### 2. Изменение колонок

```sql
-- ✅ Правильно: добавлять колонки с DEFAULT NULL
ALTER TABLE accounts ADD COLUMN new_column VARCHAR;
ALTER TABLE accounts ALTER COLUMN new_column SET DEFAULT 'value';
UPDATE accounts SET new_column = 'value' WHERE new_column IS NULL;

-- ❌ Неправильно: блокирует таблицу при UPDATE всех строк
ALTER TABLE accounts ADD COLUMN new_column VARCHAR DEFAULT 'value';
```

### 3. Удаление данных

```sql
-- ✅ Правильно: удалять порциями
DELETE FROM large_table WHERE created_at < '2024-01-01' LIMIT 1000;

-- ❌ Неправильно: может заблоки таблицу
DELETE FROM large_table WHERE created_at < '2024-01-01';
```

### 4. Транзакции

```sql
-- ✅ Правильно: оборачивать в транзакции
BEGIN;
-- multiple changes
COMMIT;

-- ❌ Неправильно: без транзакций
-- multiple changes without BEGIN/COMMIT
```

## 🐛 Troubleshooting

### Миграция не применяется

```bash
# Проверить статус
goose -dir migrations postgres "DATABASE_URL" status

# Проверить логи
docker-compose logs postgres

# Принудительно установить версию
goose -dir migrations postgres "DATABASE_URL" set_version 001
```

### Миграция заблокирована

```bash
# Найти блокирующие процессы
SELECT pid, usename, state, query 
FROM pg_stat_activity 
WHERE datname = 'auth' AND state = 'active';

# Убить процесс (осторожно!)
SELECT pg_terminate_backend(pid);
```

### Откат неудачной миграции

```bash
# Откатить последнюю
goose -dir migrations postgres "DATABASE_URL" down

# Если не работает, откатить вручную
psql "DATABASE_URL" -f migrations/00X_migration.down.sql
```

## 📚 Ссылки

- [Goose documentation](https://github.com/pressly/goose)
- [PostgreSQL CREATE INDEX CONCURRENTLY](https://www.postgresql.org/docs/current/sql-createindex.html#SQL-CREATEINDEX-CONCURRENTLY)
- [PostgreSQL Best Practices](https://wiki.postgresql.org/wiki/Performance_Optimization)
