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
var resetPasswordFormTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/forgotPassword/resetPasswordForm.gohtml"))
var passwordResetSuccessfullyTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/forgotPassword/passwordResetSuccessfully.gohtml"))

const minPasswordLength = 16

type PageData struct {
	pages.PageData
	ErrorMessage string
}

type ResetPasswordPageData struct {
	pages.PageData
	Username     string
	ErrorMessage string
	MinLength    int
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

	renderResetPasswordForm(w, r, user.Username, "")
}

func renderResetPasswordForm(w reqRes.MyWriter, r *reqRes.MyRequest, username string, errorMessage string) {
	pageData := ResetPasswordPageData{
		PageData:     pages.NewPageData(r, "Reset password"),
		Username:     username,
		ErrorMessage: errorMessage,
		MinLength:    minPasswordLength,
	}

	err := resetPasswordFormTemplate.Execute(w, pageData)
	if err != nil {
		message := fmt.Sprintf("Failed to render page: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
	}
}

func (controller Controller) PostSetPasswordHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
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

	password := r.Form.Get("password")
	if password == "" {
		http.Error(w, "password is required", http.StatusBadRequest)
		return
	}

	repeatPassword := r.Form.Get("repeatPassword")
	if repeatPassword == "" {
		http.Error(w, "repeatPassword is required", http.StatusBadRequest)
		return
	}

	if password != repeatPassword {
		renderResetPasswordForm(w, r, username, "Passwords do not match")
		return
	}

	if len(password) < minPasswordLength {
		renderResetPasswordForm(w, r, username, "Password must be at least 20 characters")
		return
	}

	var user *models.User
	result := controller.db.First(&user, "username = ?", username)
	if result.RowsAffected == 0 || errors.Is(result.Error, gorm.ErrRecordNotFound) {
		message := fmt.Sprintf("User %s not found", username)
		log.Println(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	if result.Error != nil {
		message := fmt.Sprintf("Failed to search for user: \n%v", result.Error)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	passwordHash, err := models.HashPassword(password)
	user.PasswordHash = passwordHash
	controller.db.Save(&user)

	pageData := ResetPasswordPageData{
		PageData: pages.NewPageData(r, "Password reset"),
	}
	err = passwordResetSuccessfullyTemplate.Execute(w, pageData)
	if err != nil {
		message := fmt.Sprintf("Failed to render page: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
	}
}
