package setup

import (
	"gorm.io/gorm"
	"net/http"
)

type MyRequest struct {
	http.Request

	Config   AppConfig
	Username string
	Db       *gorm.DB
}
