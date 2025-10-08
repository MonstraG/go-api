package player

import (
	"go-api/infrastructure/appConfig"

	"gorm.io/gorm"
)

type Controller struct {
	db           *gorm.DB
	explorerRoot string
}

func NewController(config appConfig.AppConfig, db *gorm.DB) *Controller {
	return &Controller{db: db, explorerRoot: config.ExplorerRoot}
}
