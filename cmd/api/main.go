//go:build !wireinject
// +build !wireinject

package main

import (
	"log"

	"github.com/google/wire"
	"github.com/taska-auth-me-go/internal/application/ports"
	"github.com/taska-auth-me-go/internal/application/services"
	"github.com/taska-auth-me-go/internal/domain"
	"github.com/taska-auth-me-go/internal/infrastructure/auth"
	"github.com/taska-auth-me-go/internal/infrastructure/persistence/gorm"
	"github.com/taska-auth-me-go/internal/interfaces/http/handlers"
	"github.com/taska-auth-me-go/internal/interfaces/http/middlewares"
	"github.com/taska-auth-me-go/pkg/config"
	"gorm.io/gorm"
)

func main() {
	// Inicializar la aplicación con inyección de dependencias
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Error al inicializar la aplicación: %v", err)
	}

	// Iniciar el servidor
	if err := app.Run(); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
