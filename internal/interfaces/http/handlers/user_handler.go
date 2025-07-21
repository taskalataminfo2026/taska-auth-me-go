package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taska-auth-me-go/internal/application/ports"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

type registerRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// Register maneja las solicitudes de registro de usuarios
// @Summary Registrar un nuevo usuario
// @Description Crea una nueva cuenta de usuario
// @Tags users
// @Accept json
// @Produce json
// @Param input body registerRequest true "Datos del usuario"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "solicitud inválida: " + err.Error()})
		return
	}

	// Aquí deberías mapear el request al modelo de dominio
	// Esto es un ejemplo simplificado
	user := &domain.User{
		Email:     req.Email,
		Password:  req.Password, // La contraseña debería estar hasheada en el servicio
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.userService.Register(user); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "el correo electrónico ya está en uso" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "usuario registrado exitosamente"})
}

// GetProfile obtiene el perfil del usuario autenticado
// @Summary Obtener perfil de usuario
// @Description Obtiene la información del perfil del usuario autenticado
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	user, err := h.userService.GetProfile(userID.(uint))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuario no encontrado" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

type updateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// UpdateProfile actualiza el perfil del usuario autenticado
// @Summary Actualizar perfil de usuario
// @Description Actualiza la información del perfil del usuario autenticado
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body updateProfileRequest true "Datos a actualizar"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "solicitud inválida: " + err.Error()})
		return
	}

	// Obtener el usuario actual
	user, err := h.userService.GetProfile(userID.(uint))
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuario no encontrado" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	// Actualizar campos si se proporcionan
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := h.userService.UpdateProfile(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar el perfil"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteAccount elimina la cuenta del usuario autenticado
// @Summary Eliminar cuenta de usuario
// @Description Elimina la cuenta del usuario autenticado
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autorizado"})
		return
	}

	if err := h.userService.DeleteAccount(userID.(uint)); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "usuario no encontrado" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cuenta eliminada exitosamente"})
}

// RegisterRoutes registra las rutas de usuarios
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/register", h.Register)
		
		// Rutas protegidas que requieren autenticación
		authorized := usersGroup.Group("")
		authorized.Use(AuthMiddleware()) // Asegúrate de implementar este middleware
		{
			authorized.GET("/me", h.GetProfile)
			authorized.PUT("/me", h.UpdateProfile)
			authorized.DELETE("/me", h.DeleteAccount)
		}
	}
}
