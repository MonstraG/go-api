package setup

import (
	"fmt"
	"go-server/setup/appConfig"
	"go-server/setup/myJwt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// App = http.ServeMux + Middleware slice
type App struct {
	mux         *http.ServeMux
	middlewares []Middleware
	config      appConfig.AppConfig
	db          *gorm.DB
	MyJwt       myJwt.Service
}

// NewApp is a constructor for App
func NewApp(appConfig appConfig.AppConfig) *App {
	db := OpenDb(appConfig)

	return &App{
		mux:         http.NewServeMux(),
		middlewares: []Middleware{},
		config:      appConfig,
		db:          db,
		MyJwt:       myJwt.CreateMyJwt(appConfig, time.Now),
	}
}

// Use adds Middleware to chain
func (app *App) Use(mw Middleware) {
	app.middlewares = append(app.middlewares, mw)
}

// HandleFunc is a wrapper around normal http.HandleFunc but calling all Middleware-s first
func (app *App) HandleFunc(pattern string, handlerFunc MyHandlerFunc) {
	app.mux.HandleFunc(pattern, MyReqResWrapperMiddleware(applyMiddlewares(handlerFunc, app.middlewares), app))
}

// applyMiddlewares runs all middlewares in order
func applyMiddlewares(handlerFunc MyHandlerFunc, middlewares []Middleware) MyHandlerFunc {
	for _, middleware := range middlewares {
		handlerFunc = middleware(handlerFunc)
	}
	return handlerFunc
}

// ListenAndServe is a wrapper around normal http.ListenAndServe
func (app *App) ListenAndServe() error {
	log.Println(fmt.Sprintf("Starting server on %s", app.config.Host))

	server := &http.Server{
		Addr:         app.config.Host,
		Handler:      app.mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return server.ListenAndServe()
}
