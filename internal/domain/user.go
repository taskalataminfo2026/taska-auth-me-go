package domain

import "time"

// User representa un usuario en el sistema
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository define los métodos que debe implementar un repositorio de usuarios
type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

// UserUseCase define los casos de uso para los usuarios
type UserUseCase interface {
	Register(user *User) error
	Login(email, password string) (string, error)
	GetProfile(userID uint) (*User, error)
	UpdateProfile(user *User) error
	DeleteAccount(userID uint) error
}
