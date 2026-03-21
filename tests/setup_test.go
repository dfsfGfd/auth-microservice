// Package tests предоставляет e2e тесты для auth-microservice.
package tests

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const (
	// Docker compose файл для интеграционных тестов
	dockerComposeFile = "docker-compose.integration.yml"

	// Таймауты
	composeUpTimeout    = 2 * time.Minute
	serviceReadyTimeout = 30 * time.Second
	composeDownTimeout  = 30 * time.Second

	// Адреса сервисов для тестов
	grpcAddress = "localhost:9091"
	httpAddress = "localhost:8081"
)

// TestMain точка входа для всех интеграционных тестов
func TestMain(m *testing.M) {
	ctx := context.Background()

	// Получаем путь к директории с docker-compose файлом
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get working directory: %v\n", err)
		os.Exit(1)
	}

	composePath := filepath.Join(wd, dockerComposeFile)

	// Setup: поднимаем контейнеры
	fmt.Println("🚀 Starting integration test infrastructure...")
	if err := composeUp(ctx, composePath); err != nil {
		fmt.Printf("❌ Failed to start containers: %v\n", err)
		os.Exit(1)
	}

	// Ждём готовности сервисов
	fmt.Println("⏳ Waiting for services to be ready...")
	if err := waitForServices(ctx); err != nil {
		fmt.Printf("❌ Services not ready: %v\n", err)
		_ = composeDown(ctx, composePath)
		os.Exit(1)
	}
	fmt.Println("✅ Services are ready!")

	// Запускаем тесты
	exitCode := m.Run()

	// Teardown: удаляем контейнеры
	fmt.Println("\n🧹 Cleaning up test infrastructure...")
	if err := composeDown(ctx, composePath); err != nil {
		fmt.Printf("⚠️  Warning: failed to cleanup containers: %v\n", err)
	}
	fmt.Println("✅ Cleanup complete")

	os.Exit(exitCode)
}

// composeUp поднимает контейнеры через docker-compose
func composeUp(ctx context.Context, composePath string) error {
	ctx, cancel := context.WithTimeout(ctx, composeUpTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "up", "-d", "--build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// composeDown удаляет контейнеры
func composeDown(ctx context.Context, composePath string) error {
	ctx, cancel := context.WithTimeout(ctx, composeDownTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composePath, "down", "-v", "--remove-orphans")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// waitForServices ждёт готовности всех сервисов
func waitForServices(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, serviceReadyTimeout)
	defer cancel()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for services")
		case <-ticker.C:
			if err := checkServices(ctx); err == nil {
				return nil
			}
		}
	}
}

// checkServices проверяет доступность сервисов
func checkServices(ctx context.Context) error {
	// Проверяем HTTP health endpoint
	cmd := exec.CommandContext(ctx, "curl", "-s", "-o", "/dev/null", "-w", "%{http_code}",
		fmt.Sprintf("http://%s/health", httpAddress))

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if string(output) != "200" {
		return fmt.Errorf("health check returned %s", string(output))
	}

	return nil
}

// getGRPCAddress возвращает адрес gRPC сервера
func getGRPCAddress() string {
	return grpcAddress
}

// getHTTPAddress возвращает адрес HTTP сервера
func getHTTPAddress() string {
	return httpAddress
}
