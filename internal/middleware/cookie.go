// Package middleware предоставляет HTTP middleware для микросервиса.
package middleware

import (
	"bufio"
	"encoding/json"
	"net"
	"net/http"

	"auth-microservice/pkg/cookies"
	"auth-microservice/pkg/proto/auth/v1"
)

// CookieMiddleware middleware для установки refresh токена в HttpOnly cookie.
// Работает с grpc-gateway ответами, перехватывая и модифицируя ответ до отправки клиенту.
func CookieMiddleware(cookieService *cookies.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Оборачиваем ResponseWriter для перехвата ответа
			wrapper := &responseWrapper{
				ResponseWriter: w,
				body:           &buffer{},
			}

			// Вызываем следующий handler
			next.ServeHTTP(wrapper, r)

			// Проверяем путь и статус — только для login и refresh
			if wrapper.status == http.StatusOK {
				switch r.URL.Path {
				case "/api/auth/login", "/api/auth/refresh":
					// Парсим ответ для извлечения refresh токена
					var resp authv1.LoginResponse
					if err := json.Unmarshal(wrapper.body.Bytes(), &resp); err == nil {
						if resp.Data != nil && resp.Data.RefreshToken != "" {
							// Устанавливаем refresh токен в HttpOnly cookie
							cookieService.SetRefreshToken(w, resp.Data.RefreshToken)
						}
					}
				case "/api/auth/logout":
					// Удаляем cookie при logout
					cookieService.DeleteRefreshToken(w)
				}
			}
		})
	}
}

// buffer простой буфер для записи тела
type buffer struct {
	data []byte
}

func (b *buffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *buffer) Bytes() []byte {
	return b.data
}

// responseWrapper оборачивает http.ResponseWriter для перехвата статуса и тела
type responseWrapper struct {
	http.ResponseWriter
	body   *buffer
	status int
	wrote  bool
}

func (rw *responseWrapper) WriteHeader(code int) {
	if rw.wrote {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wrote = true
}

func (rw *responseWrapper) Write(b []byte) (int, error) {
	if !rw.wrote {
		rw.status = http.StatusOK
		rw.wrote = true
	}
	// Сохраняем копию тела
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// Реализуем дополнительные интерфейсы для совместимости
func (rw *responseWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hj.Hijack()
}

func (rw *responseWrapper) Flush() {
	fl, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		fl.Flush()
	}
}

// Unwrap возвращает оригинальный ResponseWriter
func (rw *responseWrapper) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}
