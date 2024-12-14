package ytDlp

import (
	"errors"
	"fmt"
	"go-server/models"
	"go-server/pages/music/websockets"
	"go-server/setup/appConfig"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func SanitizeUrl(url string) string {
	index := strings.Index(url, "&")
	if index != -1 {
		return url[:index]
	}
	return url
}

func getSongId(url string) (string, error) {
	index := strings.Index(url, "v=")
	if index == -1 {
		return "", errors.New(fmt.Sprintf("didn't find `v=` in url `%s`", url))
	}
	return url[index+2:], nil
}

const ytDlpBinary = "yt-dlp"
const fileNamePattern = "%(id)s-%(duration)s.%(ext)s"

func GetDuration(id string, config appConfig.AppConfig) int {
	files, err := os.ReadDir(config.SongsFolder)
	for _, file := range files {
		fileNameWithExtension := file.Name()
		if strings.HasPrefix(fileNameWithExtension, id) {
			extension := filepath.Ext(fileNameWithExtension)
			fileName := strings.TrimSuffix(fileNameWithExtension, extension)
			durationStr := fileName[len(id)+1:] // also skip separatory `-`
			duration, err := strconv.Atoi(durationStr)
			if err != nil {
				log.Printf("failed to read song duration: %s, error: \n%v\n", fileNameWithExtension, err)
				return 0
			}
			return duration
		}
	}
	if err != nil {
		log.Printf("failed to read songs folder: %s, error: \n%v\n", config.SongsFolder, err)
	}
	return 0
}

func Download(url string, config appConfig.AppConfig, db *gorm.DB) {
	log.Println("Starting song download")
	id, err := getSongId(url)
	if err != nil {
		log.Printf("Failed to get song id:\n%v\n", err)
		return
	}
	log.Printf("Song url: %s\n", url)

	outputFile := filepath.Join(config.SongsFolder, fileNamePattern)
	log.Printf("Destination: %s\n", outputFile)

	command := exec.Command(ytDlpBinary, []string{url, "-x", "-o", outputFile}...)
	command.Stdout = os.Stdout
	err = command.Start()
	if err != nil {
		log.Printf("Failed to start download command:\n%v\n", err)
	}

	go func() {
		err = command.Wait()
		if err != nil {
			log.Printf("Error from wait on command:\n%v\n", err)
			return
		}

		duration := GetDuration(id, config)
		if duration == 0 {
			return
		}

		song := &models.Song{YoutubeId: id, Duration: duration}
		result := db.FirstOrCreate(song)
		if result.Error != nil {
			log.Printf("Failed to save song:\n%v\n", result.Error)
		}
		if result.RowsAffected == 0 {
			log.Printf("Failed to save song: 0 rows affected\n")
		}

		songQueueItem := &models.SongQueueItem{
			SongId: song.ID,
		}
		result = db.Create(songQueueItem)
		if result.Error != nil {
			log.Printf("Failed to save song queue item:\n%v\n", result.Error)
			return
		}
		if result.RowsAffected == 0 {
			log.Printf("Failed to save song: 0 rows affected\n")
			return
		}

		log.Printf("Finished downloading song \n")
		websockets.HubSingleton.Broadcast("song")
	}()
}
