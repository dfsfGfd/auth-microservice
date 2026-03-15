package cookies_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"auth-microservice/pkg/cookies"
)

func TestNewService(t *testing.T) {
	t.Run("создание сервиса с полной конфигурацией", func(t *testing.T) {
		svc := cookies.NewService(cookies.Config{
			Secure:   true,
			HTTPOnly: true,
			SameSite: http.SameSiteStrictMode,
			Domain:   "example.com",
			Path:     "/api",
			MaxAge:   3600,
		})

		require.NotNil(t, svc)
		config := svc.GetConfig()
		assert.Equal(t, true, config.Secure)
		assert.Equal(t, true, config.HTTPOnly)
		assert.Equal(t, http.SameSiteStrictMode, config.SameSite)
		assert.Equal(t, "example.com", config.Domain)
		assert.Equal(t, "/api", config.Path)
		assert.Equal(t, 3600, config.MaxAge)
	})

	t.Run("создание сервиса с конфигурацией по умолчанию", func(t *testing.T) {
		svc := cookies.NewService(cookies.Config{})

		require.NotNil(t, svc)
		config := svc.GetConfig()
		assert.Equal(t, "/", config.Path)
		assert.Equal(t, int((14*24*time.Hour).Seconds()), config.MaxAge)
	})
}

func TestService_SetRefreshToken(t *testing.T) {
	svc := cookies.NewService(cookies.Config{
		Secure:   true,
		HTTPOnly: true,
		SameSite: http.SameSiteStrictMode,
		Domain:   "example.com",
		MaxAge:   1209600, // 14 дней
	})

	t.Run("установка refresh токена", func(t *testing.T) {
		w := httptest.NewRecorder()
		token := "test_refresh_token_12345"

		svc.SetRefreshToken(w, token)

		result := w.Result()
		cookies := result.Cookies()

		require.Len(t, cookies, 1)
		cookie := cookies[0]

		assert.Equal(t, "refresh_token", cookie.Name)
		assert.Equal(t, token, cookie.Value)
		assert.Equal(t, 1209600, cookie.MaxAge)
		assert.Equal(t, "/", cookie.Path)
		assert.Equal(t, "example.com", cookie.Domain)
		assert.Equal(t, true, cookie.Secure)
		assert.Equal(t, true, cookie.HttpOnly)
		assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
	})
}

func TestService_GetRefreshToken(t *testing.T) {
	svc := cookies.NewService(cookies.Config{})

	t.Run("получение refresh токена", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		r.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: "test_token_value",
		})

		token, err := svc.GetRefreshToken(r)

		require.NoError(t, err)
		assert.Equal(t, "test_token_value", token)
	})

	t.Run("отсутствие refresh токена", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)

		token, err := svc.GetRefreshToken(r)

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestService_DeleteRefreshToken(t *testing.T) {
	svc := cookies.NewService(cookies.Config{
		Secure:   true,
		HTTPOnly: true,
	})

	t.Run("удаление refresh токена", func(t *testing.T) {
		w := httptest.NewRecorder()

		svc.DeleteRefreshToken(w)

		result := w.Result()
		cookies := result.Cookies()

		require.Len(t, cookies, 1)
		cookie := cookies[0]

		assert.Equal(t, "refresh_token", cookie.Name)
		assert.Equal(t, "", cookie.Value)
		assert.Equal(t, -1, cookie.MaxAge)
	})
}

