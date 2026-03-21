package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// BCryptCost стоимость хеширования bcrypt.
// Рекомендация OWASP 2024: минимум 12 для production.
// Увеличивайте на +1 каждые 1-2 года по мере роста мощностей CPU.
const BCryptCost = 12

// Hasher предоставляет интерфейс для хеширования и верификации паролей
type Hasher interface {
	Hash(password string, cost int) (string, error)
	Compare(hash, password string) error
}

// Service реализует Hasher используя bcrypt
type Service struct{}

// NewService создаёт новый сервис хеширования
func NewService() *Service {
	return &Service{}
}

// Hash хеширует пароль используя bcrypt
func (s *Service) Hash(password string, cost int) (string, error) {
	if cost == 0 {
		cost = BCryptCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Compare сравнивает пароль с хешем
func (s *Service) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// HashPassword хеширует пароль используя bcrypt (обратная совместимость)
func HashPassword(password string, cost int) (string, error) {
	if cost == 0 {
		cost = BCryptCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ComparePassword сравнивает пароль с хешем (обратная совместимость)
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
