package reqRes

import (
	"go-server/setup/appConfig"
	"gorm.io/gorm"
	"net/http"
)

type MyRequest struct {
	http.Request

	AppConfig appConfig.AppConfig
	Username  string
	Db        *gorm.DB
}
