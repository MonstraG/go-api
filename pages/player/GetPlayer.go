package player

import (
	"fmt"
	"go-api/infrastructure/models"
	"go-api/infrastructure/reqRes"
	"go-api/pages/fileExplorer"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

type Queue struct {
	Items []models.QueuedSong
}

var playerTemplate = template.Must(template.ParseFiles("pages/player/playerPartial.gohtml"))

func (controller *Controller) GetPlayer(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	var queuedSongs []models.QueuedSong

	result := controller.db.Find(&queuedSongs)
	if result.Error != nil {
		message := fmt.Sprintf("Failed to render player: \n%v", result.Error)
		w.Error(message, http.StatusBadRequest)
		return
	}

	pageData := Queue{
		Items: queuedSongs,
	}

	w.RenderTemplate(playerTemplate, pageData)
}

func (controller *Controller) AddSong(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	pathToFile := filepath.Join(controller.explorerRoot, pathQueryParam)

	stat, err := os.Stat(pathToFile)
	if err != nil {
		message := fmt.Sprintf("Failed to open song to add to queue: \n%v", err)
		w.Error(message, http.StatusBadRequest)
		return
	}

	isDir := stat.IsDir()
	if isDir {
		w.Error("Failed to open song to add to queue, it's a directory", http.StatusBadRequest)
		return
	}

	isSong := fileExplorer.IsSong(stat.Name())
	if !isSong {
		w.Error("Failed to open song to add to queue, it's not a song", http.StatusBadRequest)
		return
	}

	result := controller.db.Create(&models.QueuedSong{
		Path: pathToFile,
	})

	if result.Error != nil {
		message := fmt.Sprintf("Failed to insert song into queue: \n%v", result.Error)
		w.Error(message, http.StatusBadRequest)
		return
	}

	controller.GetPlayer(w, r)
}
