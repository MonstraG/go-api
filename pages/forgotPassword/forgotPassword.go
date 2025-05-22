package forgotPassword

import (
	"errors"
	"fmt"
	"go-server/models"
	"go-server/pages"
	"go-server/setup/reqRes"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
)

var forgotPasswordTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/forgotPassword/forgotPassword.gohtml"))

type PageData struct {
	pages.PageData
	ErrorMessage string
}

type Controller struct {
	db *gorm.DB
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{db: db}
}

func (controller Controller) GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	renderForgotPasswordPage(w, r, "")
}

func renderForgotPasswordPage(w reqRes.MyWriter, r *reqRes.MyRequest, errorMessage string) {
	pageData := PageData{
		PageData:     pages.NewPageData(r, "Forgot password"),
		ErrorMessage: errorMessage,
	}

	err := forgotPasswordTemplate.Execute(w, pageData)
	if err != nil {
		message := fmt.Sprintf("Failed to render page: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
	}
}

func (controller Controller) PostHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	err := r.ParseForm()
	if err != nil {
		message := fmt.Sprintf("Failed to parse form: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	var user *models.User
	result := controller.db.First(&user, "username = ?", username)
	if result.RowsAffected == 0 || errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println(fmt.Sprintf("User %s not found", username))
		renderForgotPasswordPage(w, r, "User not found")
		return
	}
	if result.Error != nil {
		message := fmt.Sprintf("Failed to search for user: \n%v", result.Error)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	if !user.CanResetPassword {
		log.Println(fmt.Sprintf("User %s cannot reset password", username))
		renderForgotPasswordPage(w, r, "User cannot reset password")
		return
	}

	w.Header().Set("Location", `/`)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
