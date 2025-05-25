package music

import (
	"fmt"
	"go-api/setup/reqRes"
	"net/http"
	"os"
	"path/filepath"
)

const kilobyte = 1 << 10

func (controller *Controller) CreateFolderHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")

	err := r.ParseMultipartForm(kilobyte)
	if err != nil {
		message := fmt.Sprintf("Failed to parse form: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	value := r.FormValue("name")
	if value == "" {
		message := fmt.Sprintf("No name provided for folder")
		w.Error(message, http.StatusBadRequest)
		return
	}

	folder := filepath.Join(controller.songsFolder, pathQueryParam)
	path := filepath.Join(folder, value)

	err = os.Mkdir(path, 0666)
	if err != nil {
		message := fmt.Sprintf("Failed to create folder: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	readDir(w, folder, pathQueryParam, "Folder created!")
}
