package main

import (
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/setup"
	"runtime/debug"
)

func main() {
	config := appConfig.ReadConfig()

	app := setup.NewApp(config)

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		myLog.Fatal.Logf("Failed to read build info")
		return
	}

	myLog.Info.Logf("Version: %s", buildInfo.Main.Version)

	app.Use(setup.LoggingMiddleware)
	app.Use(setup.VersionMiddleware(buildInfo.Main.Version))
	app.MapRoutes()

	err := app.ListenAndServe()
	if err != nil {
		myLog.Fatal.Logf("%v", err.Error())
	}
}
