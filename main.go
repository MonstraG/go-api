package main

import (
	"go-server/pages"
	"go-server/pages/index"
	"go-server/pages/notFound"
	"go-server/setup"
	"log"
)

func main() {
	config := setup.ReadConfig()

	app := setup.NewApp(config)

	var authMiddleware = setup.CreateBasicAuthMiddleware(*app)
	app.Use(authMiddleware)

	app.Use(setup.LoggingMiddleware)

	// pages
	app.HandleFunc("GET /", notFound.GetHandler)
	app.HandleFunc("GET /{$}", index.GetHandler)

	// resources
	app.HandleFunc("GET /public/{path...}", pages.PublicHandler)

	log.Fatal(app.ListenAndServe())
}
