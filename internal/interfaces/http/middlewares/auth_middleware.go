package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taska-auth-me-go/internal/application/ports"
)

// AuthMiddleware es un middleware que verifica la autenticación del usuario
func AuthMiddleware(authService ports.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del encabezado de autorización
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "se requiere token de autenticación"})
			c.Abort()
			return
		}

		// El formato del encabezado debe ser "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar el token
		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido o expirado"})
			c.Abort()
			return
		}

		// Agregar el ID del usuario al contexto para que los controladores puedan acceder a él
		c.Set("userID", userID)
		c.Next()
	}
}

// CORSMiddleware permite solicitudes CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ErrorHandler es un middleware para manejar errores
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Verificar si hay errores
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			switch err.Type {
			case gin.ErrorTypeBind:
				c.JSON(http.StatusBadRequest, gin.H{"error": "solicitud inválida: " + err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error interno del servidor"})
			}
		}
	}
}
