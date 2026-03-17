# Migrations Guide

Руководство по управлению миграциями базы данных с использованием golang-migrate.

## 📁 Структура

```
migrations/
└── 001_create_accounts_table.up.sql    # Применение миграции
└── 001_create_accounts_table.down.sql  # Откат миграции
```

## 🚀 Использование

### Применение всех миграций

```bash
# Из корня проекта
go run cmd/migrate/main.go -dsn "postgres://user:pass@localhost:5432/auth?sslmode=disable" up
```

### Откат последней миграции

```bash
go run cmd/migrate/main.go -dsn "postgres://user:pass@localhost:5432/auth?sslmode=disable" down
```

### Проверка статуса

```bash
go run cmd/migrate/main.go -dsn "postgres://user:pass@localhost:5432/auth?sslmode=disable" status
```

### Сборка утилиты

```bash
go build -o bin/migrate ./cmd/migrate

# Использование
./bin/migrate -dsn "postgres://..." up
```

## 📝 Формат миграций

### Именование файлов

```
{version}_{description}.{direction}.sql
```

- **version**: номер версии (001, 002, 003...)
- **description**: описание на английском (snake_case)
- **direction**: `up` (применение) или `down` (откат)

### Пример

**001_create_accounts_table.up.sql:**
```sql
-- +goose up
CREATE TABLE accounts (...);
CREATE INDEX idx_accounts_email ON accounts(email);
```

**001_create_accounts_table.down.sql:**
```sql
-- +goose down
DROP TABLE IF EXISTS accounts;
DROP INDEX IF EXISTS idx_accounts_email;
```

## 🔧 Создание новой миграции

### 1. Создать файлы

```bash
# Копировать существующую миграцию как шаблон
cp migrations/001_create_accounts_table.up.sql migrations/002_new_feature.up.sql
cp migrations/001_create_accounts_table.down.sql migrations/002_new_feature.down.sql
```

### 2. Отредактировать

**002_new_feature.up.sql:**
```sql
-- +goose up
-- Описание изменений

ALTER TABLE accounts ADD COLUMN new_column VARCHAR;
CREATE INDEX idx_accounts_new_column ON accounts(new_column);
```

**002_new_feature.down.sql:**
```sql
-- +goose down
-- Откат изменений

DROP INDEX IF EXISTS idx_accounts_new_column;
ALTER TABLE accounts DROP COLUMN IF EXISTS new_column;
```

### 3. Применить

```bash
go run cmd/migrate/main.go -dsn "postgres://..." up
```

## 🏗 Миграция 001_create_accounts_table

### Что создаёт

**Таблица `accounts`:**
- `id` UUID PRIMARY KEY — уникальный идентификатор
- `email` VARCHAR(254) NOT NULL UNIQUE — email пользователя
- `password` VARCHAR(72) NOT NULL — хеш пароля (bcrypt)
- `created_at` TIMESTAMPTZ — время создания
- `updated_at` TIMESTAMPTZ — время обновления

**Индексы:**
- `idx_accounts_email` — уникальный индекс на email (O(1) поиск)
- `idx_accounts_created_at` — для сортировки по дате

**Триггер:**
- `update_updated_at_column` — автоматическое обновление `updated_at`

### Производительность

- Уникальный индекс на email — быстрый поиск при login
- Триггер — автоматическое поддержание `updated_at`

## 🐛 Troubleshooting

### Ошибка "migration already applied"

```bash
# Проверить статус
go run cmd/migrate/main.go -dsn "postgres://..." status

# Принудительно установить версию
# (осторожно, только если уверены что миграция применена)
```

### Ошибка "relation already exists"

```bash
# Миграция уже применена, проверить статус
go run cmd/migrate/main.go -dsn "postgres://..." status

# Если нужно, откатить и применить заново
go run cmd/migrate/main.go -dsn "postgres://..." down
go run cmd/migrate/main.go -dsn "postgres://..." up
```

### Грязная миграция (dirty migration)

```bash
# Если миграция прервалась mid-flight
# Проверить статус
go run cmd/migrate/main.go -dsn "postgres://..." status

# Если показывает "dirty", нужно исправить вручную в БД
psql "postgres://..." -c "SELECT * FROM schema_migrations;"
```

## 📚 Ссылки

- [golang-migrate documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL CREATE TABLE](https://www.postgresql.org/docs/current/sql-createtable.html)
- [PostgreSQL CREATE INDEX](https://www.postgresql.org/docs/current/sql-createindex.html)
