package notFound

import (
	"go-api/infrastructure/reqRes"
	"go-api/pages"
	"net/http"
)

var notFoundTemplate = pages.ParsePage("notFound/notFound.gohtml")

func Show404(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pageData := pages.NewPageData(r, "404: page not found")

	w.WriteHeader(http.StatusNotFound)
	w.RenderTemplate(notFoundTemplate, pageData)
}
