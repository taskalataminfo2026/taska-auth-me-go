package providers

import (
	"fmt"
	"github.com/taskalataminfo2026/taska-auth-me-go/cmd/api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func DatabaseConnectionPostgres() (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	for retry := 0; retry < config.MaxConnectionRetries; retry++ {
		if db, err = GetDBConnectionPostgres(); err != nil {
			continue
		}

		// Probar conexión
		sqlConn, _ := db.DB()
		if err = sqlConn.Ping(); err == nil {
			break
		}
	}

	return db, err
}

func GetDBConnectionPostgres() (*gorm.DB, error) {
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.DBConfig.Host,
		config.DBConfig.Username,
		config.DBConfig.Password,
		config.DBConfig.Name,
		config.DBConfig.Port)

	conn, err := gorm.Open(postgres.Open(connString), &gorm.Config{PrepareStmt: true, QueryFields: true})
	if err != nil {
		//utils.GetLogger(context.TODO()).Error(fmt.Sprintf("Error abriendo la conexión a la base de datos: %v", err))
		return nil, err
	}

	sqlDB, err := conn.DB()
	if err != nil {
		//utils.GetLogger(context.TODO()).Error(fmt.Sprintf("Error obteniendo la base de datos: %v", err))
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(config.DBConfig.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.DBConfig.ConnMaxIdleTime)
	sqlDB.SetMaxIdleConns(config.DBConfig.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(config.DBConfig.MaxOpenConnections)

	//utils.GetLogger(context.TODO()).Info(fmt.Sprintf("Conexiones abiertas: %v", sqlDB.Stats().OpenConnections))

	env := os.Getenv("GO_ENVIRONMENT")
	if env == "test" || env == "" {
		//test_utils.CreateTestDatabase(conn)
	}

	return conn, err
}
