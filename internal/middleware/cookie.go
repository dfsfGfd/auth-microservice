// Package middleware предоставляет HTTP middleware для микросервиса.
package middleware

import (
	"net/http"
	"time"
)

// CookieConfig конфигурация cookie
type CookieConfig struct {
	Name     string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}

// SetCookie устанавливает cookie в ответ
func SetCookie(w http.ResponseWriter, config CookieConfig) {
	cookie := &http.Cookie{
		Name:     config.Name,
		Value:    config.Value,
		MaxAge:   config.MaxAge,
		Path:     config.Path,
		Domain:   config.Domain,
		Secure:   config.Secure,
		HttpOnly: config.HTTPOnly,
		SameSite: config.SameSite,
	}
	http.SetCookie(w, cookie)
}

// SetRefreshTokenCookie устанавливает refresh токен в HttpOnly cookie
func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string, secure bool, domain string) {
	SetCookie(w, CookieConfig{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   int(14 * 24 * time.Hour / time.Second), // 14 дней
		Path:     "/api/auth",
		Domain:   domain,
		Secure:   secure,
		HTTPOnly: true, // Защита от XSS
		SameSite: http.SameSiteStrictMode,
	})
}

// DeleteRefreshTokenCookie удаляет refresh токен cookie
func DeleteRefreshTokenCookie(w http.ResponseWriter, domain string) {
	SetCookie(w, CookieConfig{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/api/auth",
		Domain:   domain,
		Secure:   false,
		HTTPOnly: true,
	})
}
