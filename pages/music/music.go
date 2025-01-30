package music

import (
	"fmt"
	"go-server/models"
	"go-server/pages"
	"go-server/pages/music/ytDlp"
	"go-server/setup/reqRes"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

func PostHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		w.Error(http.StatusBadRequest, fmt.Sprintf("Failed to parse form: \n%v", err))
		return
	}

	songUrl := r.Form.Get("songUrl")
	if songUrl == "" {
		w.Error(http.StatusBadRequest, "songUrl missing")
		return
	}

	sanitizedUrl := ytDlp.SanitizeUrl(songUrl)

	ytDlp.Download(sanitizedUrl, r.AppConfig, r.Db)
}

func getSongQueue(r *reqRes.MyRequest) *gorm.DB {
	now := time.Now()
	return r.Db.Where("ends_at > ?", now).Order("song_queue_items.starts_at asc")
}

var songQueueTemplate = template.Must(template.ParseFiles("pages/music/songQueue.gohtml"))
var songQueueEmptyTemplate = template.Must(template.ParseFiles("pages/music/songQueueEmpty.gohtml"))

func GetSongQueueHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	now := time.Now()

	var songQueueItems []models.SongQueueItem
	result := getSongQueue(r).Find(&songQueueItems)
	if result.Error != nil {
		w.Error(http.StatusBadRequest, fmt.Sprintf("Failed to get song queue: \n%v", result.Error))
		return
	}

	if result.RowsAffected == 0 {
		err := songQueueEmptyTemplate.Execute(w, nil)
		if err != nil {
			w.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to execute song queue template: \n%v", err))
		}
		return
	}

	type SongQueueItemDTO struct {
		ID       uint
		Song     string
		Duration time.Duration
		StartsIn time.Duration
		EndsIn   time.Duration
	}

	songCount := len(songQueueItems)
	songs := make([]SongQueueItemDTO, songCount)
	for index, songQueueItem := range songQueueItems {
		songItem := SongQueueItemDTO{
			ID:       songQueueItem.SongId,
			Song:     songQueueItem.Song.Title,
			Duration: (time.Duration(songQueueItem.Song.Duration) * time.Second).Truncate(time.Second),
		}

		songItem.StartsIn = songQueueItem.StartsAt.Sub(now).Truncate(time.Second)
		if songItem.StartsIn < 0 {
			songItem.StartsIn = 0
		}
		songItem.EndsIn = songQueueItem.EndsAt.Sub(now).Truncate(time.Second)

		songs[index] = songItem
	}

	type List struct {
		Songs []SongQueueItemDTO
	}

	err := songQueueTemplate.Execute(w, List{songs})
	if err != nil {
		w.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to execute song queue template: \n%v", err))
	}
}

var songPlayerTemplate = template.Must(template.ParseFiles("pages/music/songPlayer.gohtml"))
var songPlayerEmptyTemplate = template.Must(template.ParseFiles("pages/music/songPlayerEmpty.gohtml"))

func GetSongPlayerHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	var songQueueItem models.SongQueueItem
	result := getSongQueue(r).First(&songQueueItem)
	if result.Error != nil {
		w.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to get current song: \n%v", result.Error))
		return
	}
	if result.RowsAffected == 0 {
		err := songPlayerEmptyTemplate.Execute(w, nil)
		if err != nil {
			w.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to execute song queue template: \n%v", err))
		}
	}

	err := songPlayerTemplate.Execute(w, songQueueItem)
	if err != nil {
		w.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to execute song queue template: \n%v", err))
	}
}

func GetSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join(r.AppConfig.SongsFolder, pathQueryParam)
	pages.ServeFile(w, r, filename)
}
