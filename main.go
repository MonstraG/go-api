package main

import (
	"go-server/pages"
	"go-server/setup"
	"log"
)

func main() {
	config := setup.ReadConfig()

	app := setup.NewApp(config)

	var authMiddleware = setup.CreateBasicAuthMiddleware(*app)
	app.Use(authMiddleware)

	app.Use(setup.LoggingMiddleware)

	pages.MapRoutes(app)

	log.Fatal(app.ListenAndServe())
}
