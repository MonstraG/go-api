package admin

import (
	"errors"
	"fmt"
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/models"
	"go-api/infrastructure/reqRes"
	"go-api/pages"
	"net/http"

	"gorm.io/gorm"
)

var indexTemplate = pages.ParsePage(
	"nav.gohtml",
	"admin/adminPage.gohtml",
)

type Controller struct {
	db           *gorm.DB
	vpsLoginLink string
}

type PageData struct {
	pages.PageData
	VpsLoginLink string
}

func NewController(config appConfig.AppConfig, db *gorm.DB) *Controller {
	return &Controller{db: db, vpsLoginLink: config.VpsLoginLink}
}

func (controller *Controller) GetAdminPage(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	var pageData = PageData{
		PageData:     pages.NewPageData(r, "Homepage"),
		VpsLoginLink: controller.vpsLoginLink,
	}

	w.RenderTemplate(indexTemplate, pageData)
}

func (controller *Controller) SetPasswordChangeStatus(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	ok := r.ParseFormRequired(w)
	if !ok {
		return
	}
	username := r.GetFormFieldRequired(w, "username")
	if username == "" {
		return
	}
	canChangePassword := r.GetFormFieldRequired(w, "canChangePassword")
	if canChangePassword == "" {
		return
	}

	user := models.User{Username: username}
	result := controller.db.First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		message := fmt.Sprintf("Failed to find user %v", result.Error)
		w.Error(message, http.StatusBadRequest)
		return
	}
	if result.Error != nil {
		message := fmt.Sprintf("Failed to find user %v", result.Error)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	user.CanResetPassword = canChangePassword == "true"
	result = controller.db.Save(&user)
	if result.Error != nil {
		message := fmt.Sprintf("Failed to update user %v", result.Error)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
