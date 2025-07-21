package services

import (
	"errors"
	"time"

	"github.com/taska-auth-me-go/internal/domain"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Register(user *domain.User) error {
	// Verificar si el usuario ya existe
	existingUser, _ := s.userRepo.FindByEmail(user.Email)
	if existingUser != nil {
		return errors.New("el correo electrónico ya está en uso")
	}

	// Hashear la contraseña (implementar función de hash)
	// hashedPassword, err := hashPassword(user.Password)
	// if err != nil {
	// 	return err
	// }
	// user.Password = hashedPassword

	// Establecer fechas
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Crear usuario
	return s.userRepo.Create(user)
}

func (s *UserService) GetProfile(userID uint) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) UpdateProfile(user *domain.User) error {
	// Verificar que el usuario exista
	_, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	// Actualizar fecha de modificación
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(user)
}

func (s *UserService) DeleteAccount(userID uint) error {
	// Verificar que el usuario exista
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	return s.userRepo.Delete(userID)
}
