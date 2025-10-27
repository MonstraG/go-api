package setup

import (
	"go-api/infrastructure/helpers"
	"go-api/infrastructure/myJwt"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/reqRes"
	"go-api/infrastructure/version"
	"net/http"
	"time"
)

// MyHandlerFunc is an alias for http.HandlerFunc argument, but with my reqRes.MyResponseWriter and reqRes.MyRequest
type MyHandlerFunc func(w reqRes.MyResponseWriter, r *reqRes.MyRequest)

// Middleware is just a MyHandlerFunc that returns a MyHandlerFunc
type Middleware func(MyHandlerFunc) MyHandlerFunc

func myReqResWrapperMiddleware(next MyHandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		myWriter := reqRes.MyResponseWriter{ResponseWriter: w, MyMetadata: &reqRes.MyMetadata{}}
		myRequest := reqRes.MyRequest{Request: *r, RequestId: helpers.RandId()}
		next(myWriter, &myRequest)
	}
}

// LoggingMiddleware is a Middleware that logs a hit and time taken to answer
func LoggingMiddleware(next MyHandlerFunc) MyHandlerFunc {
	return func(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
		start := time.Now()
		myLog.Info.Logf("{%s} Started %s %s", r.RequestId, r.Method, r.URL.Path)
		next(w, r)
		if w.StatusCode == 0 {
			w.StatusCode = 200
		}
		myLog.Info.Logf("{%s} Responded %d to %s %s in %v", r.RequestId, w.StatusCode, r.Method, r.URL.Path, time.Since(start))
	}
}

func VersionMiddleware(next MyHandlerFunc) MyHandlerFunc {
	return func(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
		w.Header().Set("X-Version", version.AppVersion)
		next(w, r)
	}
}

func createJwtAuthRequiredMiddleware(jwtService *myJwt.Service) Middleware {
	return func(next MyHandlerFunc) MyHandlerFunc {
		return func(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
			cookie, err := r.CookieIfValid(myJwt.Cookie)
			if err != nil {
				w.RedirectToLogin(r)
				return
			}

			claims, err := jwtService.ValidateJWT(cookie.Value)
			if err != nil {
				myLog.Info.Logf("Error validating JWT:\n\t%v", err)
				w.RedirectToLogin(r)
				return
			}

			r.UserId, err = claims.GetSubject()
			if err != nil {
				myLog.Info.Logf("Failed to get JWT subject, ignoring:\n\t%v", err)
			}

			r.Username = claims.Username

			next(w, r)
		}
	}
}
