package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Стандартные ошибки пакета
var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token expired")
	ErrInvalidClaims = errors.New("invalid token claims")
)

// tokenType определяет тип токена
type tokenType string

const (
	accessToken  tokenType = "access"
	refreshToken tokenType = "refresh"
)

// claims представляет claims JWT токена
type claims struct {
	Type tokenType `json:"type"`
	jwt.RegisteredClaims
}

// Config конфигурация JWT сервиса
type Config struct {
	// SecretKey секретный ключ для подписи токенов
	SecretKey string
	// AccessTokenTTL время жизни access токена
	AccessTokenTTL time.Duration
	// RefreshTokenTTL время жизни refresh токена
	RefreshTokenTTL time.Duration
	// Issuer название сервиса (iss claim)
	Issuer string
}

// Validate валидирует конфигурацию
func (c *Config) Validate() error {
	if c.SecretKey == "" {
		return errors.New("secret key is required")
	}
	if c.AccessTokenTTL <= 0 {
		return errors.New("access token TTL must be positive")
	}
	if c.RefreshTokenTTL <= 0 {
		return errors.New("refresh token TTL must be positive")
	}
	if c.Issuer == "" {
		return errors.New("issuer is required")
	}
	return nil
}

// Service сервис для работы с JWT токенами
type Service struct {
	config Config
}

// TokenPair пара токенов (access + refresh)
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
}

// NewService создаёт новый JWT сервис
func NewService(config Config) (*Service, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		config: config,
	}, nil
}

// GenerateTokens генерирует пару access и refresh токенов
func (s *Service) GenerateTokens(accountID, email string) (*TokenPair, error) {
	accessToken, err := s.generateToken(accountID, email, accessToken, s.config.AccessTokenTTL)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(accountID, email, refreshToken, s.config.RefreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.AccessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// generateToken создаёт JWT токен
func (s *Service) generateToken(accountID, email string, tokenType tokenType, ttl time.Duration) (string, error) {
	now := time.Now()

	claims := claims{
		Type: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.Issuer,
			Subject:   accountID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

// ValidateToken валидирует JWT токен и возвращает claims
func (s *Service) ValidateToken(tokenString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(s.config.SecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	c, ok := token.Claims.(*claims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return c, nil
}

// ValidateRefreshToken валидирует refresh токен
func (s *Service) ValidateRefreshToken(tokenString string) (*claims, error) {
	c, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if c.Type != refreshToken {
		return nil, ErrInvalidClaims
	}

	return c, nil
}

// RefreshTTLDuration возвращает время жизни refresh токена
func (s *Service) RefreshTTLDuration() time.Duration {
	return s.config.RefreshTokenTTL
}
