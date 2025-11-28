package setup

import (
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/myJwt"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/websockets"
	"go-api/pages"
	"go-api/pages/fileExplorer"
	"go-api/pages/forgotPassword"
	"go-api/pages/index"
	"go-api/pages/login"
	"go-api/pages/logout"
	"go-api/pages/notFound"
	"go-api/pages/player"
	"go-api/pages/users"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

// App = http.ServeMux + Middleware slice
type App struct {
	mux         *http.ServeMux
	middlewares []Middleware
	Config      appConfig.AppConfig
	Db          *gorm.DB
}

// NewApp is a constructor for App
func NewApp(appConfig appConfig.AppConfig) *App {
	db := OpenDb(appConfig)

	err := os.MkdirAll(appConfig.ExplorerRoot, 0766)
	if err != nil {
		myLog.Fatal.Logf("Failed to ensure explorer root folder exists")
	}

	return &App{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
		Config:      appConfig,
		Db:          db,
	}
}

// Use adds Middleware to chain
func (app *App) Use(mw Middleware) {
	app.middlewares = append(app.middlewares, mw)
}

// handleFunc is a wrapper around normal http.HandleFunc but calling all Middleware-s first
func (app *App) handleFunc(pattern string, handlerFunc MyHandlerFunc) {
	app.mux.HandleFunc(pattern, myReqResWrapperMiddleware(applyMiddlewares(handlerFunc, app.middlewares)))
}

// applyMiddlewares runs all middlewares in order
func applyMiddlewares(handlerFunc MyHandlerFunc, middlewares []Middleware) MyHandlerFunc {
	for _, middleware := range middlewares {
		handlerFunc = middleware(handlerFunc)
	}
	return handlerFunc
}

func (app *App) MapRoutes() {
	jwtService := myJwt.CreateMyJwt(app.Config, time.Now)

	// todo: see if I want to create db context per-request
	authRequired := createJwtAuthRequiredMiddleware(&jwtService, app.Db)
	adminRequired := createAdminRequiredMiddleware(&jwtService, app.Db)

	app.handleFunc("GET /", notFound.Show404)
	app.handleFunc("POST /", notFound.Show404)

	forgotPasswordController := forgotPassword.NewController(app.Db)

	app.handleFunc("GET /forgot-password", forgotPasswordController.GetForgotPasswordForm)
	app.handleFunc("POST /forgot-password", forgotPasswordController.SubmitForgotPasswordForm)
	app.handleFunc("POST /set-password", forgotPasswordController.SetPassword)

	indexController := index.NewController(app.Config)
	app.handleFunc("GET /{$}", authRequired(indexController.GetHandler))

	loginController := login.NewController(&jwtService, app.Db)
	app.handleFunc("GET /login", loginController.GetHandler)
	app.handleFunc("POST /login", loginController.PostHandler)

	app.handleFunc("GET /logout", logout.GetHandler)

	explorerController := fileExplorer.NewController(app.Config)
	app.handleFunc("GET /exploreAt/{path...}", authRequired(explorerController.ExploreAt))
	app.handleFunc("GET /file/{path...}", explorerController.GetFile)
	app.handleFunc("PUT /file/{path...}", authRequired(explorerController.PutFile))
	app.handleFunc("DELETE /file/{path...}", authRequired(explorerController.DeleteFile))
	app.handleFunc("PUT /directory/{path...}", authRequired(explorerController.PutDirectory))

	playerController := player.NewController(app.Config, app.Db)
	app.handleFunc("GET /player", authRequired(playerController.GetPlayer))
	app.handleFunc("POST /enqueueSong/{path...}", authRequired(playerController.EnqueueSong))
	app.handleFunc("POST /enqueueFolder/{path...}", authRequired(playerController.EnqueueFolder))
	app.handleFunc("DELETE /removeSong/{id}", authRequired(playerController.RemoveSong))
	app.handleFunc("POST /reportSongDuration/{queuedSongId}", authRequired(playerController.ReportSongDuration))

	usersController := users.NewController(app.Db)
	app.handleFunc("PUT /users/setPasswordChangeStatus", adminRequired(usersController.SetPasswordChangeStatus))

	app.handleFunc("GET /public/{path...}", pages.PublicHandler)

	app.handleFunc("GET /ws", websockets.HandleWebSocket)
}

// ListenAndServe is a wrapper around normal http.ListenAndServe
func (app *App) ListenAndServe() error {
	myLog.Info.Logf("Starting server on %s", app.Config.Host)

	server := &http.Server{
		Addr:         app.Config.Host,
		Handler:      app.mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return server.ListenAndServe()
}
