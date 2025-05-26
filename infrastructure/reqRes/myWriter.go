package reqRes

import (
	"bufio"
	"errors"
	"fmt"
	"go-api/infrastructure/myJwt"
	"go-api/infrastructure/myLog"
	"html/template"
	"net"
	"net/http"
)

type MyResponseWriter struct {
	http.ResponseWriter
}

func (myResponseWriter *MyResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := myResponseWriter.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func (myResponseWriter *MyResponseWriter) RenderTemplate(tmpl *template.Template, data any) bool {
	err := tmpl.Execute(myResponseWriter, data)
	if err != nil {
		message := fmt.Sprintf("Failed to render template: \n%v", err)
		myResponseWriter.Error(message, http.StatusInternalServerError)
		return false
	}
	return true
}

func (myResponseWriter *MyResponseWriter) RedirectToLogin(r *MyRequest) {
	http.Redirect(myResponseWriter, &r.Request, "/login", http.StatusTemporaryRedirect)
}

func (myResponseWriter *MyResponseWriter) Error(message string, code int) {
	if code >= 500 {
		myLog.Error.SkipLog(1, message)
	} else {
		myLog.Info.SkipLog(1, message)
	}
	http.Error(myResponseWriter, message, code)
}

func (myResponseWriter *MyResponseWriter) IssueCookie(value string, age int) {
	cookie := http.Cookie{
		Name:     myJwt.Cookie,
		Value:    value,
		Path:     "/",
		MaxAge:   age,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(myResponseWriter, &cookie)
}

func (myResponseWriter *MyResponseWriter) ExpireCookie() {
	myResponseWriter.IssueCookie("", -1)
}
