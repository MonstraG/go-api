package music

import (
	"go-server/models"
	"go-server/pages/music/websockets"
	"go-server/pages/music/ytDlp"
	"go-server/setup/reqRes"
	"log"
	"net/http"
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

func PongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	websockets.HubSingleton.Broadcast("pong")
	w.WriteHeader(http.StatusOK)
}

func GetSongsHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	var users []models.SongQueueItem
	result := r.Db.Joins("Song").Find(users)
	if result.Error != nil {
		log.Printf("Failed to get song queue: \n%v\n", result.Error)
	}

	type SongQueueItemDTO struct {
		SongId   int `json:"songId"`
		Duration int `json:"duration"`
	}

}

func GetSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	//var song []models.Song

}
