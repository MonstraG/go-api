package pages

import (
	"go-api/infrastructure/reqRes"
)

type PageData struct {
	PageTitle string

	StylesHash string

	ErrorMessage string

	Username    string
	UserInitial string

	IsAdmin bool
}

func NewPageData(request *reqRes.MyRequest, pageTitle string) PageData {
	return PageData{
		PageTitle:   pageTitle,
		Username:    request.User.Username,
		StylesHash:  StylesHash,
		UserInitial: getInitialFromUsername(request.User.Username),
		IsAdmin:     request.User.IsAdmin(),
	}
}

func getInitialFromUsername(username string) string {
	if username == "" {
		return ""
	}
	return string(username[0])
}
