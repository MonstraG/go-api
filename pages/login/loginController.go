package login

import (
	"database/sql"
	"errors"
	"fmt"
	"go-server/models"
	"go-server/pages"
	"go-server/setup/myJwt"
	"go-server/setup/reqRes"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var loginTemplate = template.Must(template.ParseFiles("pages/base.gohtml", "pages/login/login.gohtml"))

type PageData struct {
	pages.PageData
	ErrorMessage string
}

func GetHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	renderLoginPage(w, r, "")
}

func renderLoginPage(w reqRes.MyWriter, r *reqRes.MyRequest, errorMessage string) {
	loginPageData := PageData{
		PageData:     pages.NewPageData(r, "Login"),
		ErrorMessage: errorMessage,
	}

	err := loginTemplate.Execute(w, loginPageData)
	if err != nil {
		message := fmt.Sprintf("Failed to render login page: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
	}
}

func PostHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
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
		http.Error(w, "password", http.StatusBadRequest)
		return
	}

	lowercaseUsername := strings.ToLower(username)
	user := models.User{}
	result := r.Db.Where("lower(username) = @name", sql.Named("name", lowercaseUsername)).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		renderLoginPage(w, r, "Username or password is invalid")
		return
	}
	if result.Error != nil {
		message := fmt.Sprintf("Hit error when searching for user '%v': \n%v", lowercaseUsername, result.Error)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		fmt.Printf("Rows affected 0 when searching for user '%v'", lowercaseUsername)
		renderLoginPage(w, r, "Username or password is invalid")
		return
	}

	ok := user.CheckPasswordHash(password)
	if !ok {
		renderLoginPage(w, r, "Username or password is invalid")
		return
	}

	jwtToken, err := myJwt.Singleton.CreateJwt(user, r.AppConfig)
	if err != nil {
		message := fmt.Sprintf("Error generating jwt token for user '%v': \n%v", lowercaseUsername, result.Error)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:     "jwtToken",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   myJwt.MaxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, &r.Request, "/", http.StatusSeeOther)
}
