package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

const BCryptCost = 12

// Service реализует хеширование паролей используя bcrypt.
type Service struct{}

// NewService создаёт новый сервис хеширования.
func NewService() *Service {
	return &Service{}
}

// Hash хеширует пароль используя bcrypt.
// Если cost == 0, используется BCryptCost по умолчанию.
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

// Compare сравнивает пароль с хешем.
func (s *Service) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
