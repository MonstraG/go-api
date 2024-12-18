package setup

import (
	"go-server/pages/notFound"
	"go-server/setup/myJwt"
	"go-server/setup/reqRes"
	"log"
	"net/http"
	"strings"
	"time"
)

// HandlerFn is an alias for http.HandlerFunc argument, but with my ytDlp.MyWriter
type HandlerFn func(w reqRes.MyWriter, r *reqRes.MyRequest)

// Middleware is just a HandlerFn that returns a HandlerFn
type Middleware func(HandlerFn) HandlerFn

func MyReqResWrapperMiddleware(next HandlerFn, app *App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		myWriter := reqRes.MyWriter{ResponseWriter: w}
		myRequest := reqRes.MyRequest{Request: *r, Db: app.db, AppConfig: app.config}
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

			cookie, err := r.CookieIfValid(myJwt.Cookie)
			if err != nil {
				redirectToLogin(w)
				return
			}

			_, err = myJwt.Singleton.ValidateJWT(cookie.Value, app.config)
			if err != nil {
				log.Printf("Error validating JWT:\n%v\n", err)
				redirectToLogin(w)
				return
			}

			next(w, r)
		}
	}
}

func redirectToLogin(w reqRes.MyWriter) {
	w.Header().Set("Location", `/login`)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
