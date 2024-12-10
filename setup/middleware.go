package setup

import (
	"database/sql"
	"errors"
	"go-server/helpers"
	"go-server/models"
	"go-server/pages/notFound"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

// HandlerFn is an alias for http.HandlerFunc argument, but with my helpers.MyWriter
type HandlerFn func(w helpers.MyWriter, r *helpers.MyRequest)

// Middleware is just a HandlerFn that returns a HandlerFn
type Middleware func(HandlerFn) HandlerFn

func MyWriterWrapperMiddleware(next HandlerFn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		myWriter := helpers.MyWriter{ResponseWriter: w}
		myRequest := helpers.MyRequest{Request: *r}
		next(myWriter, &myRequest)
	}
}

// LoggingMiddleware is a Middleware that logs a hit and time taken to answer
func LoggingMiddleware(next HandlerFn) HandlerFn {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	return func(w helpers.MyWriter, r *helpers.MyRequest) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	}
}

// HtmxPartialMiddleware guards against direct browser navigations to partials
// It returns notFound if request wasn't made by htmx (Hx-Request header)
func HtmxPartialMiddleware(next HandlerFn) HandlerFn {
	return func(w helpers.MyWriter, r *helpers.MyRequest) {
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
		return func(w helpers.MyWriter, r *helpers.MyRequest) {
			username, password, ok := r.BasicAuth()
			if !ok {
				unauthorizedResponse(w)
				return
			}

			lowercaseUsername := strings.ToLower(username)
			user := models.User{}
			result := app.db.Where("lower(username) = @name", sql.Named("name", lowercaseUsername)).First(&user)
			if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Printf("Hit error when searching for user '%v':\n%v\n", lowercaseUsername, result.Error)
				w.WriteHeader(500)
				return
			}
			if result.RowsAffected == 0 {
				// didn't find them
				unauthorizedResponse(w)
				return
			}

			ok = user.CheckPasswordHash(password)

			if !ok {
				unauthorizedResponse(w)
				return
			}

			r.Username = username

			next(w, r)
		}
	}
}

func unauthorizedResponse(w helpers.MyWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="server", charset="UTF-8"`)
	w.WriteHeader(401)
}
