// Package middleware предоставляет HTTP и gRPC middleware для приложения.
package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
)

// CORSConfig конфигурация CORS
type CORSConfig struct {
	// AllowedOrigins разрешённые origin
	AllowedOrigins []string
	// AllowedMethods разрешённые методы
	AllowedMethods []string
	// AllowedHeaders разрешённые заголовки
	AllowedHeaders []string
	// AllowCredentials разрешает отправку credentials
	AllowCredentials bool
	// MaxAge для preflight запросов
	MaxAge int
	// Debug включает debug логирование
	Debug bool
}

// Validate валидирует и устанавливает значения по умолчанию
func (c *CORSConfig) Validate() {
	if len(c.AllowedOrigins) == 0 {
		c.AllowedOrigins = []string{"*"}
	}
	if len(c.AllowedMethods) == 0 {
		c.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if len(c.AllowedHeaders) == 0 {
		c.AllowedHeaders = []string{"Authorization", "Content-Type", "X-Request-ID"}
	}
	if c.MaxAge <= 0 {
		c.MaxAge = 86400
	}
}

// NewCORS создаёт CORS middleware
func NewCORS(config CORSConfig) func(http.Handler) http.Handler {
	config.Validate()

	c := cors.New(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   config.AllowedMethods,
		AllowedHeaders:   config.AllowedHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           config.MaxAge,
		Debug:            config.Debug,
		// Обработка wildcard origin
		AllowOriginFunc: func(origin string) bool {
			// Если есть wildcard, разрешаем все
			for _, allowed := range config.AllowedOrigins {
				if allowed == "*" {
					return true
				}
				// Поддержка wildcard поддоменов
				if strings.HasPrefix(allowed, "*.") {
					suffix := allowed[1:] // .example.com
					if strings.HasSuffix(origin, suffix) {
						return true
					}
				}
				// Точное совпадение
				if allowed == origin {
					return true
				}
			}
			return false
		},
	})

	return c.Handler
}

// WithDefaultCORS создаёт CORS middleware с настройками по умолчанию
func WithDefaultCORS() func(http.Handler) http.Handler {
	return NewCORS(CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: false,
		MaxAge:           86400,
		Debug:            false,
	})
}

// WithSecureCORS создаёт безопасный CORS middleware для production
func WithSecureCORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return NewCORS(CORSConfig{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"},
		AllowCredentials: true,
		MaxAge:           86400,
		Debug:            false,
	})
}

// Middleware создаёт CORS handler с таймаутом
func Middleware(handler http.Handler, config CORSConfig, readTimeout, writeTimeout time.Duration) http.Handler {
	config.Validate()

	c := cors.New(cors.Options{
		AllowedOrigins:   config.AllowedOrigins,
		AllowedMethods:   config.AllowedMethods,
		AllowedHeaders:   config.AllowedHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           config.MaxAge,
		Debug:            config.Debug,
		AllowOriginFunc: func(origin string) bool {
			for _, allowed := range config.AllowedOrigins {
				if allowed == "*" {
					return true
				}
				if strings.HasPrefix(allowed, "*.") {
					suffix := allowed[1:]
					if strings.HasSuffix(origin, suffix) {
						return true
					}
				}
				if allowed == origin {
					return true
				}
			}
			return false
		},
	})

	// Оборачиваем в CORS middleware
	corsHandler := c.Handler(handler)

	// Добавляем timeout handler
	return http.TimeoutHandler(corsHandler, readTimeout+writeTimeout, "request timeout")
}
