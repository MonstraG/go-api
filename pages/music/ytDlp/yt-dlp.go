package ytDlp

import (
	"errors"
	"fmt"
	"go-server/models"
	"go-server/setup/appConfig"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

func findSongFilename(id string, config appConfig.AppConfig) string {
	files, err := os.ReadDir(config.SongsFolder)
	for _, file := range files {
		fileNameWithExtension := file.Name()
		if strings.HasPrefix(fileNameWithExtension, id) {
			extension := filepath.Ext(fileNameWithExtension)
			fileName := strings.TrimSuffix(fileNameWithExtension, extension)
			return fileName
		}
	}
	if err != nil {
		log.Printf("failed to read songs folder: %s, error: \n%v\n", config.SongsFolder, err)
	} else {
		log.Printf("failed to find song: %s in songs folder %s\n", id, config.SongsFolder)
	}
	return ""
}

func getSongModel(id string, config appConfig.AppConfig) *models.Song {
	filename := findSongFilename(id, config)
	if filename == "" {
		return nil
	}

	filenameParts := strings.Split(filename, "|")
	if len(filenameParts) != 3 {
		log.Printf("failed to split song name into parts: %s, got parts length = %v\n", filename, len(filenameParts))
		return nil
	}

	// 0 is id and we already know it
	title := filenameParts[1]
	duration, err := strconv.Atoi(filenameParts[2])
	if err != nil {
		log.Printf("failed to read song duration: %s, error: \n%v\n", filename, err)
		return nil
	}

	return &models.Song{
		YoutubeId: id,
		Duration:  duration,
		Title:     title,
	}
}

const ytDlpBinary = "yt-dlp"
const fileNamePattern = "%(id)s|%(title)s|%(duration)s.%(ext)s"

const startDelay = 10 * time.Second

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

	command := exec.Command(ytDlpBinary, []string{url, "-x", "--embed-metadata", "--embed-thumbnail", "-o", outputFile}...)
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

		log.Printf("Command finished successfully")

		song := getSongModel(id, config)
		if song == nil {
			log.Printf("Failed to get song model")
			return
		} else {
			log.Printf("Loaded song: %v\n", song)
		}

		result := db.FirstOrCreate(song)
		if result.Error != nil {
			log.Printf("Failed to save song:\n%v\n", result.Error)
		}
		if result.RowsAffected == 0 {
			log.Printf("Failed to save song: 0 rows affected\n")
		}

		startsAt := time.Now().Add(startDelay)
		songQueueItem := &models.SongQueueItem{
			SongId:   song.ID,
			StartsAt: startsAt,
			EndsAt:   startsAt.Add(time.Duration(song.Duration) * time.Second),
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
	}()
}
