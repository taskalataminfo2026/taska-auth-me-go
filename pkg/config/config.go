package config

// Config contiene toda la configuración de la aplicación
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

// ServerConfig contiene la configuración del servidor HTTP
type ServerConfig struct {
    Port string
}

// DatabaseConfig contiene la configuración de la base de datos
type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

// JWTConfig contiene la configuración para JWT
type JWTConfig struct {
    SecretKey       string
    ExpirationHours int
}

// LoadConfig carga la configuración de las variables de entorno
func LoadConfig() (*Config, error) {
    // Por ahora devolvemos una configuración por defecto
    // Más adelante implementaremos la carga desde variables de entorno
    return &Config{
        Server: ServerConfig{
            Port: "8080",
        },
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     "5432",
            User:     "postgres",
            Password: "postgres",
            DBName:   "taska_auth",
            SSLMode:  "disable",
        },
        JWT: JWTConfig{
            SecretKey:       "your-secret-key-change-this-in-production",
            ExpirationHours: 24,
        },
    }, nil
}
