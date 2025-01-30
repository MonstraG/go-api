package notFound

import (
	"fmt"
	"go-server/pages"
	"go-server/setup/reqRes"
	"html/template"
	"net/http"
)

var notFoundTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/notFound/notFound.gohtml"))
var notFoundPageData = pages.PageData{
	PageTitle: "404: page not found",
}

func GetHandler(w reqRes.MyWriter, _ *reqRes.MyRequest) {
	err := notFoundTemplate.Execute(w, notFoundPageData)
	if err != nil {
		w.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to render 404 page: \n%v", err))
	}
}
