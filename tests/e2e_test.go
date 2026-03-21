//go:build integration

package tests

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"auth-microservice/pkg/proto/auth/v1"

	"github.com/stretchr/testify/suite"
)

// E2ETestSuite набор e2e тестов для auth-сервиса
type E2ETestSuite struct {
	suite.Suite
	ctx        context.Context
	grpcClient authv1.AuthServiceClient
	grpcConn   *grpc.ClientConn
}

// SetupSuite выполняется один раз перед всеми тестами
func (s *E2ETestSuite) SetupSuite() {
	s.ctx = context.Background()

	// Создаём gRPC соединение с тестовым сервисом
	var err error
	s.grpcConn, err = grpc.NewClient(
		getGRPCAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(10*time.Second),
	)
	s.Require().NoError(err, "Failed to create gRPC connection")

	s.grpcClient = authv1.NewAuthServiceClient(s.grpcConn)
}

// TearDownSuite выполняется один раз после всех тестов
func (s *E2ETestSuite) TearDownSuite() {
	if s.grpcConn != nil {
		_ = s.grpcConn.Close()
	}
}

// TestRegister тестирует регистрацию нового пользователя
func (s *E2ETestSuite) TestRegister() {
	tests := []struct {
		name        string
		email       string
		password    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid registration",
			email:    "test.e2e@example.com",
			password: "SecurePass123!",
			wantErr:  false,
		},
		{
			name:        "invalid email format",
			email:       "invalid-email",
			password:    "SecurePass123!",
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name:        "password too short",
			email:       "test2.e2e@example.com",
			password:    "123",
			wantErr:     true,
			errContains: "password too short",
		},
		{
			name:        "empty email",
			email:       "",
			password:    "SecurePass123!",
			wantErr:     true,
			errContains: "email",
		},
		{
			name:        "empty password",
			email:       "test3.e2e@example.com",
			password:    "",
			wantErr:     true,
			errContains: "invalid password",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.grpcClient.Register(s.ctx, &authv1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})

			if tt.wantErr {
				s.Require().Error(err)
				if tt.errContains != "" {
					s.Contains(err.Error(), tt.errContains)
				}
				return
			}

			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().NotNil(resp.Data)
			s.NotEmpty(resp.Data.GetAccountId())
			s.Equal(tt.email, resp.Data.GetEmail())
			s.NotEmpty(resp.Data.GetCreatedAt())
		})
	}
}

// TestRegister_ValidPasswords тестирует регистрацию с различными валидными паролями
// После упрощения валидации (NIST 800-63B): только длина ≥8, без требований к регистру/цифрам
func (s *E2ETestSuite) TestRegister_ValidPasswords() {
	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"no uppercase", "test.upper@example.com", "password123"},
		{"no lowercase", "test.lower@example.com", "PASSWORD123"},
		{"no digit", "test.digit@example.com", "Password"},
		{"all lowercase", "test.alllower@example.com", "abcdefgh"},
		{"simple 8 chars", "test.simple@example.com", "qwerty12"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.grpcClient.Register(s.ctx, &authv1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})

			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().NotNil(resp.Data)
			s.NotEmpty(resp.Data.GetAccountId())
			s.Equal(tt.email, resp.Data.GetEmail())
		})
	}
}

// TestLogin тестирует вход пользователя
func (s *E2ETestSuite) TestLogin() {
	// Сначала регистрируем уникального пользователя
	email := "login.e2e@example.com"
	password := "SecurePass123!"

	_, err := s.grpcClient.Register(s.ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	s.Require().NoError(err)

	tests := []struct {
		name        string
		email       string
		password    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid login",
			email:    email,
			password: password,
			wantErr:  false,
		},
		{
			name:        "wrong password",
			email:       email,
			password:    "WrongPass123!",
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name:        "non-existent user",
			email:       "nonexistent.e2e@example.com",
			password:    password,
			wantErr:     true,
			errContains: "invalid",
		},
		{
			name:        "empty email",
			email:       "",
			password:    password,
			wantErr:     true,
			errContains: "invalid",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.grpcClient.Login(s.ctx, &authv1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			})

			if tt.wantErr {
				s.Require().Error(err)
				if tt.errContains != "" {
					s.Contains(err.Error(), tt.errContains)
				}
				return
			}

			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().NotNil(resp.Data)
			s.NotEmpty(resp.Data.GetAccessToken())
			s.NotEmpty(resp.Data.GetRefreshToken())
			s.Equal("Bearer", resp.Data.GetTokenType())
			s.Greater(resp.Data.GetExpiresIn(), int32(0))
		})
	}
}

