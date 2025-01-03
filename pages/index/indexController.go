package index

import (
	"go-server/pages"
	"go-server/setup/reqRes"
	"html/template"
	"log"
)

var indexTemplate = template.Must(template.ParseFiles(
	"pages/base.gohtml",
	"pages/nav.gohtml",
	"pages/index/index.gohtml",
))

func GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	var indexPageData = pages.PageData{
		PageTitle:    "Homepage",
		VpsLoginLink: r.AppConfig.VpsLoginLink,
	}

	err := indexTemplate.Execute(w, indexPageData)
	if err != nil {
		log.Printf("Failed to render index page:\n%v\n", err)
	}
}
