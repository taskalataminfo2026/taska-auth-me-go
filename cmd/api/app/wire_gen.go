//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/taska-auth-me-go/internal/application/ports"
	"github.com/taska-auth-me-go/internal/application/services"
	"github.com/taska-auth-me-go/internal/domain"
	"github.com/taska-auth-me-go/internal/infrastructure/auth"
	"github.com/taska-auth-me-go/internal/infrastructure/persistence/gorm"
	"github.com/taska-auth-me-go/internal/interfaces/http/handlers"
	"github.com/taska-auth-me-go/internal/interfaces/http/middlewares"
	"gorm.io/gorm"
)

// InitializeApp inicializa la aplicación con todas sus dependencias
func InitializeApp() (*App, error) {
	wire.Build(
		// Configuración
		config.LoadConfig,

		// Infraestructura
		gorm.InitDB,
		gorm.NewUserRepository,
		newJWTManager,
		newPasswordHasher,

		// Servicios de aplicación
		services.NewAuthService,
		services.NewUserService,

		// Handlers HTTP
		handlers.NewAuthHandler,
		handlers.NewUserHandler,

		// Middlewares
		middlewares.NewAuthMiddleware,

		// Aplicación
		NewApp,
	)
	return &App{}, nil
}

// newJWTManager provee una instancia de JWTManager
func newJWTManager(cfg *config.Config) domain.TokenService {
	return auth.NewJWTManager(
		cfg.JWT.SecretKey,
		time.Duration(cfg.JWT.ExpirationHours)*time.Hour,
	)
}

// newPasswordHasher provee una instancia de PasswordHasher
func newPasswordHasher() *auth.PasswordHasher {
	return auth.NewPasswordHasher()
}

// newAuthMiddleware provee una instancia de AuthMiddleware
func newAuthMiddleware(authService ports.AuthService) *middlewares.AuthMiddleware {
	return middlewares.NewAuthMiddleware(authService)
}
