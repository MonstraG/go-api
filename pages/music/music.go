package music

import (
	"fmt"
	"go-server/models"
	"go-server/pages"
	"go-server/pages/music/ytDlp"
	"go-server/setup/reqRes"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func PostHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		log.Printf("Failed to parse form:\n%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	songUrl := r.Form.Get("songUrl")
	if songUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteSilent([]byte("songUrl missing"))
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
		message := fmt.Sprintf("Failed to get song queue: \n%v", result.Error)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		err := songQueueEmptyTemplate.Execute(w, nil)
		if err != nil {
			// todo: use this pattern everywhere
			message := fmt.Sprintf("Failed to execute song queue template:\n%v\n", err)
			http.Error(w, message, http.StatusInternalServerError)
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
		// todo: use this pattern everywhere
		message := fmt.Sprintf("Failed to execute song queue template:\n%v\n", err)
		http.Error(w, message, http.StatusInternalServerError)
	}
}

var songPlayerTemplate = template.Must(template.ParseFiles("pages/music/songPlayer.gohtml"))
var songPlayerEmptyTemplate = template.Must(template.ParseFiles("pages/music/songPlayerEmpty.gohtml"))

func GetSongPlayerHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	var songQueueItem models.SongQueueItem
	result := getSongQueue(r).First(&songQueueItem)
	if result.Error != nil {
		message := fmt.Sprintf("Failed to get current song: \n%v", result.Error)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	if result.RowsAffected == 0 {
		err := songPlayerEmptyTemplate.Execute(w, nil)
		if err != nil {
			// todo: use this pattern everywhere
			message := fmt.Sprintf("Failed to execute song queue template:\n%v\n", err)
			http.Error(w, message, http.StatusInternalServerError)
		}
	}

	err := songPlayerTemplate.Execute(w, songQueueItem)
	if err != nil {
		// todo: use this pattern everywhere
		message := fmt.Sprintf("Failed to execute song queue template:\n%v\n", err)
		http.Error(w, message, http.StatusInternalServerError)
	}
}

func GetSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join(r.AppConfig.SongsFolder, pathQueryParam)
	pages.ServeFile(w, r, filename)
}
