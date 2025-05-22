package reqRes

import (
	"net/http"
)

type MyRequest struct {
	http.Request
	Username string
}

func (myRequest *MyRequest) CookieIfValid(name string) (*http.Cookie, error) {
	cookie, err := myRequest.Cookie(name)
	if err != nil {
		return nil, err
	}
	err = cookie.Valid()
	if err != nil {
		return nil, err
	}
	return cookie, nil
}
