package ports

import "github.com/taska-auth-me-go/internal/domain"

// AuthService es la interfaz que define el contrato para el servicio de autenticación
type AuthService interface {
    Login(email, password string) (*domain.Token, error)
    Logout(tokenString string) error
    RefreshToken(tokenString string) (*domain.Token, error)
    ValidateToken(tokenString string) (uint, error)
}
