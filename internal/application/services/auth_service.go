package services

import (
	"errors"
	"time"

	"github.com/taska-auth-me-go/internal/domain"
)

type AuthService struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenService
}

func NewAuthService(userRepo domain.UserRepository, tokenRepo domain.TokenService) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *AuthService) Login(email, password string) (*domain.Token, error) {
	// Buscar usuario por email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar contraseña (implementar función de verificación)
	// if !checkPasswordHash(password, user.Password) {
	// 	return nil, errors.New("credenciales inválidas")
	// }

	// Generar token
	token, err := s.tokenRepo.CreateToken(user.ID)
	if err != nil {
		return nil, errors.New("error al generar el token")
	}

	return token, nil
}

func (s *AuthService) Logout(tokenString string) error {
	// En una implementación real, podrías invalidar el token aquí
	return nil
}

func (s *AuthService) RefreshToken(tokenString string) (*domain.Token, error) {
	return s.tokenRepo.RefreshToken(tokenString)
}

func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	return s.tokenRepo.ValidateToken(tokenString)
}
