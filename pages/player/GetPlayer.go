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
	"strconv"
	"time"
)

type Queue struct {
	CurrentSong models.QueuedSong
	Items       []models.QueuedSong
}

var playerTemplate = template.Must(template.ParseFiles("pages/player/playerPartial.gohtml"))

func (controller *Controller) GetPlayer(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	var queuedSongs []models.QueuedSong

	result := controller.db.Where("duration = 0 OR datetime(ends_at) > datetime()").Find(&queuedSongs)
	if result.Error != nil {
		message := fmt.Sprintf("Failed to render player: \n%v", result.Error)
		w.Error(message, http.StatusBadRequest)
		return
	}

	var currentSong models.QueuedSong
	if len(queuedSongs) > 0 {
		currentSong = queuedSongs[0]
	}

	pageData := Queue{
		CurrentSong: currentSong,
		Items:       queuedSongs,
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
		Path: pathQueryParam,
	})

	if result.Error != nil {
		message := fmt.Sprintf("Failed to insert song into queue: \n%v", result.Error)
		w.Error(message, http.StatusBadRequest)
		return
	}

	controller.GetPlayer(w, r)
}

// ReportSongDuration should be called by client to tell the server when the song actually ends
// It would have been nice to be able to figure out duration server-side, but that seems to not be that easy
func (controller *Controller) ReportSongDuration(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	queuedSongId := r.PathValue("queuedSongId")
	durationStr := r.URL.Query().Get("duration")

	var song models.QueuedSong
	result := controller.db.First(&song, queuedSongId)
	if result.RowsAffected == 0 {
		message := fmt.Sprintf("Failed to report song duration, song not found")
		w.Error(message, http.StatusBadRequest)
		return
	}

	if song.Duration == 0 {
		duration, err := strconv.ParseFloat(durationStr, 64)
		if err != nil {
			message := fmt.Sprintf("Failed to parse duration: %v", err.Error())
			w.Error(message, http.StatusBadRequest)
			return
		}

		song.Duration = time.Duration(duration * float64(time.Second))
		song.EndsAt = song.CreatedAt.Add(song.Duration)
		controller.db.Updates(&song)
	} else {
		// duration already known, ignore
	}

	w.Header().Set("HX-Trigger", "playerReloadEvent")
}
