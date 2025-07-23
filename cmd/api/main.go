package main

import (
	"github.com/taskalataminfo2026/tool-kit-lib-go/pkg/loggers"
	"os"

	"github.com/taskalataminfo2026/taska-auth-me-go/cmd/api/app"
)

// @title Taska API
// @version 1.0
// @description Taska API for use with his admin project.
// @contact.email wilsonev.saldarriaga88@gmail.com
func main() {
	appInstance, err := app.Start()
	if err != nil {
		loggers.Error("Failed to initialize application", err)
		os.Exit(1)
	}

	port := ":8080"
	loggers.Info("Starting server", "port", port)

	if err := appInstance.Start(port); err != nil {
		loggers.Error("Server failed to start", err)
		os.Exit(1)
	}
	loggers.Info("Server started successfully", "port", port)
}
