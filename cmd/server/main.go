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

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"auth-microservice/internal/di"
	"auth-microservice/pkg/logger"
	"auth-microservice/pkg/proto/auth/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	defer app.CleanUp(ctx)

	log := app.Logger

	// Создаём gRPC сервер
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryLogger(log)),
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

	// Создаём mux для REST (grpc-gateway)
	gwMux := runtime.NewServeMux()

	// Создаём HTTP сервер для REST (запускается позже)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Server.HTTPPort),
		Handler:      gwMux,
		ReadTimeout:  time.Duration(app.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(app.Config.Server.IdleTimeout) * time.Second,
	}

	// Канал для сигналов завершения
	errCh := make(chan error, 2)

	// Запускаем gRPC сервер
	go func() {
		log.Info("starting gRPC server", "port", app.Config.Server.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			errCh <- fmt.Errorf("serve gRPC: %w", err)
		}
	}()

	// Даем gRPC серверу время запуститься
	time.Sleep(100 * time.Millisecond)

	// Настраиваем grpc-gateway (после запуска gRPC)
	gw, err := createGateway(ctx, app.Config.Server.GRPCPort)
	if err != nil {
		return fmt.Errorf("create gateway: %w", err)
	}

	// Создаём корневой mux и добавляем health check + gateway
	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	rootMux.Handle("/", gw)

	// Обновляем handler HTTP сервера
	httpServer.Handler = rootMux

	// Запускаем HTTP сервер (REST)
	go func() {
		log.Info("starting REST server (grpc-gateway)", "port", app.Config.Server.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("serve HTTP: %w", err)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Info("shutdown signal received", "signal", sig)
	case err := <-errCh:
		log.Error("server error", "error", err)
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
	log.Info("gRPC server stopped")

	// Останавливаем HTTP сервер
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	log.Info("REST server stopped")

	log.Info("application shutdown complete")
	return nil
}

// createGateway создаёт grpc-gateway mux
func createGateway(ctx context.Context, grpcPort int) (*runtime.ServeMux, error) {
	gwMux := runtime.NewServeMux()

	// Создаём gRPC connection для gateway
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial grpc: %w", err)
	}
	defer conn.Close()

	// Регистрируем gateway
	if err := authv1.RegisterAuthServiceHandler(ctx, gwMux, conn); err != nil {
		return nil, fmt.Errorf("register gateway: %w", err)
	}

	return gwMux, nil
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
		log.Info(
			"gRPC request",
			"method", info.FullMethod,
			"duration", duration.String(),
			"error", err,
		)

		return resp, err
	}
}
