package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/taska-auth-me-go/internal/domain"
)

// JWTManager implementa la interfaz TokenService del dominio
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// UserClaims contiene los claims personalizados para el token JWT
type UserClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// NewJWTManager crea una nueva instancia de JWTManager
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

// CreateToken genera un nuevo token JWT para el usuario
func (m *JWTManager) CreateToken(userID uint) (*domain.Token, error) {
	// Crear claims
	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "taska-auth-service",
		},
	}

	// Crear token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token con la clave secreta
	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return nil, err
	}

	return &domain.Token{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   int64(m.tokenDuration.Seconds()),
	}, nil
}

// ValidateToken valida el token JWT y devuelve el ID del usuario
func (m *JWTManager) ValidateToken(tokenString string) (uint, error) {
	// Parsear el token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validar el algoritmo de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("método de firma de token inválido")
			}
			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		return 0, err
	}

	// Validar claims
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("token inválido")
}

// RefreshToken renueva un token JWT existente
func (m *JWTManager) RefreshToken(tokenString string) (*domain.Token, error) {
	// Primero validamos el token existente
	userID, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Si el token es válido, generamos uno nuevo
	return m.CreateToken(userID)
}
