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

// WriteSilent calls w.Write without telling you the result
func (w MyWriter) WriteSilent(content []byte) {
	_, err := w.ResponseWriter.Write(content)
	if err != nil {
		log.Printf("Write failed:\n%v\n", err)
	}
}

func (w MyWriter) WriteResponse(status int, content []byte) {
	w.WriteHeader(status)
	w.WriteSilent(content)
}

func (w MyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}
