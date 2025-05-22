package main

import (
	"go-server/pages"
	"go-server/pages/forgotPassword"
	"go-server/pages/index"
	"go-server/pages/login"
	"go-server/pages/logout"
	"go-server/pages/music"
	"go-server/pages/notFound"
	"go-server/setup"
	"go-server/setup/appConfig"
	"go-server/setup/myJwt"
	"go-server/setup/websockets"
	"log"
	"time"
)

func main() {
	config := appConfig.ReadConfig()

	app := setup.NewApp(config)

	app.Use(setup.LoggingMiddleware)

	mapRoutes(app)

	log.Fatal(app.ListenAndServe())
}

func mapRoutes(app *setup.App) {
	var jwtService = myJwt.CreateMyJwt(app.Config, time.Now)

	var authRequired = setup.CreateJwtAuthRequiredMiddleware(&jwtService)

	app.HandleFunc("GET /", notFound.GetHandler)
	app.HandleFunc("POST /", notFound.GetHandler)

	var forgotPasswordController = forgotPassword.NewController(app.Db)

	app.HandleFunc("GET /forgot-password", forgotPasswordController.GetHandler)
	app.HandleFunc("POST /forgot-password", forgotPasswordController.PostHandler)
	app.HandleFunc("POST /set-password", forgotPasswordController.PostSetPasswordHandler)

	var indexController = index.NewController(app.Config)
	app.HandleFunc("GET /{$}", authRequired(indexController.GetHandler))

	var loginController = login.NewController(&jwtService, app.Db)
	app.HandleFunc("GET /login", loginController.GetHandler)
	app.HandleFunc("POST /login", loginController.PostHandler)

	app.HandleFunc("GET /logout", logout.GetHandler)

	var musicController = music.NewController(app.Config)
	app.HandleFunc("GET /listSongs/{path...}", authRequired(musicController.GetSongs))
	app.HandleFunc("GET /song/{path...}", musicController.GetSongHandler)
	app.HandleFunc("PUT /song/{path...}", authRequired(musicController.PutSongHandler))
	app.HandleFunc("DELETE /song/{path...}", authRequired(musicController.DeleteSongHandler))
	app.HandleFunc("PUT /songFolder/{path...}", authRequired(musicController.CreateFolderHandler))

	app.HandleFunc("GET /public/{path...}", pages.PublicHandler)

	app.HandleFunc("GET /ws", websockets.HandleWebSocket)
}
