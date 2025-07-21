package gorm

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/taska-auth-me-go/pkg/config"
)

// InitDB inicializa y devuelve una instancia de conexión a la base de datos
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Construir la cadena de conexión
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	// Conectar a la base de datos
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar a la base de datos: %v", err)
	}

	// Configurar el pool de conexiones
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error al acceder a la conexión de base de datos: %v", err)
	}

	// Configurar el pool de conexiones
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Verificar la conexión
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error al hacer ping a la base de datos: %v", err)
	}

	log.Println("Conexión a la base de datos establecida correctamente")

	return db, nil
}

// AutoMigrate realiza las migraciones automáticamente
func AutoMigrate(db *gorm.DB) error {
	// Agrega aquí tus modelos para migrar
	err := db.AutoMigrate(
		// &domain.User{},
		// Agrega más modelos aquí según sea necesario
	)

	if err != nil {
		return fmt.Errorf("error al realizar migraciones: %v", err)
	}

	log.Println("Migraciones completadas exitosamente")
	return nil
}
