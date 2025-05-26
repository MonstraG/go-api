package setup

import (
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/myJwt"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/websockets"
	"go-api/pages"
	"go-api/pages/forgotPassword"
	"go-api/pages/index"
	"go-api/pages/login"
	"go-api/pages/logout"
	"go-api/pages/music"
	"go-api/pages/notFound"
	"gorm.io/gorm"
	"net/http"
	"time"
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
	var jwtService = myJwt.CreateMyJwt(app.Config, time.Now)

	var authRequired = createJwtAuthRequiredMiddleware(&jwtService)

	app.handleFunc("GET /", notFound.GetHandler)
	app.handleFunc("POST /", notFound.GetHandler)

	var forgotPasswordController = forgotPassword.NewController(app.Db)

	app.handleFunc("GET /forgot-password", forgotPasswordController.GetHandler)
	app.handleFunc("POST /forgot-password", forgotPasswordController.PostHandler)
	app.handleFunc("POST /set-password", forgotPasswordController.PostSetPasswordHandler)

	var indexController = index.NewController(app.Config)
	app.handleFunc("GET /{$}", authRequired(indexController.GetHandler))

	var loginController = login.NewController(&jwtService, app.Db)
	app.handleFunc("GET /login", loginController.GetHandler)
	app.handleFunc("POST /login", loginController.PostHandler)

	app.handleFunc("GET /logout", logout.GetHandler)

	var musicController = music.NewController(app.Config)
	app.handleFunc("GET /listSongs/{path...}", authRequired(musicController.GetSongs))
	app.handleFunc("GET /song/{path...}", musicController.GetSongHandler)
	app.handleFunc("PUT /song/{path...}", authRequired(musicController.PutSongHandler))
	app.handleFunc("DELETE /song/{path...}", authRequired(musicController.DeleteSongHandler))
	app.handleFunc("PUT /songFolder/{path...}", authRequired(musicController.CreateFolderHandler))

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
