package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-microservice/internal/di"
	"auth-microservice/internal/middleware"
	"auth-microservice/pkg/logger"
	"auth-microservice/pkg/proto/auth/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Инициализация приложения через DI
	app, err := di.InitializeApplication()
	if err != nil {
		return fmt.Errorf("initialize application: %w", err)
	}

	ctx := context.Background()

	log := app.Logger
	defer func() {
		if err := app.CleanUp(); err != nil {
			log.Error("cleanup_error", "err", err)
		}
	}()

	// Создаём gRPC сервер с rate limiting и логированием
	rateLimitInterceptor := middleware.UnaryServerInterceptor(
		app.RateLimiter,
		middleware.MethodEndpointFunc(),
		middleware.MethodKeyFunc(),
	)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(rateLimitInterceptor, unaryLogger(log)),
	)

	// Регистрируем Auth handler
	authv1.RegisterAuthServiceServer(grpcServer, app.AuthHandler)

	// Включаем reflection для gRPC
	reflection.Register(grpcServer)

	// Создаём listener для gRPC
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.Config.Server.GRPCPort))
	if err != nil {
		return fmt.Errorf("create gRPC listener: %w", err)
	}

	// Создаём HTTP сервер для REST (запускается позже)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Server.HTTPPort),
		ReadTimeout:  time.Duration(app.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(app.Config.Server.IdleTimeout) * time.Second,
	}

	// Канал для сигналов завершения
	errCh := make(chan error, 2)

	// Канал готовности gRPC сервера
	grpcReady := make(chan struct{})

	// Запускаем gRPC сервер
	go func() {
		close(grpcReady)  // Сигнал готовности
		log.Info("grpc_server_start", "port", app.Config.Server.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			errCh <- fmt.Errorf("serve gRPC: %w", err)
		}
	}()

	// Ждём готовности gRPC сервера
	<-grpcReady

	// Настраиваем grpc-gateway (после запуска gRPC)
	gw, grpcConn, err := createGateway(ctx, app.Config.Server.GRPCPort)
	if err != nil {
		return fmt.Errorf("create gateway: %w", err)
	}
	defer grpcConn.Close()

	// Оборачиваем gateway в rate limiting middleware
	gwWithRateLimit := middleware.HTTPRateLimitMiddleware(
		app.RateLimiter,
		middleware.HTTPPathEndpointFunc(),
		middleware.IPKeyFunc(),
	)(gw)

	// Создаём корневой mux и добавляем health check + gateway
	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	rootMux.Handle("/", gwWithRateLimit)

	// Устанавливаем handler HTTP сервера
	httpServer.Handler = rootMux

	// Запускаем HTTP сервер (REST)
	go func() {
		log.Info("rest_server_start", "port", app.Config.Server.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("serve HTTP: %w", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Info("shutdown_signal", "signal", sig)
	case err := <-errCh:
		log.Error("server_error", "err", err)
		return err
	}

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		time.Duration(app.Config.Shutdown.Timeout)*time.Second,
	)
	defer shutdownCancel()

	// Останавливаем gRPC сервер
	grpcServer.GracefulStop()
	log.Info("grpc_server_stop")

	// Останавливаем HTTP сервер
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	log.Info("rest_server_stop")

	log.Info("application_shutdown")
	return nil
}

// createGateway создаёт grpc-gateway mux и возвращает соединение
func createGateway(ctx context.Context, grpcPort int) (*runtime.ServeMux, *grpc.ClientConn, error) {
	gwMux := runtime.NewServeMux()

	// Создаём gRPC соединение для gateway
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("dial grpc: %w", err)
	}

	// Регистрируем gateway
	if err := authv1.RegisterAuthServiceHandler(ctx, gwMux, conn); err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("register gateway: %w", err)
	}

	return gwMux, conn, nil
}

// unaryLogger создаёт interceptor для логирования gRPC запросов
func unaryLogger(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		// Логгируем с оптимизированным форматом
		if err != nil {
			log.Error("grpc_request",
				"method", info.FullMethod,
				"dur_ms", duration.Milliseconds(),
				"err", err,
			)
		} else {
			log.Info("grpc_request",
				"method", info.FullMethod,
				"dur_ms", duration.Milliseconds(),
			)
		}

		return resp, err
	}
}
