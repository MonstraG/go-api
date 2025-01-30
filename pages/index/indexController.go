package index

import (
	"fmt"
	"go-server/pages"
	"go-server/setup/reqRes"
	"html/template"
	"net/http"
)

var indexTemplate = template.Must(template.ParseFiles(
	"pages/base.gohtml",
	"pages/nav.gohtml",
	"pages/index/index.gohtml",
))

type PageData struct {
	pages.PageData
	VpsLoginLink string
}

func GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	var indexPageData = PageData{
		PageData:     pages.PageData{PageTitle: "Homepage"},
		VpsLoginLink: r.AppConfig.VpsLoginLink,
	}

	err := indexTemplate.Execute(w, indexPageData)
	if err != nil {
		w.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to render index page:\n%v", err))
	}
}
