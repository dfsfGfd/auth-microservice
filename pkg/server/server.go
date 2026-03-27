// Package server предоставляет запуск gRPC и HTTP серверов.
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-microservice/internal/config"
	"auth-microservice/internal/middleware"
	"auth-microservice/pkg/logger"
	"auth-microservice/pkg/proto/auth/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Server управляет gRPC и HTTP серверами.
type Server struct {
	cfg       *config.Config
	log       *logger.Logger
	rateLimit *middleware.RateLimiter
	handler   authv1.AuthServiceServer
}

// NewServer создаёт новый сервер.
func NewServer(cfg *config.Config, log *logger.Logger, rateLimit *middleware.RateLimiter, handler authv1.AuthServiceServer) *Server {
	return &Server{
		cfg:       cfg,
		log:       log,
		rateLimit: rateLimit,
		handler:   handler,
	}
}

// Run запускает сервера и ожидает сигнал завершения.
func (s *Server) Run() error {
	ctx := context.Background()

	// gRPC сервер
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Server.GRPCPort))
	if err != nil {
		return fmt.Errorf("create gRPC listener: %w", err)
	}

	rateLimitInterceptor := middleware.UnaryServerInterceptor(
		s.rateLimit,
		middleware.MethodEndpointFunc(),
		middleware.MethodKeyFunc(),
	)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(rateLimitInterceptor, s.unaryLogger()),
	)

	authv1.RegisterAuthServiceServer(grpcServer, s.handler)
	reflection.Register(grpcServer)

	// Канал готовности
	grpcReady := make(chan struct{})

	// Запуск gRPC
	go func() {
		close(grpcReady)
		s.log.Info("grpc_server_start", "port", s.cfg.Server.GRPCPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			s.log.Error("grpc_server_error", "err", err)
		}
	}()

	// Ждём готовности gRPC
	<-grpcReady

	// HTTP Gateway
	gw, conn, err := s.createGateway(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	gwWithRateLimit := middleware.HTTPRateLimitMiddleware(
		s.rateLimit,
		middleware.HTTPPathEndpointFunc(),
		middleware.IPKeyFunc(),
	)(gw)

	rootMux := http.NewServeMux()
	rootMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	rootMux.Handle("/", gwWithRateLimit)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Server.HTTPPort),
		Handler:      rootMux,
		ReadTimeout:  s.cfg.Server.ReadTimeoutDuration(),
		WriteTimeout: s.cfg.Server.WriteTimeoutDuration(),
		IdleTimeout:  s.cfg.Server.IdleTimeoutDuration(),
	}

	// Запуск HTTP
	go func() {
		s.log.Info("rest_server_start", "port", s.cfg.Server.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("rest_server_error", "err", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.Shutdown.TimeoutDuration())
	defer cancel()

	grpcServer.GracefulStop()
	s.log.Info("grpc_server_stop")

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	s.log.Info("rest_server_stop")

	return nil
}

func (s *Server) createGateway(ctx context.Context) (*runtime.ServeMux, *grpc.ClientConn, error) {
	gwMux := runtime.NewServeMux()

	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", s.cfg.Server.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("dial grpc: %w", err)
	}

	if err := authv1.RegisterAuthServiceHandler(ctx, gwMux, conn); err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("register gateway: %w", err)
	}

	return gwMux, conn, nil
}

func (s *Server) unaryLogger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			s.log.Error("grpc_request", "method", info.FullMethod, "dur_ms", duration.Milliseconds(), "err", err)
		} else {
			s.log.Info("grpc_request", "method", info.FullMethod, "dur_ms", duration.Milliseconds())
		}

		return resp, err
	}
}