func TestService_SetAccessToken(t *testing.T) {
	svc := cookies.NewService(cookies.Config{
		Secure:   true,
		HTTPOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	t.Run("установка access токена", func(t *testing.T) {
		w := httptest.NewRecorder()
		token := "test_access_token_67890"
		expiresIn := 900 // 15 минут

		svc.SetAccessToken(w, token, expiresIn)

		result := w.Result()
		cookies := result.Cookies()

		require.Len(t, cookies, 1)
		cookie := cookies[0]

		assert.Equal(t, "access_token", cookie.Name)
		assert.Equal(t, token, cookie.Value)
		assert.Equal(t, expiresIn, cookie.MaxAge)
		assert.Equal(t, true, cookie.Secure)
		assert.Equal(t, true, cookie.HttpOnly)
		assert.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
	})
}

func TestService_GetAccessToken(t *testing.T) {
	svc := cookies.NewService(cookies.Config{})

	t.Run("получение access токена", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/resource", nil)
		r.AddCookie(&http.Cookie{
			Name:  "access_token",
			Value: "access_token_value",
		})

		token, err := svc.GetAccessToken(r)

		require.NoError(t, err)
		assert.Equal(t, "access_token_value", token)
	})

	t.Run("отсутствие access токена", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/resource", nil)

		token, err := svc.GetAccessToken(r)

		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestService_DeleteAccessToken(t *testing.T) {
	svc := cookies.NewService(cookies.Config{})

	t.Run("удаление access токена", func(t *testing.T) {
		w := httptest.NewRecorder()

		svc.DeleteAccessToken(w)

		result := w.Result()
		cookies := result.Cookies()

		require.Len(t, cookies, 1)
		cookie := cookies[0]

		assert.Equal(t, "access_token", cookie.Name)
		assert.Equal(t, "", cookie.Value)
		assert.Equal(t, -1, cookie.MaxAge)
	})
}

func TestService_DeleteAll(t *testing.T) {
	svc := cookies.NewService(cookies.Config{})

	t.Run("удаление всех cookie", func(t *testing.T) {
		w := httptest.NewRecorder()

		svc.DeleteAll(w)

		result := w.Result()
		cookies := result.Cookies()

		require.Len(t, cookies, 2)

		// Проверяем refresh_token
		var refreshToken, accessToken *http.Cookie
		for _, c := range cookies {
			if c.Name == "refresh_token" {
				refreshToken = c
			}
			if c.Name == "access_token" {
				accessToken = c
			}
		}

		require.NotNil(t, refreshToken)
		assert.Equal(t, "", refreshToken.Value)
		assert.Equal(t, -1, refreshToken.MaxAge)

		require.NotNil(t, accessToken)
		assert.Equal(t, "", accessToken.Value)
		assert.Equal(t, -1, accessToken.MaxAge)
	})
}

func TestService_ConfigValidation(t *testing.T) {
	t.Run("Path по умолчанию", func(t *testing.T) {
		svc := cookies.NewService(cookies.Config{})
		config := svc.GetConfig()
		assert.Equal(t, "/", config.Path)
	})

	t.Run("MaxAge по умолчанию (14 дней)", func(t *testing.T) {
		svc := cookies.NewService(cookies.Config{})
		config := svc.GetConfig()
		expectedMaxAge := int((14 * 24 * time.Hour).Seconds())
		assert.Equal(t, expectedMaxAge, config.MaxAge)
	})

	t.Run("кастомный Path", func(t *testing.T) {
		svc := cookies.NewService(cookies.Config{Path: "/api/auth"})
		config := svc.GetConfig()
		assert.Equal(t, "/api/auth", config.Path)
	})

	t.Run("кастомный MaxAge", func(t *testing.T) {
		svc := cookies.NewService(cookies.Config{MaxAge: 7200})
		config := svc.GetConfig()
		assert.Equal(t, 7200, config.MaxAge)
	})
}

func TestService_SecureConfig(t *testing.T) {
	t.Run("продакшен конфигурация", func(t *testing.T) {
		svc := cookies.NewService(cookies.Config{
			Secure:   true,
			HTTPOnly: true,
			SameSite: http.SameSiteStrictMode,
			Domain:   "auth.example.com",
			Path:     "/",
			MaxAge:   1209600,
		})

		config := svc.GetConfig()
		assert.True(t, config.Secure, "Secure должен быть включён")
		assert.True(t, config.HTTPOnly, "HttpOnly должен быть включён")
		assert.Equal(t, http.SameSiteStrictMode, config.SameSite)
	})

	t.Run("development конфигурация", func(t *testing.T) {
		svc := cookies.NewService(cookies.Config{
			Secure:   false,
			HTTPOnly: true,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		})

		config := svc.GetConfig()
		assert.False(t, config.Secure, "Secure может быть выключен в dev")
		assert.True(t, config.HTTPOnly, "HttpOnly должен быть включён")
		assert.Equal(t, http.SameSiteLaxMode, config.SameSite)
	})
}

func TestService_CookieFlow(t *testing.T) {
	svc := cookies.NewService(cookies.Config{
		Secure:   true,
		HTTPOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	t.Run("полный цикл: установка -> получение -> удаление", func(t *testing.T) {
		// Установка
		w := httptest.NewRecorder()
		token := "full_cycle_token"
		svc.SetRefreshToken(w, token)

		// Получение из ответа
		result := w.Result()
		cookies := result.Cookies()
		require.Len(t, cookies, 1)

		// Создание нового запроса с cookie
		r := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		r.AddCookie(cookies[0])

		// Чтение
		retrievedToken, err := svc.GetRefreshToken(r)
		require.NoError(t, err)
		assert.Equal(t, token, retrievedToken)

		// Удаление
		w2 := httptest.NewRecorder()
		svc.DeleteRefreshToken(w2)

		deleteCookies := w2.Result().Cookies()
		require.Len(t, deleteCookies, 1)
		assert.Equal(t, -1, deleteCookies[0].MaxAge)
	})
}
