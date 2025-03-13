package main

import (
	"go-server/pages"
	"go-server/pages/index"
	"go-server/pages/login"
	"go-server/pages/logout"
	"go-server/pages/music"
	"go-server/pages/notFound"
	"go-server/setup"
	"go-server/setup/appConfig"
	"go-server/setup/websockets"
	"log"
)

func main() {
	config := appConfig.ReadConfig()

	app := setup.NewApp(config)

	app.Use(setup.LoggingMiddleware)

	mapRoutes(app)

	log.Fatal(app.ListenAndServe())
}

func mapRoutes(app *setup.App) {
	var authRequired = setup.CreateJwtAuthRequiredMiddleware(*app)

	app.HandleFunc("GET /", authRequired(notFound.GetHandler))

	app.HandleFunc("GET /{$}", authRequired(index.GetHandler))

	app.HandleFunc("GET /login", login.GetHandler)
	app.HandleFunc("POST /login", login.PostHandler)
	app.HandleFunc("GET /logout", logout.GetHandler)

	app.HandleFunc("GET /listSongs/{path...}", authRequired(music.GetSongs))

	app.HandleFunc("GET /song/{path...}", music.GetSongHandler)
	app.HandleFunc("PUT /song/{path...}", authRequired(music.PutSongHandler))
	app.HandleFunc("DELETE /song/{path...}", authRequired(music.DeleteSongHandler))

	app.HandleFunc("GET /public/{path...}", pages.PublicHandler)

	app.HandleFunc("GET /ws", websockets.HandleWebSocket)
}
