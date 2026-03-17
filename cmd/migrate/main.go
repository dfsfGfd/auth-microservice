// Package main предоставляет утилиту для управления миграциями БД.
//
// Использование:
//
//	# Применить все миграции
//	migrate up
//
//	# Откатить последнюю миграцию
//	migrate down
//
//	# Проверить статус
//	migrate status
//
//	# Применить до конкретной версии
//	migrate up 1
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Флаги командной строки
	dsn := flag.String("dsn", "", "Database URL (обязательно)")
	path := flag.String("path", "migrations", "Путь к миграциям")
	flag.Parse()

	// Проверка DSN
	if *dsn == "" {
		fmt.Fprintln(os.Stderr, "Ошибка: требуется флаг -dsn")
		fmt.Fprintln(os.Stderr, "Пример: migrate -dsn postgres://user:pass@localhost:5432/db?sslmode=disable up")
		os.Exit(1)
	}

	// Получение команды (up, down, status)
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Ошибка: требуется команда (up, down, status)")
		fmt.Fprintln(os.Stderr, "Пример: migrate -dsn <url> up")
		os.Exit(1)
	}

	command := flag.Arg(0)

	// Создание миграции
	m, err := migrate.New(
		fmt.Sprintf("file://%s", *path),
		*dsn,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания миграции: %v\n", err)
		os.Exit(1)
	}
	defer m.Close()

	// Выполнение команды
	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Fprintf(os.Stderr, "Ошибка применения миграций: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Миграции успешно применены")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Fprintf(os.Stderr, "Ошибка отката миграций: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Миграции успешно откатаны")

	case "status":
		version, dirty, err := m.Version()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка получения статуса: %v\n", err)
			os.Exit(1)
		}
		if dirty {
			fmt.Printf("Версия: %d (грязная)\n", version)
		} else {
			fmt.Printf("Версия: %d\n", version)
		}

	default:
		fmt.Fprintf(os.Stderr, "Неизвестная команда: %s\n", command)
		fmt.Fprintln(os.Stderr, "Доступные команды: up, down, status")
		os.Exit(1)
	}
}
