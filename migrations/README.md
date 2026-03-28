# Migrations Guide

Управление миграциями PostgreSQL с использованием golang-migrate.

---

## 📁 Структура

```
migrations/
├── 00001_create_accounts_table.up.sql   # Применение
└── 00001_create_accounts_table.down.sql # Откат
```

**Формат имён:** `{5-значная версия}_{description}.{direction}.sql`

- **version:** 00001, 00002, 00003...
- **direction:** `up` или `down`

---

## 🚀 Использование

### Через Taskfile (рекомендуется)

```bash
# Применить все миграции в Docker
task migrate:up

# Откатить последнюю миграцию
task migrate:down

# Показать статус миграций
task migrate:status
```

### Локально (требуется PostgreSQL)

```bash
# Применить
go run cmd/migrate/main.go -dsn "postgres://user:pass@localhost:5432/auth?sslmode=disable" -path migrations up

# Откатить
go run cmd/migrate/main.go -dsn "postgres://..." -path migrations down

# Статус
go run cmd/migrate/main.go -dsn "postgres://..." -path migrations status
```

### Сборка утилиты

```bash
go build -o bin/migrate ./cmd/migrate
./bin/migrate -dsn "postgres://..." -path migrations up
```

---

## 🔧 Создание новой миграции

1. **Создать файлы миграции:**

```bash
cp migrations/00001_create_accounts_table.up.sql migrations/00002_new_feature.up.sql
cp migrations/00001_create_accounts_table.down.sql migrations/00002_new_feature.down.sql
```

2. **Отредактировать файлы** под новую функциональность

3. **Применить миграцию:**

```bash
task migrate:up
```

---

## 🏗 Миграция 00001

**Таблица `accounts`:**

| Колонка | Тип | Описание |
|---------|-----|----------|
| `id` | UUID PRIMARY KEY | Уникальный идентификатор |
| `email` | VARCHAR(254) NOT NULL | Email пользователя |
| `password` | VARCHAR(72) NOT NULL | Хеш пароля (bcrypt) |
| `created_at` | TIMESTAMPTZ | Время создания |
| `updated_at` | TIMESTAMPTZ | Время обновления |

**Индексы:**
- `idx_accounts_email` — уникальный индекс для быстрого поиска по email

**Триггер:**
- `update_accounts_updated_at` — автоматически обновляет `updated_at` при изменении записи

---

## 📚 Ссылки

- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Основная документация](../docs/README.md)
