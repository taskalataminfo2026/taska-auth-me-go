package config

import (
	"github.com/taskalataminfo2026/taska-auth-me-go/cmd/api/constants"
	"os"
	"time"
)

var (
	RustyConfig          RustyClientConfig
	DBConfig             ConnectionConfig
	MaxIdleConnections   int
	MaxOpenConnections   int
	ConnMaxLifetime      time.Duration
	ConnMaxIdleTime      time.Duration
	MaxBatchSize         int
	MaxConnectionRetries int
)

type ConnectionConfig struct {
	Username           string
	Password           string
	Host               string
	Name               string
	Port               string
	MaxIdleConnections int
	MaxOpenConnections int
	MaxBatchSize       int
	ConnMaxLifetime    time.Duration
	ConnMaxIdleTime    time.Duration
}

type RustyClientConfig struct {
	DefaultTimeOut time.Duration
	RetryCount     int
}

func init() {

	// DB.
	MaxIdleConnections = 500
	MaxOpenConnections = 500
	ConnMaxLifetime = 600 * time.Second
	ConnMaxIdleTime = 600 * time.Second
	MaxBatchSize = 100
	MaxConnectionRetries = 3

	// Rusty client
	RustyConfig.DefaultTimeOut = 11 * time.Second
	RustyConfig.RetryCount = 3

	if os.Getenv("GO_ENVIRONMENT") == "" ||
		os.Getenv("GO_ENVIRONMENT") == "test" ||
		os.Getenv("GO_ENVIRONMENT") == constants.ScopeLocal {

		DBConfig = ConnectionConfig{
			Username:           "postgres",
			Password:           "root",
			Host:               "localhost",
			Name:               "local_tareaya",
			Port:               "5432",
			MaxIdleConnections: MaxIdleConnections,
			MaxOpenConnections: MaxOpenConnections,
			ConnMaxLifetime:    time.Second * ConnMaxLifetime,
			ConnMaxIdleTime:    time.Second * ConnMaxIdleTime,
			MaxBatchSize:       MaxBatchSize,
		}
	}

	if os.Getenv("GO_ENVIRONMENT") == constants.ScopeBeta {
		DBConfig = ConnectionConfig{
			Username:           "beta_user",
			Password:           "beta_password",
			Host:               "beta-host:3306",
			Name:               "beta_database",
			Port:               "5432",
			MaxIdleConnections: MaxIdleConnections,
			MaxOpenConnections: MaxOpenConnections,
			ConnMaxLifetime:    time.Second * ConnMaxLifetime,
			ConnMaxIdleTime:    time.Second * ConnMaxIdleTime,
			MaxBatchSize:       MaxBatchSize,
		}
	}

	if os.Getenv("GO_ENVIRONMENT") == constants.ScopeProd {
		DBConfig = ConnectionConfig{
			Username:           "prod_user",
			Password:           "prod_password",
			Host:               "prod-host:3306",
			Name:               "prod_database",
			Port:               "5432",
			MaxIdleConnections: MaxIdleConnections,
			MaxOpenConnections: MaxOpenConnections,
			ConnMaxLifetime:    time.Second * ConnMaxLifetime,
			ConnMaxIdleTime:    time.Second * ConnMaxIdleTime,
			MaxBatchSize:       MaxBatchSize,
		}
	}
}
