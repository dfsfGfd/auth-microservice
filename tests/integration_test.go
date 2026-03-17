// Package integration предоставляет интеграционные тесты для auth-microservice.
//
// Тесты используют testcontainers для запуска PostgreSQL и Redis в Docker.
//
// Запуск тестов:
//
//	go test ./tests -v
//
// Запуск с очисткой контейнеров после тестов:
//
//	go test ./tests -v -cleanup
package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	// Тестовые константы
	testEmail    = "test@example.com"
	testPassword = "TestPassword123!"
	testUsername = "testuser"

	// Таймауты
	startupTimeout   = 60 * time.Second
	shutdownTimeout  = 30 * time.Second
	requestTimeout   = 10 * time.Second
	testTimeout      = 2 * time.Minute
)

// TestMain запускается перед всеми тестами.
// Создаёт общие тестовые контейнеры PostgreSQL и Redis.
func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Запуск контейнеров
	psqlContainer, psqlURL, err := setupPostgres(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup PostgreSQL: %v\n", err)
		os.Exit(1)
	}

	redisContainer, redisURL, err := setupRedis(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to setup Redis: %v\n", err)
		os.Exit(1)
	}

	// Установка переменных окружения для тестов
	os.Setenv("DATABASE_URL", psqlURL)
	os.Setenv("REDIS_URL", redisURL)
	os.Setenv("JWT_SECRET", "test-secret-key-minimum-32-characters-long")
	os.Setenv("APP_ENV", "test")

	// Применение миграций
	if err := runMigrations(ctx, psqlURL); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run migrations: %v\n", err)
		cleanupContainers(ctx, psqlContainer, redisContainer)
		os.Exit(1)
	}

	// Запуск тестов
	exitCode := m.Run()

	// Очистка
	cleanupContainers(ctx, psqlContainer, redisContainer)
	os.Exit(exitCode)
}

// setupPostgres создаёт контейнер PostgreSQL для тестов.
func setupPostgres(ctx context.Context) (testcontainers.Container, string, error) {
	fmt.Println("🐘 Starting PostgreSQL container...")

	container, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("auth_test"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(startupTimeout),
		),
	)
	if err != nil {
		return nil, "", fmt.Errorf("create postgres container: %w", err)
	}

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, "", fmt.Errorf("get connection string: %w", err)
	}

	fmt.Printf("✅ PostgreSQL ready: %s\n", connectionString)
	return container, connectionString, nil
}

// setupRedis создаёт контейнер Redis для тестов.
func setupRedis(ctx context.Context) (testcontainers.Container, string, error) {
	fmt.Println("📦 Starting Redis container...")

	container, err := redis.Run(ctx,
		"redis:7.4-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(startupTimeout),
		),
	)
	if err != nil {
		return nil, "", fmt.Errorf("create redis container: %w", err)
	}

	connectionString, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("get connection string: %w", err)
	}

	fmt.Printf("✅ Redis ready: %s\n", connectionString)
	return container, connectionString, nil
}

// runMigrations применяет миграции к тестовой БД.
func runMigrations(ctx context.Context, dsn string) error {
	fmt.Println("🔄 Running migrations...")

	// Получаем корень проекта (родитель tests/)
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// Используем migrate утилиту
	cmd := exec.CommandContext(ctx, "go", "run",
		"../cmd/migrate/main.go",
		"-dsn", dsn,
		"-path", "../migrations",
		"up",
	)
	cmd.Dir = wd // Устанавливаем рабочую директорию

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("run migrations: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✅ Migrations applied: %s\n", string(output))
	return nil
}

// cleanupContainers останавливает и удаляет контейнеры.
func cleanupContainers(ctx context.Context, containers ...testcontainers.Container) {
	fmt.Println("🧹 Cleaning up containers...")

	for i, container := range containers {
		if container == nil {
			continue
		}

		ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()

		if err := container.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to terminate container %d: %v\n", i, err)
		}
	}

	fmt.Println("✅ Cleanup complete")
}

