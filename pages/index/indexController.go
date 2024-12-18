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
var indexPageData = pages.PageData{
	PageTitle: "Homepage",
}

func GetHandler(w reqRes.MyWriter, _ *reqRes.MyRequest) {
	err := indexTemplate.Execute(w, indexPageData)
	if err != nil {
		log.Printf("Failed to render index page:\n%v\n", err)
	}
}
