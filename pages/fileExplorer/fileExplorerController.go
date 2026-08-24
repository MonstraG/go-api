package fileExplorer

import "go-api/infrastructure/appConfig"

type Controller struct {
	explorerRoot string
}

func NewController(config appConfig.AppConfig) *Controller {
	return &Controller{explorerRoot: config.ExplorerRoot}
}
