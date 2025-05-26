package music

import (
	"fmt"
	"go-api/infrastructure/reqRes"
	"net/http"
	"os"
	"path/filepath"
)

const kilobyte = 1 << 10

const filePermissions = 0666

func (controller *Controller) CreateFolderHandler(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")

	err := r.ParseMultipartForm(kilobyte)
	if err != nil {
		message := fmt.Sprintf("Failed to parse form: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	value := r.FormValue("name")
	if value == "" {
		w.Error("No name provided for folder", http.StatusBadRequest)
		return
	}

	folder := filepath.Join(controller.songsFolder, pathQueryParam)
	path := filepath.Join(folder, value)

	err = os.Mkdir(path, filePermissions)
	if err != nil {
		message := fmt.Sprintf("Failed to create folder: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	readDir(w, folder, pathQueryParam, "Folder created!")
}
