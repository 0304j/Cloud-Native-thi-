package service

import (
	"auth-service/internal/domain"
	"auth-service/internal/ports"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      ports.UserRepository
	jwtSecret string
	expiry    time.Duration
	JWTSecret string
}

func (s *AuthService) Authenticate(email, password string) (*domain.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func NewAuthService(repo ports.UserRepository, jwtSecret string, expirySeconds int) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		expiry:    time.Duration(expirySeconds) * time.Second,
	}
}

func (s *AuthService) Register(email, password, role string) error {
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hash),
		Role:     role, // <-- Rolle aus Parameter!
	}
	return s.repo.Create(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(s.expiry).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}
