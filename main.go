package main

import (
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/myJwt"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/setup"
	"go-api/infrastructure/websockets"
	"go-api/pages"
	"go-api/pages/forgotPassword"
	"go-api/pages/index"
	"go-api/pages/login"
	"go-api/pages/logout"
	"go-api/pages/music"
	"go-api/pages/notFound"
	"time"
)

func main() {
	config := appConfig.ReadConfig()

	app := setup.NewApp(config)

	app.Use(setup.LoggingMiddleware)

	mapRoutes(app)

	err := app.ListenAndServe()
	if err != nil {
		myLog.Fatal.Log(err.Error())
	}
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
