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
	app.MapRoutes()

	err := app.ListenAndServe()
	if err != nil {
		myLog.Fatal.Logf(err.Error())
	}
}
