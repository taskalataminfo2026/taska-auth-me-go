package domain

// Token representa un token JWT en el sistema
type Token struct {
    AccessToken string `json:"access_token"`
    TokenType   string `json:"token_type"`
    ExpiresIn   int64  `json:"expires_in"`
}

// TokenService define los métodos para manejar tokens
type TokenService interface {
    CreateToken(userID uint) (*Token, error)
    ValidateToken(tokenString string) (uint, error)
    RefreshToken(tokenString string) (*Token, error)
}

// AuthUseCase define los casos de uso para la autenticación
type AuthUseCase interface {
    Login(email, password string) (*Token, error)
    Logout(tokenString string) error
    RefreshToken(tokenString string) (*Token, error)
    ValidateToken(tokenString string) (uint, error)
}
