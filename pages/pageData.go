package pages

import "go-api/infrastructure/reqRes"

type PageData struct {
	PageTitle string

	ErrorMessage string

	Username    string
	UserInitial string
}

func NewPageData(request *reqRes.MyRequest, pageTitle string) PageData {
	return PageData{
		PageTitle:   pageTitle,
		Username:    request.Username,
		UserInitial: getInitialFromUsername(request.Username),
	}
}

func getInitialFromUsername(username string) string {
	if username == "" {
		return ""
	}
	return string(username[0])
}
