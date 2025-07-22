package providers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func ProviderRouter() *echo.Echo {
	router := echo.New()

	router.GET("/swagger/*", echoSwagger.WrapHandler)

	router.Use(middleware.Recover())
	router.Use(middleware.Logger())

	//api := router.Group("/v1/api")
	{
		//users := api.Group("/users	")

	}
	return router
}
