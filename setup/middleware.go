package setup

import (
	"database/sql"
	"errors"
	"go-server/models"
	"go-server/pages/notFound"
	"go-server/setup/reqRes"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

// HandlerFn is an alias for http.HandlerFunc argument, but with my helpers.MyWriter
type HandlerFn func(w reqRes.MyWriter, r *reqRes.MyRequest)

// Middleware is just a HandlerFn that returns a HandlerFn
type Middleware func(HandlerFn) HandlerFn

func MyReqResWrapperMiddleware(next HandlerFn, app *App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		myWriter := reqRes.MyWriter{ResponseWriter: w}
		myRequest := reqRes.MyRequest{Request: *r, Db: app.db}
		next(myWriter, &myRequest)
	}
}

// LoggingMiddleware is a Middleware that logs a hit and time taken to answer
func LoggingMiddleware(next HandlerFn) HandlerFn {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	return func(w reqRes.MyWriter, r *reqRes.MyRequest) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// HtmxPartialMiddleware guards against direct browser navigations to partials
// It returns notFound if request wasn't made by htmx (Hx-Request header)
func HtmxPartialMiddleware(next HandlerFn) HandlerFn {
	return func(w reqRes.MyWriter, r *reqRes.MyRequest) {
		isHtmxRequest := r.Header.Get("Hx-Request") == "true"
		if !isHtmxRequest {
			notFound.GetHandler(w, r)
			return
		}

		next(w, r)
	}
}

// CreateBasicAuthMiddleware returns middleware that requires basic auth
func CreateBasicAuthMiddleware(app App) Middleware {
	return func(next HandlerFn) HandlerFn {
		return func(w reqRes.MyWriter, r *reqRes.MyRequest) {
			if r.URL.Path == "/login" ||
				strings.HasPrefix(r.URL.Path, "/public") {
				next(w, r)
				return
			}

			username, password, ok := r.BasicAuth()
			if !ok {
				redirectToLogin(w)
				return
			}

			lowercaseUsername := strings.ToLower(username)
			user := models.User{}
			result := app.db.Where("lower(username) = @name", sql.Named("name", lowercaseUsername)).First(&user)
			if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Printf("Hit error when searching for user '%v':\n%v\n", lowercaseUsername, result.Error)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if result.RowsAffected == 0 {
				// didn't find them
				redirectToLogin(w)
				return
			}

			ok = user.CheckPasswordHash(password)

			if !ok {
				redirectToLogin(w)
				return
			}

			r.Username = username

			next(w, r)
		}
	}
}

// todo: remember url?
func redirectToLogin(w reqRes.MyWriter) {
	w.Header().Set("Location", `/login`)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
