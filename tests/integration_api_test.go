//go:build integration
// +build integration

// Package integration предоставляет интеграционные тесты для auth-microservice.
//
// Запуск тестов:
//
//	# 1. Поднять инфраструктуру
//	docker-compose -f docker-compose.integration.yml up -d
//
//	# 2. Запустить тесты
//	go test -tags integration -v
//
//	# 3. Очистить
//	docker-compose -f docker-compose.integration.yml down -v
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Базовый URL API (из docker-compose)
	baseURL = "http://localhost:8080"

	// Endpoints
	registerEndpoint = "/api/auth/register"
	loginEndpoint    = "/api/auth/login"
	logoutEndpoint   = "/api/auth/logout"
	refreshEndpoint  = "/api/auth/refresh"
	healthEndpoint   = "/health"

	// Тестовые данные
	testEmail           = "test@example.com"
	testPassword        = "TestPassword123!"
	testPasswordInvalid = "wrongpassword"

	// Таймауты
	requestTimeout   = 30 * time.Second
	testTimeout      = 5 * time.Minute
	shutdownTimeout  = 10 * time.Second
)

// TestMain запускается перед всеми тестами.
// Проверяет, что сервис готов к тестам.
func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Ожидание готовности сервиса
	if err := waitForService(ctx, baseURL+healthEndpoint); err != nil {
		fmt.Fprintf(os.Stderr, "Service not ready: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Service is ready for integration tests")

	// Запуск тестов
	exitCode := m.Run()

	// Очистка (опционально)
	cleanupTestData(ctx)

	os.Exit(exitCode)
}

// waitForService ожидает доступности сервиса.
func waitForService(ctx context.Context, url string) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(60 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for service")
		case <-ticker.C:
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				return nil
			}
			fmt.Println("⏳ Waiting for service to be ready...")
		}
	}
}

// cleanupTestData очищает тестовые данные после тестов.
func cleanupTestData(ctx context.Context) {
	fmt.Println("🧹 Cleaning up test data...")
	// Очистка будет реализована через прямые запросы к БД
}

// ============================================================================
// Health Check тесты
// ============================================================================

// TestIntegration_HealthCheck проверяет health check endpoint.
func TestIntegration_HealthCheck(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+healthEndpoint, nil)
	require.NoError(t, err)

	client := &http.Client{Timeout: requestTimeout}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "ok", result["status"])
	t.Log("✅ Health check passed")
}

// ============================================================================
// Registration тесты
// ============================================================================

// TestIntegration_Registration_Success тестирует успешную регистрацию.
func TestIntegration_Registration_Success(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Уникальный email для каждого теста
	email := fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())

	payload := map[string]string{
		"email":    email,
		"password": testPassword,
	}

	resp := makeRequest(ctx, t, http.MethodPost, baseURL+registerEndpoint, payload)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, float64(200), result["status_code"])
	assert.Contains(t, result["message"], "registered")
	assert.NotNil(t, result["data"])

	t.Logf("✅ Registration successful: %s", email)
}

// TestIntegration_Registration_DuplicateEmail тестирует регистрацию с существующим email.
func TestIntegration_Registration_DuplicateEmail(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	payload := map[string]string{
		"email":    testEmail,
		"password": testPassword,
	}

	// Первая регистрация (успешная)
	resp1 := makeRequest(ctx, t, http.MethodPost, baseURL+registerEndpoint, payload)
	resp1.Body.Close()

	// Вторая регистрация (должна вернуть ошибку)
	resp2 := makeRequestWithStatus(ctx, t, http.MethodPost, baseURL+registerEndpoint, payload, http.StatusConflict)
	defer resp2.Body.Close()

	assert.Equal(t, http.StatusConflict, resp2.StatusCode)
	t.Log("✅ Duplicate email registration correctly rejected")
}

// TestIntegration_Registration_InvalidEmail тестирует регистрацию с невалидным email.
func TestIntegration_Registration_InvalidEmail(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	payload := map[string]string{
		"email":    "invalid-email",
		"password": testPassword,
	}

	resp := makeRequestWithStatus(ctx, t, http.MethodPost, baseURL+registerEndpoint, payload, http.StatusBadRequest)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	t.Log("✅ Invalid email registration correctly rejected")
}

// TestIntegration_Registration_WeakPassword тестирует регистрацию со слабым паролем.
func TestIntegration_Registration_WeakPassword(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	payload := map[string]string{
		"email":    fmt.Sprintf("test_weak_%d@example.com", time.Now().UnixNano()),
		"password": "123", // Слишком короткий
	}

	resp := makeRequestWithStatus(ctx, t, http.MethodPost, baseURL+registerEndpoint, payload, http.StatusBadRequest)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	t.Log("✅ Weak password registration correctly rejected")
}

// ============================================================================
// Login тесты
// ============================================================================

