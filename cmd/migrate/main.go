// Package main предоставляет утилиту для управления миграциями БД.
//
// Использование:
//
//	# Применить все миграции
//	migrate -dsn postgres://user:pass@localhost:5432/db?sslmode=disable up
//
//	# Откатить последнюю миграцию
//	migrate -dsn postgres://user:pass@localhost:5432/db?sslmode=disable down
//
//	# Проверить статус
//	migrate -dsn postgres://user:pass@localhost:5432/db?sslmode=disable status
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //nolint:revive
	_ "github.com/golang-migrate/migrate/v4/source/file"       //nolint:revive
)

const (
	exitSuccess     = 0
	exitError       = 1
	exitUsageError  = 2
)

func main() {
	os.Exit(run())
}

func run() int {
	// Флаги командной строки
	dsn := flag.String("dsn", "", "Database URL (обязательно)")
	path := flag.String("path", "migrations", "Путь к миграциям")
	version := flag.Bool("version", false, "Показать версию утилиты")
	help := flag.Bool("help", false, "Показать помощь")
	flag.Usage = usage
	flag.Parse()

	// Показать версию
	if *version {
		fmt.Println("migrate version 1.0.0")
		return exitSuccess
	}

	// Показать помощь
	if *help {
		usage()
		return exitSuccess
	}

	// Проверка DSN
	if *dsn == "" {
		fmt.Fprintln(os.Stderr, "Ошибка: требуется флаг -dsn")
		usage()
		return exitUsageError
	}

	// Получение команды (up, down, status)
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Ошибка: требуется команда (up, down, status)")
		usage()
		return exitUsageError
	}

	command := flag.Arg(0)

	// Создание миграции
	m, err := migrate.New(
		fmt.Sprintf("file://%s", *path),
		*dsn,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания миграции: %v\n", err)
		return exitError
	}
	defer m.Close()

	// Выполнение команды
	switch command {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("Нет новых миграций для применения")
				return exitSuccess
			}
			fmt.Fprintf(os.Stderr, "Ошибка применения миграций: %v\n", err)
			return exitError
		}
		fmt.Println("Миграции успешно применены")

	case "down":
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("Нет миграций для отката")
				return exitSuccess
			}
			fmt.Fprintf(os.Stderr, "Ошибка отката миграций: %v\n", err)
			return exitError
		}
		fmt.Println("Миграции успешно откатаны")

	case "status":
		version, dirty, err := m.Version()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка получения статуса: %v\n", err)
			return exitError
		}
		if dirty {
			fmt.Printf("Версия: %d (грязная)\n", version)
		} else {
			fmt.Printf("Версия: %d\n", version)
		}

	default:
		fmt.Fprintf(os.Stderr, "Неизвестная команда: %s\n", command)
		fmt.Fprintln(os.Stderr, "Доступные команды: up, down, status")
		return exitUsageError
	}

	return exitSuccess
}

func usage() {
	fmt.Fprintf(os.Stderr, `Утилита для управления миграциями базы данных.

Использование:
  migrate -dsn <url> <command>

Команды:
  up       Применить все миграции
  down     Откатить последнюю миграцию
  status   Показать текущую версию

Флаги:
  -dsn     Database URL (обязательно)
  -path    Путь к миграциям (по умолчанию: migrations)
  -version Показать версию утилиты
  -help    Показать эту справку

Примеры:
  migrate -dsn postgres://user:pass@localhost:5432/auth?sslmode=disable up
  migrate -dsn postgres://user:pass@localhost:5432/auth?sslmode=disable down
  migrate -dsn postgres://user:pass@localhost:5432/auth?sslmode=disable status
`)
}
