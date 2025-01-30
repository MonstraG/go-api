package notFound

import (
	"fmt"
	"go-server/pages"
	"go-server/setup/reqRes"
	"html/template"
	"log"
	"net/http"
)

var notFoundTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/notFound/notFound.gohtml"))

func GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	var notFoundPageData = pages.NewPageData(r, "404: page not found")
	err := notFoundTemplate.Execute(w, notFoundPageData)
	if err != nil {
		message := fmt.Sprintf("Failed to render 404 page: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
	}
}
