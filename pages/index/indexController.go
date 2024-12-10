package index

import (
	"go-server/pages"
	"go-server/setup"
	"html/template"
	"log"
)

var indexTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/index/index.gohtml"))
var indexPageData = pages.PageData{
	PageTitle: "Homepage",
}

func GetHandler(w setup.MyWriter, _ *setup.MyRequest) {
	err := indexTemplate.Execute(w, indexPageData)
	if err != nil {
		log.Fatal("Failed to render index page:\n", err)
	}
}
