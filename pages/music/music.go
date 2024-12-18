package music

import (
	"fmt"
	"go-server/models"
	"go-server/pages/music/ytDlp"
	"go-server/setup/reqRes"
	"html/template"
	"log"
	"net/http"
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

var songQueueTemplate = template.Must(template.ParseFiles("pages/music/songQueue.gohtml"))
var songQueueEmptyTemplate = template.Must(template.ParseFiles("pages/music/songQueueEmpty.gohtml"))

func GetSongQueueHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	now := time.Now()

	var songQueueItems []models.SongQueueItem
	result := r.Db.Joins("Song").Where("ends_at > ?", now).Find(&songQueueItems)
	if result.Error != nil {
		log.Printf("Failed to get song queue: \n%v\n", result.Error)
	}

	if result.RowsAffected == 0 {
		err := songQueueEmptyTemplate.Execute(w, songQueueItems)
		if err != nil {
			// todo: use this pattern everywhere
			message := fmt.Sprintf("Failed to execute song queue template:\n%v\n", err)
			http.Error(w, message, http.StatusInternalServerError)
		}
		return
	}

	type SongQueueItemDTO struct {
		Song     string
		Duration time.Duration
		StartsIn time.Duration
		EndsIn   time.Duration
	}

	songCount := len(songQueueItems)
	songs := make([]SongQueueItemDTO, songCount)
	for index, songQueueItem := range songQueueItems {
		songItem := SongQueueItemDTO{
			Song:     songQueueItem.Song.YoutubeId,
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

	log.Printf("SongQueue: \n%v\n", songs)

	err := songQueueTemplate.Execute(w, List{songs})
	if err != nil {
		// todo: use this pattern everywhere
		message := fmt.Sprintf("Failed to execute song queue template:\n%v\n", err)
		http.Error(w, message, http.StatusInternalServerError)
	}
}
