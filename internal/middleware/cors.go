// Package middleware предоставляет HTTP и gRPC middleware для приложения.
package middleware

import (
	"net/http"
	"strings"

	"github.com/rs/cors"
)

// CORSConfig конфигурация CORS
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
	Debug            bool
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
		AllowOriginFunc: func(origin string) bool {
			for _, allowed := range config.AllowedOrigins {
				if allowed == "*" {
					if config.AllowCredentials {
						return false
					}
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

	return c.Handler
}
