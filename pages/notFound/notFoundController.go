package notFound

import (
	"go-api/infrastructure/reqRes"
	"go-api/pages"
	"html/template"
)

var notFoundTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/notFound/notFound.gohtml"))

func GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pageData := pages.NewPageData(r, "404: page not found")

	w.RenderTemplate(notFoundTemplate, pageData)
}
