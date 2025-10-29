package login

import (
	"database/sql"
	"errors"
	"fmt"
	"go-api/infrastructure/models"
	"go-api/infrastructure/myJwt"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/reqRes"
	"go-api/pages"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

var loginTemplate = pages.ParsePage("login/login.gohtml")

type PageData struct {
	pages.PageData
	ErrorMessage string
}

type Controller struct {
	MyJwt *myJwt.Service
	Db    *gorm.DB
}

func NewController(myJwt *myJwt.Service, Db *gorm.DB) *Controller {
	return &Controller{MyJwt: myJwt, Db: Db}
}

func (controller *Controller) GetHandler(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	renderLoginPage(w, r, "")
}

func renderLoginPage(w reqRes.MyResponseWriter, r *reqRes.MyRequest, errorMessage string) {
	pageData := PageData{
		PageData:     pages.NewPageData(r, "Login"),
		ErrorMessage: errorMessage,
	}

	if errorMessage != "" {
		myLog.Error.Logf("Login failed with errorMessage: %s", errorMessage)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.RenderTemplate(loginTemplate, pageData)
}

func (controller *Controller) PostHandler(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
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

	lowercaseUsername := strings.ToLower(username)
	user := models.User{}
	result := controller.Db.Where("lower(username) = @name", sql.Named("name", lowercaseUsername)).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		renderLoginPage(w, r, "Username or password is invalid")
		return
	}
	if result.Error != nil {
		message := fmt.Sprintf("Hit error when searching for user '%v': \n%v", lowercaseUsername, result.Error)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		fmt.Printf("Rows affected 0 when searching for user '%v'", lowercaseUsername)
		renderLoginPage(w, r, "Username or password is invalid")
		return
	}

	ok = user.CheckPasswordHash(password)
	if !ok {
		renderLoginPage(w, r, "Username or password is invalid")
		return
	}

	jwtToken, err := controller.MyJwt.CreateJwt(user)
	if err != nil {
		message := fmt.Sprintf("Error generating jwt token for user '%v': \n%v", lowercaseUsername, result.Error)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	w.IssueCookie(jwtToken, myJwt.MaxAge)

	http.Redirect(w, &r.Request, "/", http.StatusSeeOther)
}
