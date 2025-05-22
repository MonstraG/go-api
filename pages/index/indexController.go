package index

import (
	"fmt"
	"go-server/pages"
	"go-server/setup/appConfig"
	"go-server/setup/reqRes"
	"html/template"
	"log"
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

type Controller struct {
	vpsLoginLink string
}

func NewController(config appConfig.AppConfig) *Controller {
	return &Controller{
		vpsLoginLink: config.VpsLoginLink,
	}
}

func (controller *Controller) GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	var indexPageData = PageData{
		PageData:     pages.NewPageData(r, "Homepage"),
		VpsLoginLink: controller.vpsLoginLink,
	}

	err := indexTemplate.Execute(w, indexPageData)
	if err != nil {
		message := fmt.Sprintf("Failed to render index page:\n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
	}
}
