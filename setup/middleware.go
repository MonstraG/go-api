package setup

import (
	"go-server/setup/myJwt"
	"go-server/setup/myLog"
	"go-server/setup/reqRes"
	"net/http"
	"time"
)

// MyHandlerFunc is an alias for http.HandlerFunc argument, but with my reqRes.MyWriter and reqRes.MyRequest
type MyHandlerFunc func(w reqRes.MyWriter, r *reqRes.MyRequest)

// Middleware is just a MyHandlerFunc that returns a MyHandlerFunc
type Middleware func(MyHandlerFunc) MyHandlerFunc

func MyReqResWrapperMiddleware(next MyHandlerFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		myWriter := reqRes.MyWriter{ResponseWriter: w}
		myRequest := reqRes.MyRequest{Request: *r, RequestId: RandId()}
		next(myWriter, &myRequest)
	}
}

// LoggingMiddleware is a Middleware that logs a hit and time taken to answer
func LoggingMiddleware(next MyHandlerFunc) MyHandlerFunc {
	return func(w reqRes.MyWriter, r *reqRes.MyRequest) {
		start := time.Now()
		myLog.Info.Logf("{%s} Started %s %s", r.RequestId, r.Method, r.URL.Path)
		next(w, r)
		myLog.Info.Logf("{%s} Completed %s %s in %v", r.RequestId, r.Method, r.URL.Path, time.Since(start))
	}
}

func CreateJwtAuthRequiredMiddleware(jwtService *myJwt.Service) Middleware {
	return func(next MyHandlerFunc) MyHandlerFunc {
		return func(w reqRes.MyWriter, r *reqRes.MyRequest) {
			cookie, err := r.CookieIfValid(myJwt.Cookie)
			if err != nil {
				w.RedirectToLogin(r)
				return
			}

			claims, err := jwtService.ValidateJWT(cookie.Value)
			if err != nil {
				myLog.Info.Logf("Error validating JWT:\n%v\n", err)
				w.RedirectToLogin(r)
				return
			}

			r.UserId, err = claims.GetSubject()
			if err != nil {
				myLog.Info.Logf("Failed to get JWT subject, ignoring:\n%v\n", err)
			}

			r.Username = claims.Username

			next(w, r)
		}
	}
}
