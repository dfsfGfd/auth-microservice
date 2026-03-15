// Package cookies предоставляет утилиты для работы с HTTP cookie.
//
// Пример использования:
//
//	// Создание сервиса
//	cookieService := cookies.NewService(cookies.Config{
//	    Secure:   true,
//	    HTTPOnly: true,
//	    SameSite: http.SameSiteStrictMode,
//	    Domain:   "example.com",
//	})
//
//	// Установка cookie
//	cookieService.SetRefreshToken(w, refreshToken)
//
//	// Получение cookie
//	refreshToken := cookieService.GetRefreshToken(r)
package cookies

import (
	"net/http"
	"time"
)

// Config конфигурация cookie сервиса
type Config struct {
	// Secure устанавливает флаг Secure (передача только по HTTPS)
	Secure bool
	// HTTPOnly устанавливает флаг HttpOnly (недоступны для JavaScript)
	HTTPOnly bool
	// SameSite политика SameSite (Strict, Lax, None)
	SameSite http.SameSite
	// Domain домен cookie (опционально)
	Domain string
	// Path путь cookie
	Path string
	// MaxAge время жизни cookie в секундах
	MaxAge int
}

// Validate валидирует и устанавливает значения по умолчанию
func (c *Config) Validate() {
	if c.Path == "" {
		c.Path = "/"
	}
	if c.MaxAge == 0 {
		// 14 дней по умолчанию для refresh токена
		c.MaxAge = int((14 * 24 * time.Hour).Seconds())
	}
}

// Service сервис для работы с cookie
type Service struct {
	config Config
}

// NewService создаёт новый cookie сервис
func NewService(config Config) *Service {
	config.Validate()
	return &Service{config: config}
}

// SetRefreshToken устанавливает refresh токен в cookie
func (s *Service) SetRefreshToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token,
		MaxAge:   s.config.MaxAge,
		Path:     s.config.Path,
		Domain:   s.config.Domain,
		Secure:   s.config.Secure,
		HttpOnly: s.config.HTTPOnly,
		SameSite: s.config.SameSite,
	})
}

// GetRefreshToken получает refresh токен из cookie
func (s *Service) GetRefreshToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// DeleteRefreshToken удаляет refresh токен из cookie
func (s *Service) DeleteRefreshToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     s.config.Path,
		Domain:   s.config.Domain,
		Secure:   s.config.Secure,
		HttpOnly: s.config.HTTPOnly,
		SameSite: s.config.SameSite,
	})
}

// SetAccessToken устанавливает access токен в cookie (опционально)
// Обычно access токен передаётся в Authorization header
func (s *Service) SetAccessToken(w http.ResponseWriter, token string, expiresIn int) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		MaxAge:   expiresIn,
		Path:     s.config.Path,
		Domain:   s.config.Domain,
		Secure:   s.config.Secure,
		HttpOnly: s.config.HTTPOnly,
		SameSite: s.config.SameSite,
	})
}

// GetAccessToken получает access токен из cookie
func (s *Service) GetAccessToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("access_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// DeleteAccessToken удаляет access токен из cookie
func (s *Service) DeleteAccessToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		Path:     s.config.Path,
		Domain:   s.config.Domain,
		Secure:   s.config.Secure,
		HttpOnly: s.config.HTTPOnly,
		SameSite: s.config.SameSite,
	})
}

// DeleteAll удаляет все cookie
func (s *Service) DeleteAll(w http.ResponseWriter) {
	s.DeleteRefreshToken(w)
	s.DeleteAccessToken(w)
}

// GetConfig возвращает конфигурацию сервиса
func (s *Service) GetConfig() Config {
	return s.config
}
