package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/taska-auth-me-go/internal/interfaces/http/handlers"
	"github.com/taska-auth-me-go/internal/interfaces/http/middlewares"
)

// App representa la aplicación principal
type App struct {
	config         *config.Config
	router         *gin.Engine
	authHandler    *handlers.AuthHandler
	userHandler    *handlers.UserHandler
	authMiddleware *middlewares.AuthMiddleware
}

// NewApp crea una nueva instancia de la aplicación
func NewApp(
	cfg *config.Config,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware *middlewares.AuthMiddleware,
) *App {
	// Configurar Gin
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Configurar middlewares globales
	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.CORSMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	return &App{
		config:         cfg,
		router:         router,
		authHandler:    authHandler,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
	}
}

// SetupRoutes configura las rutas de la aplicación
func (a *App) SetupRoutes() {
	// Grupo de rutas de API
	api := a.router.Group("/api/v1")

	// Rutas públicas
	a.authHandler.RegisterRoutes(api)
	a.userHandler.RegisterRoutes(api)

	// Rutas protegidas que requieren autenticación
	// Ejemplo:
	// protected := api.Group("")
	// protected.Use(a.authMiddleware.HandlerFunc())
	// {
	// 	// Aquí irían las rutas protegidas
	// }
}

// Run inicia el servidor HTTP
func (a *App) Run() error {
	a.SetupRoutes()

	addr := ":" + a.config.Server.Port
	log.Printf("Servidor iniciado en http://localhost%s\n", addr)

	return a.router.Run(addr)
}
