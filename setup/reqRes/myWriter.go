package reqRes

import (
	"bufio"
	"errors"
	"fmt"
	"go-server/setup/myJwt"
	"go-server/setup/myLog"
	"html/template"
	"net"
	"net/http"
)

type MyWriter struct {
	http.ResponseWriter
}

func (myWriter *MyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := myWriter.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func (myWriter *MyWriter) RenderTemplate(tmpl *template.Template, data any) bool {
	err := tmpl.Execute(myWriter, data)
	if err != nil {
		message := fmt.Sprintf("Failed to render template: \n%v", err)
		myWriter.Error(message, http.StatusInternalServerError)
		return false
	}
	return true
}

func (myWriter *MyWriter) RedirectToLogin(r *MyRequest) {
	http.Redirect(myWriter, &r.Request, "/login", http.StatusTemporaryRedirect)
}

func (myWriter *MyWriter) Error(message string, code int) {
	if code >= 500 {
		myLog.Error.SkipLog(1, message)
	} else {
		myLog.Info.SkipLog(1, message)
	}
	http.Error(myWriter, message, code)
}

func (myWriter *MyWriter) IssueCookie(value string, age int) {
	cookie := http.Cookie{
		Name:     myJwt.Cookie,
		Value:    value,
		Path:     "/",
		MaxAge:   age,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(myWriter, &cookie)
}

func (myWriter *MyWriter) ExpireCookie() {
	myWriter.IssueCookie("", -1)
}
