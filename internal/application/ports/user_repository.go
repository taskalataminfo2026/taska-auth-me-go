package ports

import "github.com/taska-auth-me-go/internal/domain"

// UserRepository es la interfaz que define el contrato para el repositorio de usuarios
type UserRepository interface {
    Create(user *domain.User) error
    FindByEmail(email string) (*domain.User, error)
    FindByID(id uint) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id uint) error
}
