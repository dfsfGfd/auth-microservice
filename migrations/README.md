# Migrations Guide

Управление миграциями PostgreSQL с использованием golang-migrate.

---

## 📁 Структура

```
migrations/
├── 001_create_accounts_table.up.sql   # Применение
└── 001_create_accounts_table.down.sql # Откат
```

---

## 🚀 Использование

```bash
# Применить все миграции
go run cmd/migrate/main.go -dsn "postgres://user:pass@localhost:5432/auth?sslmode=disable" up

# Откатить последнюю
go run cmd/migrate/main.go -dsn "postgres://..." down

# Статус
go run cmd/migrate/main.go -dsn "postgres://..." status
```

### Сборка утилиты

```bash
go build -o bin/migrate ./cmd/migrate
./bin/migrate -dsn "postgres://..." up
```

---

## 📝 Формат файлов

```
{version}_{description}.{direction}.sql
```

- **version:** 001, 002, 003...
- **direction:** `up` или `down`

### Пример

**001_create_accounts_table.up.sql:**
```sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    email VARCHAR(254) NOT NULL UNIQUE,
    password VARCHAR(72) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_accounts_email ON accounts(email);
```

**001_create_accounts_table.down.sql:**
```sql
DROP TABLE IF EXISTS accounts;
```

---

## 🔧 Создание новой миграции

```bash
# Копировать шаблон
cp migrations/001_create_accounts_table.up.sql migrations/002_new_feature.up.sql
cp migrations/001_create_accounts_table.down.sql migrations/002_new_feature.down.sql

# Отредактировать файлы

# Применить
go run cmd/migrate/main.go -dsn "postgres://..." up
```

---

## 🏗 Миграция 001

**Таблица `accounts`:**
- `id` UUID PRIMARY KEY
- `email` VARCHAR(254) UNIQUE
- `password` VARCHAR(72) (bcrypt hash)
- `created_at` TIMESTAMPTZ
- `updated_at` TIMESTAMPTZ

**Индексы:**
- `idx_accounts_email` — для быстрого поиска по email

---

## 📚 Ссылки

- [golang-migrate](https://github.com/golang-migrate/migrate)
- [README](docs/README.md)
