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
		Username:    request.Token.Username,
		StylesHash:  StylesHash,
		UserInitial: getInitialFromUsername(request.Token.Username),
		IsAdmin:     request.Token.IsAdmin,
	}
}

func getInitialFromUsername(username string) string {
	if username == "" {
		return ""
	}
	return string(username[0])
}
