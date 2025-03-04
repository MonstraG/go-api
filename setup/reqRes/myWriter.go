package reqRes

import (
	"bufio"
	"errors"
	"go-server/setup/myJwt"
	"net"
	"net/http"
)

type MyWriter struct {
	http.ResponseWriter
}

func (w MyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func (w MyWriter) RedirectToLogin() {
	w.Header().Set("Location", `/login`)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (w MyWriter) IssueCookie(value string, age int) {
	cookie := http.Cookie{
		Name:     myJwt.Cookie,
		Value:    value,
		Path:     "/",
		MaxAge:   age,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
}

func (w MyWriter) ExpireCookie() {
	w.IssueCookie("", -1)
}
