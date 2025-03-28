package music

import (
	"errors"
	"fmt"
	"go-server/setup/reqRes"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func DeleteSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	pathToRemove := filepath.Join(r.AppConfig.SongsFolder, pathQueryParam)

	stat, err := os.Stat(pathToRemove)
	if errors.Is(err, os.ErrNotExist) {
		message := fmt.Sprintf("Failed to remove: \n%v", err)
		log.Println(message)
		http.Error(w, "Path not found", http.StatusBadRequest)
	}

	isDir := stat.IsDir()
	if isDir {
		err = os.RemoveAll(pathToRemove)
	} else {
		err = os.Remove(pathToRemove)
	}
	if err != nil {
		message := fmt.Sprintf("Failed to remove: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
	}

	fileSystemFolder := filepath.Dir(pathToRemove)
	queryFolder := filepath.Dir(pathQueryParam)
	if queryFolder == "." {
		queryFolder = ""
	}

	var resultMessage string
	if isDir {
		resultMessage = "Folder deleted!"
	} else {
		resultMessage = "File deleted!"
	}

	readDir(w, fileSystemFolder, queryFolder, resultMessage)
}
