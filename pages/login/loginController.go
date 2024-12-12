package login

import (
	"database/sql"
	"errors"
	"go-server/models"
	"go-server/pages"
	"go-server/setup/reqRes"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetHandler(w reqRes.MyWriter, _ *reqRes.MyRequest) {
	var indexTemplate = template.Must(template.ParseFiles("pages/login/login.gohtml"))
	var indexPageData = pages.PageData{
		PageTitle: "Homepage",
	}

	err := indexTemplate.Execute(w, indexPageData)
	if err != nil {
		log.Fatal("Failed to render login page:\n", err)
	}
}

func PostHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed to parse form:\n%v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	username := r.Form.Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	password := r.Form.Get("password")
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lowercaseUsername := strings.ToLower(username)
	user := models.User{}
	result := r.Db.Where("lower(username) = @name", sql.Named("name", lowercaseUsername)).First(&user)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("Hit error when searching for user '%v':\n%v\n", lowercaseUsername, result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		// didn't find them
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	ok := user.CheckPasswordHash(password)
	if !ok {
		// wrong password
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	jwtToken, err := createJwt(user, r.AppConfig, now)
	if err != nil {
		log.Printf("Error generating jwt token for user '%v':\n%v\n", lowercaseUsername, result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = jwtToken
	w.WriteHeader(http.StatusOK)
}

func now() time.Time {
	return time.Now()
}
