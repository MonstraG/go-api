package reqRes

import (
	"bufio"
	"errors"
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
