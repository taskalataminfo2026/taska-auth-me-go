package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/taska-auth-me-go/internal/application/ports"
)

type AuthHandler struct {
	authService ports.AuthService
}

func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login maneja las solicitudes de inicio de sesión
// @Summary Iniciar sesión
// @Description Inicia sesión con correo electrónico y contraseña
// @Tags auth
// @Accept json
// @Produce json
// @Param input body loginRequest true "Credenciales de inicio de sesión"
// @Success 200 {object} domain.Token
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "solicitud inválida: " + err.Error()})
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	c.JSON(http.StatusOK, token)
}

// Logout maneja las solicitudes de cierre de sesión
// @Summary Cerrar sesión
// @Description Cierra la sesión del usuario actual
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token no proporcionado"})
		return
	}

	if err := h.authService.Logout(tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al cerrar sesión"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sesión cerrada correctamente"})
}

// RefreshToken maneja la renovación de tokens
// @Summary Renovar token
// @Description Renueva el token de acceso usando un token de actualización
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.Token
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token no proporcionado"})
		return
	}

	token, err := h.authService.RefreshToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
		return
	}

	c.JSON(http.StatusOK, token)
}

// RegisterRoutes registra las rutas de autenticación
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
		authGroup.POST("/logout", h.Logout)
		authGroup.POST("/refresh", h.RefreshToken)
	}
}