// ============================================================================
// Интеграционные тесты
// ============================================================================

// TestIntegration_Registration тестирует процесс регистрации пользователя.
func TestIntegration_Registration(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	t.Run("successful_registration", func(t *testing.T) {
		// Тест будет реализован после создания service layer
		t.Skip("Registration test will be implemented after service layer")
	})

	t.Run("duplicate_email", func(t *testing.T) {
		t.Skip("Duplicate email test will be implemented after service layer")
	})

	t.Run("invalid_email", func(t *testing.T) {
		t.Skip("Invalid email test will be implemented after service layer")
	})

	t.Run("weak_password", func(t *testing.T) {
		t.Skip("Weak password test will be implemented after service layer")
	})

	_ = ctx // Используется в тестах
}

// TestIntegration_Login тестирует процесс входа пользователя.
func TestIntegration_Login(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	t.Run("successful_login", func(t *testing.T) {
		t.Skip("Login test will be implemented after service layer")
	})

	t.Run("invalid_credentials", func(t *testing.T) {
		t.Skip("Invalid credentials test will be implemented after service layer")
	})

	t.Run("account_not_found", func(t *testing.T) {
		t.Skip("Account not found test will be implemented after service layer")
	})

	_ = ctx
}

// TestIntegration_TokenRefresh тестирует обновление токенов.
func TestIntegration_TokenRefresh(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	t.Run("successful_refresh", func(t *testing.T) {
		t.Skip("Refresh test will be implemented after service layer")
	})

	t.Run("invalid_refresh_token", func(t *testing.T) {
		t.Skip("Invalid refresh token test will be implemented after service layer")
	})

	t.Run("expired_refresh_token", func(t *testing.T) {
		t.Skip("Expired refresh token test will be implemented after service layer")
	})

	_ = ctx
}

// TestIntegration_Logout тестирует процесс выхода.
func TestIntegration_Logout(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	t.Run("successful_logout", func(t *testing.T) {
		t.Skip("Logout test will be implemented after service layer")
	})

	t.Run("logout_invalid_token", func(t *testing.T) {
		t.Skip("Invalid logout token test will be implemented after service layer")
	})

	_ = ctx
}

// TestIntegration_Database тестирует работу с базой данных.
func TestIntegration_Database(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	dsn := os.Getenv("DATABASE_URL")
	require.NotEmpty(t, dsn, "DATABASE_URL should be set")

	t.Run("database_connection", func(t *testing.T) {
		// Проверка подключения к БД будет реализована
		t.Skip("Database connection test will be implemented")
	})

	t.Run("concurrent_operations", func(t *testing.T) {
		// Тест конкурентных операций
		t.Skip("Concurrent operations test will be implemented")
	})

	_ = ctx
}

// TestIntegration_Redis тестирует работу с Redis.
func TestIntegration_Redis(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	dsn := os.Getenv("REDIS_URL")
	require.NotEmpty(t, dsn, "REDIS_URL should be set")

	t.Run("redis_connection", func(t *testing.T) {
		// Проверка подключения к Redis будет реализована
		t.Skip("Redis connection test will be implemented")
	})

	t.Run("token_caching", func(t *testing.T) {
		// Тест кэширования токенов
		t.Skip("Token caching test will be implemented")
	})

	_ = ctx
}

// TestIntegration_RateLimiting тестирует rate limiting.
func TestIntegration_RateLimiting(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	t.Run("rate_limit_login", func(t *testing.T) {
		// Тест ограничения login запросов
		t.Skip("Rate limit login test will be implemented")
	})

	t.Run("rate_limit_register", func(t *testing.T) {
		// Тест ограничения register запросов
		t.Skip("Rate limit register test will be implemented")
	})

	_ = ctx
}

// ============================================================================
// Вспомогательные функции
// ============================================================================
