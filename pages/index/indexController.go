package index

import (
	"go-api/pages"
	"go-api/setup/appConfig"
	"go-api/setup/reqRes"
	"html/template"
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
	var pageData = PageData{
		PageData:     pages.NewPageData(r, "Homepage"),
		VpsLoginLink: controller.vpsLoginLink,
	}

	w.RenderTemplate(indexTemplate, pageData)
}
