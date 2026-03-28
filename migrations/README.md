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
# Установить CLI (требуется один раз)
task migrate:install

# Применить все миграции
task migrate:up

# Откатить последнюю миграцию
task migrate:down

# Показать статус миграций
task migrate:status

# Принудительно установить версию (если нужно)
task migrate:force VERSION=3
```

### Напрямую через CLI

```bash
# Установить golang-migrate
go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.1

# Применить
migrate -path migrations -database "postgres://user:pass@localhost:5432/auth?sslmode=disable" up

# Откатить
migrate -path migrations -database "postgres://..." down

# Статус
migrate -path migrations -database "postgres://..." status

# Force (установить версию)
migrate -path migrations -database "postgres://..." force 3
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
| `email` | VARCHAR(254) NOT NULL | Email пользователя (уникальный) |
| `password` | VARCHAR(72) NOT NULL | Хеш пароля (bcrypt) |
| `created_at` | TIMESTAMPTZ | Время создания |
| `updated_at` | TIMESTAMPTZ | Время обновления |

**Индексы:**
- `idx_accounts_email` — уникальный индекс для быстрого поиска по email

---

## 📚 Ссылки

- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [Основная документация](../docs/README.md)
