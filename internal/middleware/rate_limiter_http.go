package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// HTTPRateLimitMiddleware создаёт HTTP middleware для rate limiting
// Возвращает функцию-обёртку для http.Handler
func HTTPRateLimitMiddleware(rl *RateLimiter, getEndpoint func(r *http.Request) string, getKey func(r *http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			endpoint := getEndpoint(r)
			key := getKey(r)

			allowed, remaining, resetTime, err := rl.Allow(r.Context(), endpoint, key)
			if err != nil {
				// Для критичных endpoint'ов (login, register) блокируем запрос
				// Для остальных — пропускаем (fail-open)
				if endpoint == "login" || endpoint == "register" {
					http.Error(w, `{"error":"rate limiter unavailable"}`, http.StatusServiceUnavailable)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			// Добавляем заголовки rate limiting
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(getLimit(rl, endpoint)))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

			if !allowed {
				retryAfter := int(time.Until(resetTime).Seconds())
				if retryAfter < 0 {
					retryAfter = 0
				}
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getLimit возвращает лимит для endpoint
func getLimit(rl *RateLimiter, endpoint string) int {
	if config, ok := rl.configs[endpoint]; ok {
		return config.Limit
	}
	return 0
}

// UnaryServerInterceptor создаёт gRPC interceptor для rate limiting
func UnaryServerInterceptor(rl *RateLimiter, getEndpoint func(ctx context.Context, fullMethod string) string, getKey func(ctx context.Context, fullMethod string) string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		endpoint := getEndpoint(ctx, info.FullMethod)
		key := getKey(ctx, info.FullMethod)

		allowed, remaining, resetTime, err := rl.Allow(ctx, endpoint, key)
		if err != nil {
			// Для критичных endpoint'ов (login, register) блокируем запрос
			// Для остальных — пропускаем (fail-open)
			if endpoint == "login" || endpoint == "register" {
				return nil, status.Error(codes.Unavailable, "rate limiter unavailable")
			}
			return handler(ctx, req)
		}

		// Добавляем метаданные rate limiting
		// Игнорируем ошибку, так как это не критично для работы
		_ = grpc.SendHeader(ctx, metadata.Pairs(
			"X-RateLimit-Limit", strconv.Itoa(getLimit(rl, endpoint)),
			"X-RateLimit-Remaining", strconv.Itoa(remaining),
			"X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10),
		))

		if !allowed {
			retryAfter := int(time.Until(resetTime).Seconds())
			if retryAfter < 0 {
				retryAfter = 0
			}
			// Игнорируем ошибку отправки заголовка — не критично
			_ = grpc.SendHeader(ctx, metadata.Pairs("Retry-After", strconv.Itoa(retryAfter)))
			return nil, status.Error(codes.ResourceExhausted, "rate limit exceeded")
		}

		return handler(ctx, req)
	}
}

// MethodKeyFunc возвращает функцию для получения ключа из gRPC метода + IP
func MethodKeyFunc() func(ctx context.Context, fullMethod string) string {
	return func(ctx context.Context, fullMethod string) string {
		// Пытаемся получить IP из метаданных
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			xff := md.Get("X-Forwarded-For")
			if len(xff) > 0 {
				return "method:" + fullMethod + ":ip:" + xff[0]
			}
		}

		return "method:" + fullMethod
	}
}

// MethodEndpointFunc возвращает функцию для получения endpoint из gRPC метода
func MethodEndpointFunc() func(ctx context.Context, fullMethod string) string {
	return func(ctx context.Context, fullMethod string) string {
		// Извлекаем имя метода из полного пути (например, /auth.v1.AuthService/Login -> login)
		switch fullMethod {
		case "/auth.v1.AuthService/Register":
			return "register"
		case "/auth.v1.AuthService/Login":
			return "login"
		case "/auth.v1.AuthService/Refresh":
			return "refresh"
		case "/auth.v1.AuthService/Logout":
			return "logout"
		default:
			return ""
		}
	}
}

// HTTPPathEndpointFunc возвращает функцию для получения endpoint из HTTP пути
func HTTPPathEndpointFunc() func(r *http.Request) string {
	return func(r *http.Request) string {
		switch r.URL.Path {
		case "/api/auth/register":
			return "register"
		case "/api/auth/login":
			return "login"
		case "/api/auth/refresh":
			return "refresh"
		case "/api/auth/logout":
			return "logout"
		default:
			return ""
		}
	}
}

// IPKeyFunc возвращает функцию для получения ключа из IP адреса
func IPKeyFunc() func(r *http.Request) string {
	return func(r *http.Request) string {
		// Проверяем X-Forwarded-For для proxy
		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			return "ip:" + xff
		}

		// Проверяем X-Real-IP
		xri := r.Header.Get("X-Real-IP")
		if xri != "" {
			return "ip:" + xri
		}

		// Используем RemoteAddr
		return "ip:" + r.RemoteAddr
	}
}
