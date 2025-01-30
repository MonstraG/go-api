package reqRes

import (
	"bufio"
	"errors"
	"log"
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

// Error is a wrapper around http.Error that also logs the message
func (w MyWriter) Error(statusCode int, message string) {
	log.Println(message)
	http.Error(w, message, statusCode)
}