// TestIntegration_Login_Success тестирует успешный вход.
func TestIntegration_Login_Success(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Сначала регистрируемся
	email := fmt.Sprintf("login_test_%d@example.com", time.Now().UnixNano())
	registerPayload := map[string]string{
		"email":    email,
		"password": testPassword,
	}

	registerResp := makeRequest(ctx, t, http.MethodPost, baseURL+registerEndpoint, registerPayload)
	registerResp.Body.Close()

	// Теперь входим
	loginPayload := map[string]string{
		"email":    email,
		"password": testPassword,
	}

	resp := makeRequest(ctx, t, http.MethodPost, baseURL+loginEndpoint, loginPayload)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, float64(200), result["status_code"])
	assert.NotNil(t, result["data"])

	data := result["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])
	assert.Equal(t, "Bearer", data["token_type"])

	t.Log("✅ Login successful")
}

// TestIntegration_Login_InvalidCredentials тестирует вход с неверными данными.
func TestIntegration_Login_InvalidCredentials(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	payload := map[string]string{
		"email":    "nonexistent@example.com",
		"password": testPasswordInvalid,
	}

	resp := makeRequestWithStatus(ctx, t, http.MethodPost, baseURL+loginEndpoint, payload, http.StatusUnauthorized)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	t.Log("✅ Invalid credentials login correctly rejected")
}

// ============================================================================
// Token Refresh тесты
// ============================================================================

// TestIntegration_TokenRefresh_Success тестирует успешное обновление токена.
func TestIntegration_TokenRefresh_Success(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Регистрируемся и входим
	email := fmt.Sprintf("refresh_test_%d@example.com", time.Now().UnixNano())

	registerPayload := map[string]string{
		"email":    email,
		"password": testPassword,
	}
	registerResp := makeRequest(ctx, t, http.MethodPost, baseURL+registerEndpoint, registerPayload)
	registerResp.Body.Close()

	loginResp := makeRequest(ctx, t, http.MethodPost, baseURL+loginEndpoint, map[string]string{
		"email":    email,
		"password": testPassword,
	})

	var loginResult map[string]interface{}
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	loginResp.Body.Close()

	loginData := loginResult["data"].(map[string]interface{})
	refreshToken := loginData["refresh_token"].(string)

	// Обновляем токен
	refreshPayload := map[string]string{
		"refresh_token": refreshToken,
	}

	resp := makeRequest(ctx, t, http.MethodPost, baseURL+refreshEndpoint, refreshPayload)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, float64(200), result["status_code"])
	assert.NotNil(t, result["data"])

	data := result["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["refresh_token"])

	t.Log("✅ Token refresh successful")
}

// ============================================================================
// Logout тесты
// ============================================================================

// TestIntegration_Logout_Success тестирует успешный выход.
func TestIntegration_Logout_Success(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	// Регистрируемся и входим
	email := fmt.Sprintf("logout_test_%d@example.com", time.Now().UnixNano())

	registerPayload := map[string]string{
		"email":    email,
		"password": testPassword,
	}
	registerResp := makeRequest(ctx, t, http.MethodPost, baseURL+registerEndpoint, registerPayload)
	registerResp.Body.Close()

	loginResp := makeRequest(ctx, t, http.MethodPost, baseURL+loginEndpoint, map[string]string{
		"email":    email,
		"password": testPassword,
	})

	var loginResult map[string]interface{}
	json.NewDecoder(loginResp.Body).Decode(&loginResult)
	loginResp.Body.Close()

	loginData := loginResult["data"].(map[string]interface{})
	refreshToken := loginData["refresh_token"].(string)

	// Выходим
	logoutPayload := map[string]string{
		"refresh_token": refreshToken,
	}

	resp := makeRequest(ctx, t, http.MethodPost, baseURL+logoutEndpoint, logoutPayload)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	t.Log("✅ Logout successful")
}

// ============================================================================
// Вспомогательные функции
// ============================================================================

// makeRequest отправляет HTTP запрос и проверяет статус 200.
func makeRequest(ctx context.Context, t *testing.T, method, url string, payload interface{}) *http.Response {
	return makeRequestWithStatus(ctx, t, method, url, payload, http.StatusOK)
}

// makeRequestWithStatus отправляет HTTP запрос и проверяет ожидаемый статус.
func makeRequestWithStatus(ctx context.Context, t *testing.T, method, url string, payload interface{}, expectedStatus int) *http.Response {
	var body []byte
	var err error

	if payload != nil {
		body, err = json.Marshal(payload)
		require.NoError(t, err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: requestTimeout}
	resp, err := client.Do(req)
	require.NoError(t, err)

	if expectedStatus > 0 {
		assert.Equal(t, expectedStatus, resp.StatusCode, "Unexpected status code")
	}

	return resp
}
