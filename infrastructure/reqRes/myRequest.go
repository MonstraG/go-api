package reqRes

import (
	"fmt"
	"net/http"
)

type MyRequest struct {
	http.Request
	RequestId string
	UserId    string
	Username  string
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

func (myRequest *MyRequest) GetFormFieldRequired(w MyWriter, key string) string {
	value := myRequest.Form.Get(key)
	if value == "" {
		message := fmt.Sprintf("%s is required", key)
		w.Error(message, http.StatusBadRequest)
		return ""
	}
	return value
}

func (myRequest *MyRequest) ParseFormRequired(w MyWriter) bool {
	err := myRequest.ParseForm()
	if err != nil {
		message := fmt.Sprintf("Failed to parse form: \n%v", err)
		w.Error(message, http.StatusBadRequest)
		return false
	}
	return true
}
