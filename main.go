package main

import (
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/setup"
)

func main() {
	config := appConfig.ReadConfig()

	app := setup.NewApp(config)

	app.Use(setup.LoggingMiddleware)
	app.Use(setup.VersionMiddleware)
	err := app.MapRoutes()
	if err != nil {
		myLog.Fatal.Logf("Failed to create services/routes: %v", err.Error())
	}

	err = app.ListenAndServe()
	if err != nil {
		myLog.Fatal.Logf("Server exited with error: %v", err.Error())
	}
}
