package music

import (
	"encoding/json"
	"fmt"
	"go-server/models"
	"go-server/pages/music/websockets"
	"go-server/pages/music/ytDlp"
	"go-server/setup/reqRes"
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

func PongHandler(w reqRes.MyWriter, _ *reqRes.MyRequest) {
	websockets.HubSingleton.Broadcast("pong")
	w.WriteHeader(http.StatusOK)
}

const startDelay = 10 * time.Second

func GetSongQueueHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	var songQueueItems []models.SongQueueItem
	result := r.Db.Joins("Song").Find(songQueueItems)
	if result.Error != nil {
		log.Printf("Failed to get song queue: \n%v\n", result.Error)
	}

	type SongQueueItemDTO struct {
		QueueItemId uint      `json:"queueItemId"`
		SongId      uint      `json:"songId"`
		StartsAt    time.Time `json:"startsAt"`
	}

	songCount := len(songQueueItems)
	songQueueItemDTOs := make([]SongQueueItemDTO, songCount)
	for index, songQueueItem := range songQueueItems {
		songQueueItemDTOs[index] = SongQueueItemDTO{
			QueueItemId: songQueueItem.ID,
			SongId:      songQueueItem.SongId,
			StartsAt:    songQueueItem.CreatedAt.Add(startDelay),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&songQueueItemDTOs)
	if err != nil {
		// todo: use this pattern everywhere
		message := fmt.Sprintf("Failed to encode song queue:\n%v\n", err)
		http.Error(w, message, http.StatusInternalServerError)
	}
}

func GetSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	//var song []models.Song
}
