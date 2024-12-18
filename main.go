package main

import (
	"go-server/pages"
	"go-server/pages/index"
	"go-server/pages/login"
	"go-server/pages/music"
	"go-server/pages/notFound"
	"go-server/setup"
	"go-server/setup/appConfig"
	"log"
)

func main() {
	config := appConfig.ReadConfig()

	app := setup.NewApp(config)

	var authMiddleware = setup.CreateJwtAuthMiddleware(*app)
	app.Use(authMiddleware)

	app.Use(setup.LoggingMiddleware)

	mapRoutes(app)

	log.Fatal(app.ListenAndServe())
}

func mapRoutes(app *setup.App) {
	app.HandleFunc("GET /", notFound.GetHandler)
	app.HandleFunc("GET /{$}", index.GetHandler)
	app.HandleFunc("POST /songQueue", music.PostHandler)
	app.HandleFunc("GET /songQueue", music.GetSongQueueHandler)
	app.HandleFunc("GET /player", music.GetSongPlayerHandler)
	app.HandleFunc("GET /song/{id}", music.GetSongHandler)
	app.HandleFunc("GET /login", login.GetHandler)
	app.HandleFunc("POST /login", login.PostHandler)

	// resources
	app.HandleFunc("GET /public/{path...}", pages.PublicHandler)
}
