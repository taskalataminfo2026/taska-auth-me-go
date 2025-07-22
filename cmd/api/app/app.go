package app

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/taskalataminfo2026/taska-auth-me-go/cmd/api/app/providers"
)

var databaseSet = wire.NewSet(
	providers.DatabaseConnectionPostgres,
)

var ClientRouterSet = wire.NewSet()

var RustyClientSet = wire.NewSet()

var routerSet = wire.NewSet(
	ClientRouterSet,
	providers.ProviderRouter,
)

func Start() (*echo.Echo, error) {
	panic(wire.Build(
		databaseSet,
		routerSet,
		RustyClientSet,
	))
	return nil, nil
}
