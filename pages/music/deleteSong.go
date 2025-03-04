package music

import (
	"fmt"
	"go-server/setup/reqRes"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func DeleteSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	file := filepath.Join(r.AppConfig.SongsFolder, pathQueryParam)
	fileSystemFolder := filepath.Dir(file)

	queryFolder := filepath.Dir(pathQueryParam)
	if queryFolder == "." {
		queryFolder = ""
	}

	err := os.Remove(file)
	if err != nil {
		message := fmt.Sprintf("Failed to remove file: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
	}

	readDir(w, fileSystemFolder, queryFolder, "File deleted!")
}
