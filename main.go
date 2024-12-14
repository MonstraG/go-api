package main

import (
	"go-server/pages"
	"go-server/pages/index"
	"go-server/pages/login"
	"go-server/pages/music"
	"go-server/pages/music/websockets"
	"go-server/pages/notFound"
	"go-server/setup"
	"go-server/setup/appConfig"
	"log"
)

func main() {
	config := appConfig.ReadConfig()

	app := setup.NewApp(config)

	var authMiddleware = setup.CreateBasicAuthMiddleware(*app)
	app.Use(authMiddleware)

	app.Use(setup.LoggingMiddleware)

	mapRoutes(app)

	log.Fatal(app.ListenAndServe())
}

func mapRoutes(app *setup.App) {
	app.HandleFunc("GET /", notFound.GetHandler)
	app.HandleFunc("GET /{$}", index.GetHandler)
	app.HandleFunc("POST /music", music.PostHandler)
	app.HandleFunc("POST /ping", music.PongHandler)
	app.HandleFunc("GET /login", login.GetHandler)
	app.HandleFunc("POST /login", login.PostHandler)

	app.HandleFunc("GET /ws", websockets.HubSingleton.ServeWs)

	// resources
	app.HandleFunc("GET /public/{path...}", pages.PublicHandler)
}
