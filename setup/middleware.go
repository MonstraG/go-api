package setup

import (
	"go-server/setup/myJwt"
	"go-server/setup/reqRes"
	"log"
	"net/http"
	"time"
)

// MyHandlerFunc is an alias for http.HandlerFunc argument, but with my reqRes.MyWriter and reqRes.MyRequest
type MyHandlerFunc func(w reqRes.MyWriter, r *reqRes.MyRequest)

// Middleware is just a MyHandlerFunc that returns a MyHandlerFunc
type Middleware func(MyHandlerFunc) MyHandlerFunc

func MyReqResWrapperMiddleware(next MyHandlerFunc, app *App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		myWriter := reqRes.MyWriter{ResponseWriter: w}
		myRequest := reqRes.MyRequest{Request: *r, Db: app.db, AppConfig: app.config}
		next(myWriter, &myRequest)
	}
}

// LoggingMiddleware is a Middleware that logs a hit and time taken to answer
func LoggingMiddleware(next MyHandlerFunc) MyHandlerFunc {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	return func(w reqRes.MyWriter, r *reqRes.MyRequest) {
		start := time.Now()
		id := RandId()
		log.Printf("{%s} Started %s %s", id, r.Method, r.URL.Path)
		next(w, r)
		log.Printf("{%s} Completed %s %s in %v", id, r.Method, r.URL.Path, time.Since(start))
	}
}

func CreateJwtAuthRequiredMiddleware(app App) Middleware {
	return func(next MyHandlerFunc) MyHandlerFunc {
		return func(w reqRes.MyWriter, r *reqRes.MyRequest) {
			cookie, err := r.CookieIfValid(myJwt.Cookie)
			if err != nil {
				w.RedirectToLogin()
				return
			}

			claims, err := app.MyJwt.ValidateJWT(cookie.Value)
			if err != nil {
				log.Printf("Error validating JWT:\n%v\n", err)
				w.RedirectToLogin()
				return
			}

			r.Username = claims.Username

			next(w, r)
		}
	}
}
