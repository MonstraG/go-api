package forgotPassword

import (
	"errors"
	"fmt"
	"go-api/infrastructure/crypto"
	"go-api/infrastructure/models"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/reqRes"
	"go-api/pages"
	"net/http"

	"gorm.io/gorm"
)

var forgotPasswordTemplate = pages.ParsePage("forgotPassword/forgotPassword.gohtml")
var resetPasswordFormTemplate = pages.ParsePage("forgotPassword/resetPasswordForm.gohtml")
var passwordResetSuccessfullyTemplate = pages.ParsePage("forgotPassword/passwordResetSuccessfully.gohtml")

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

func (controller *Controller) GetForgotPasswordForm(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	renderForgotPasswordPage(w, r, "")
}

func renderForgotPasswordPage(w reqRes.MyResponseWriter, r *reqRes.MyRequest, errorMessage string) {
	pageData := PageData{
		PageData:     pages.NewPageData(r, "Forgot password"),
		ErrorMessage: errorMessage,
	}

	if errorMessage != "" {
		myLog.Error.Logf("Forgot password start failed with errorMessage: %s", errorMessage)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.RenderTemplate(forgotPasswordTemplate, pageData)
}

func (controller *Controller) SubmitForgotPasswordForm(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
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

	//if !user.CanResetPassword {
	//	myLog.Info.Logf("User %s cannot reset password", username)
	//	renderForgotPasswordPage(w, r, "User cannot reset password")
	//	return
	//}

	renderResetPasswordForm(w, r, user.Username, "")
}

func renderResetPasswordForm(w reqRes.MyResponseWriter, r *reqRes.MyRequest, username string, errorMessage string) {
	pageData := ResetPasswordPageData{
		PageData:     pages.NewPageData(r, "Reset password"),
		Username:     username,
		ErrorMessage: errorMessage,
		MinLength:    minPasswordLength,
	}

	if errorMessage != "" {
		myLog.Error.Logf("Reset password failed with errorMessage: %s", errorMessage)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.RenderTemplate(resetPasswordFormTemplate, pageData)
}

func (controller *Controller) SetPassword(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
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

	user.PasswordSalt = crypto.NewSalt()
	user.PasswordHash = crypto.HashPassword(password, user.PasswordSalt)
	user.CanResetPassword = false
	controller.db.Save(&user)

	pageData := ResetPasswordPageData{
		PageData: pages.NewPageData(r, "Password reset"),
	}

	w.RenderTemplate(passwordResetSuccessfullyTemplate, pageData)
}
