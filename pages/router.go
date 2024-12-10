package pages

import (
	"go-server/pages/index"
	"go-server/pages/login"
	"go-server/pages/notFound"
	"go-server/setup"
)

func MapRoutes(app *setup.App) {
	// pages
	app.HandleFunc("GET /", notFound.GetHandler)
	app.HandleFunc("GET /{$}", index.GetHandler)
	app.HandleFunc("GET /login", login.GetHandler)
	app.HandleFunc("POST /login", login.PostHandler)

	// resources
	app.HandleFunc("GET /public/{path...}", PublicHandler)
}
