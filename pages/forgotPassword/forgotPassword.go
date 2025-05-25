package forgotPassword

import (
	"errors"
	"fmt"
	"go-api/models"
	"go-api/pages"
	"go-api/setup/myLog"
	"go-api/setup/reqRes"
	"gorm.io/gorm"
	"html/template"
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

func (controller *Controller) GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	renderForgotPasswordPage(w, r, "")
}

func renderForgotPasswordPage(w reqRes.MyWriter, r *reqRes.MyRequest, errorMessage string) {
	pageData := PageData{
		PageData:     pages.NewPageData(r, "Forgot password"),
		ErrorMessage: errorMessage,
	}

	if errorMessage == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.RenderTemplate(forgotPasswordTemplate, pageData)
}

func (controller *Controller) PostHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	ok := r.ParseFormRequired(w)
	if !ok {
		return
	}

	username := r.GetFormFieldRequired(w, "username")
	if username == "" {
		return
	}

	var user *models.User
	result := controller.db.First(&user, "username = ?", username)
	if result.RowsAffected == 0 || errors.Is(result.Error, gorm.ErrRecordNotFound) {
		myLog.Info.Logf("User %s not found", username)
		renderForgotPasswordPage(w, r, "User not found")
		return
	}
	if result.Error != nil {
		message := fmt.Sprintf("Failed to search for user: \n%v", result.Error)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	if !user.CanResetPassword {
		myLog.Info.Logf("User %s cannot reset password", username)
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

	if errorMessage != "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.RenderTemplate(resetPasswordFormTemplate, pageData)
}

func (controller *Controller) PostSetPasswordHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	ok := r.ParseFormRequired(w)
	if !ok {
		return
	}

	username := r.GetFormFieldRequired(w, "username")
	if username == "" {
		return
	}
	password := r.GetFormFieldRequired(w, "password")
	if password == "" {
		return
	}
	repeatPassword := r.GetFormFieldRequired(w, "repeatPassword")
	if repeatPassword == "" {
		return
	}

	if password != repeatPassword {
		renderResetPasswordForm(w, r, username, "Passwords do not match")
		return
	}

	if len(password) < minPasswordLength {
		errorMessage := fmt.Sprintf("Password must be at least %d characters", minPasswordLength)
		renderResetPasswordForm(w, r, username, errorMessage)
		return
	}

	var user *models.User
	result := controller.db.First(&user, "username = ?", username)
	if result.RowsAffected == 0 || errors.Is(result.Error, gorm.ErrRecordNotFound) {
		message := fmt.Sprintf("User %s not found", username)
		w.Error(message, http.StatusBadRequest)
		return
	}
	if result.Error != nil {
		message := fmt.Sprintf("Failed to search for user: \n%v", result.Error)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	passwordHash, err := models.HashPassword(password)
	if err != nil {
		message := fmt.Sprintf("Failed to hash password: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}
	user.PasswordHash = passwordHash
	user.CanResetPassword = false
	controller.db.Save(&user)

	pageData := ResetPasswordPageData{
		PageData: pages.NewPageData(r, "Password reset"),
	}

	w.RenderTemplate(passwordResetSuccessfullyTemplate, pageData)
}
