package helpers

import "net/http"

type MyRequest struct {
	Username string
	http.Request
}