// TestRefresh тестирует обновление токенов
func (s *E2ETestSuite) TestRefresh() {
	// Регистрируем и логинимся
	email := "refresh.e2e@example.com"
	password := "SecurePass123!"

	_, err := s.grpcClient.Register(s.ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	s.Require().NoError(err)

	loginResp, err := s.grpcClient.Login(s.ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	s.Require().NoError(err)

	tests := []struct {
		name         string
		refreshToken string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "valid refresh",
			refreshToken: loginResp.Data.GetRefreshToken(),
			wantErr:      false,
		},
		{
			name:         "invalid token",
			refreshToken: "invalid-token-xyz-123",
			wantErr:      true,
			errContains:  "invalid",
		},
		{
			name:         "empty token",
			refreshToken: "",
			wantErr:      true,
			errContains:  "token",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			resp, err := s.grpcClient.Refresh(s.ctx, &authv1.RefreshRequest{
				RefreshToken: tt.refreshToken,
			})

			if tt.wantErr {
				s.Require().Error(err)
				if tt.errContains != "" {
					s.Contains(err.Error(), tt.errContains)
				}
				return
			}

			s.Require().NoError(err)
			s.Require().NotNil(resp)
			s.Require().NotNil(resp.Data)
			s.NotEmpty(resp.Data.GetAccessToken())
			s.NotEmpty(resp.Data.GetRefreshToken())
			s.Equal("Bearer", resp.Data.GetTokenType())
		})
	}
}

// TestLogout тестирует выход пользователя
func (s *E2ETestSuite) TestLogout() {
	// Регистрируем и логинимся
	email := "logout.e2e@example.com"
	password := "SecurePass123!"

	_, err := s.grpcClient.Register(s.ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	s.Require().NoError(err)

	loginResp, err := s.grpcClient.Login(s.ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	s.Require().NoError(err)

	// Тест 1: валидный logout
	s.Run("valid logout", func() {
		resp, err := s.grpcClient.Logout(s.ctx, &authv1.LogoutRequest{
			RefreshToken: loginResp.Data.GetRefreshToken(),
		})

		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.True(resp.Data.GetSuccess())
	})

	// Тест 2: попытка использовать отозванный токен - сервис может возвращать разное поведение
	s.Run("logout with revoked token", func() {
		// Некоторые сервисы позволяют повторный logout без ошибки (idempotent)
		_, _ = s.grpcClient.Logout(s.ctx, &authv1.LogoutRequest{
			RefreshToken: loginResp.Data.GetRefreshToken(),
		})

		// Проверяем что токен больше не работает для refresh
		refreshResp, refreshErr := s.grpcClient.Refresh(s.ctx, &authv1.RefreshRequest{
			RefreshToken: loginResp.Data.GetRefreshToken(),
		})
		s.Require().Error(refreshErr)
		s.Nil(refreshResp)
	})

	// Тест 3: logout с несуществующим токеном
	s.Run("logout with invalid token", func() {
		_, err := s.grpcClient.Logout(s.ctx, &authv1.LogoutRequest{
			RefreshToken: "nonexistent-token",
		})

		s.Require().Error(err)
		s.Contains(err.Error(), "invalid")
	})
}

// TestDuplicateRegistration тестирует защиту от дубликатов
func (s *E2ETestSuite) TestDuplicateRegistration() {
	email := "duplicate.e2e@example.com"
	password := "SecurePass123!"

	// Первая регистрация успешна
	resp1, err := s.grpcClient.Register(s.ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	s.Require().NoError(err)
	s.Require().NotNil(resp1)

	// Вторая регистрация с тем же email должна вернуть ошибку
	resp2, err := s.grpcClient.Register(s.ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password + "2", // другой пароль
	})
	s.Require().Error(err)
	s.Nil(resp2)
	s.Contains(err.Error(), "already exists")
}

// TestFullAuthFlow тестирует полный цикл аутентификации
func (s *E2ETestSuite) TestFullAuthFlow() {
	email := "fullflow.e2e@example.com"
	password := "SecurePass123!"

	// 1. Регистрация
	s.Run("step 1: register", func() {
		resp, err := s.grpcClient.Register(s.ctx, &authv1.RegisterRequest{
			Email:    email,
			Password: password,
		})
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.NotEmpty(resp.Data.GetAccountId())
	})

	// 2. Login
	var accessToken, refreshToken string
	s.Run("step 2: login", func() {
		resp, err := s.grpcClient.Login(s.ctx, &authv1.LoginRequest{
			Email:    email,
			Password: password,
		})
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		accessToken = resp.Data.GetAccessToken()
		refreshToken = resp.Data.GetRefreshToken()
		s.NotEmpty(accessToken)
		s.NotEmpty(refreshToken)
	})

	// 3. Refresh
	s.Run("step 3: refresh", func() {
		resp, err := s.grpcClient.Refresh(s.ctx, &authv1.RefreshRequest{
			RefreshToken: refreshToken,
		})
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.NotEmpty(resp.Data.GetAccessToken())
		s.NotEmpty(resp.Data.GetRefreshToken())
	})

	// 4. Logout
	s.Run("step 4: logout", func() {
		resp, err := s.grpcClient.Logout(s.ctx, &authv1.LogoutRequest{
			RefreshToken: refreshToken,
		})
		s.Require().NoError(err)
		s.Require().NotNil(resp)
		s.True(resp.Data.GetSuccess())
	})

	// 5. Попытка refresh после logout должна вернуть ошибку
	s.Run("step 5: refresh after logout fails", func() {
		resp, err := s.grpcClient.Refresh(s.ctx, &authv1.RefreshRequest{
			RefreshToken: refreshToken,
		})
		s.Require().Error(err)
		s.Nil(resp)
	})
}

// TestHealthCheck тестирует health check endpoint через HTTP
func (s *E2ETestSuite) TestHealthCheck() {
	s.T().Skip("Health check tested in setup_test.go")
}

// TestIntegration запуск e2e тестов
func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
