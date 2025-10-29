package index

import (
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/reqRes"
	"go-api/pages"
)

var indexTemplate = pages.ParsePage(
	"nav.gohtml",
	"index/index.gohtml",
)

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

func (controller *Controller) GetHandler(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	var pageData = PageData{
		PageData:     pages.NewPageData(r, "Homepage"),
		VpsLoginLink: controller.vpsLoginLink,
	}

	w.RenderTemplate(indexTemplate, pageData)
}
