package setup

import (
	"go-api/infrastructure/helpers"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/myToken"
	"go-api/infrastructure/reqRes"
	"go-api/infrastructure/version"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// MyHandlerFunc is an alias for http.HandlerFunc, but with my reqRes.MyResponseWriter and reqRes.MyRequest
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

func parseToken(myTokenService *myToken.Service, w reqRes.MyResponseWriter, r *reqRes.MyRequest) *myToken.TokenPayload {
	cookie, err := r.CookieIfValid(myToken.Cookie)
	if err != nil {
		w.RedirectToLogin(r)
		return nil
	}

	payload, err := myTokenService.ParseCookie(cookie.Value)
	if err != nil {
		myLog.Info.Logf("Error parsing token:\n\t%v", err)
		w.RedirectToLogin(r)
		return nil
	}
	return &payload
}

func newAuthRequiredMiddleware(myTokenService *myToken.Service, db *gorm.DB) Middleware {
	return func(next MyHandlerFunc) MyHandlerFunc {
		return func(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
			token := parseToken(myTokenService, w, r)
			if token == nil {
				return
			}

			refreshedToken, err := myTokenService.RefreshTokenIfOld(*token, db)
			if err != nil {
				myLog.Error.Logf("Error refreshing token:\n\t%v", err)
				// the most probable cause is that user was deleted, old token should not be used in this case anyway
				w.Error("Internal server error", http.StatusInternalServerError)
				return
			}

			if refreshedToken != *token {
				cookieValue, err := myTokenService.CreateCookie(refreshedToken)
				if err != nil {
					myLog.Error.Logf("Error creating cookie for refreshed token:\n\t%v", err)
					w.Error("Internal server error", http.StatusInternalServerError)
					return
				}

				w.IssueCookie(cookieValue, myToken.DefaultCookieAge)
			}

			r.Token = refreshedToken

			next(w, r)
		}
	}
}

func newAdminRequiredMiddleware(authRequiredMiddleware Middleware) Middleware {
	return func(next MyHandlerFunc) MyHandlerFunc {
		return func(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
			authRequiredMiddleware(next)

			if !r.Token.IsAdmin {
				w.Error("Forbidden", http.StatusForbidden)
				return
			}

			next(w, r)
		}
	}
}
