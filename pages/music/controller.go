package music

import "go-api/infrastructure/appConfig"

type Controller struct {
	songsFolder string
}

func NewController(config appConfig.AppConfig) *Controller {
	return &Controller{songsFolder: config.SongsFolder}
}
