package music

import (
	"fmt"
	"go-api/infrastructure/helpers"
	"go-api/infrastructure/reqRes"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const hundredMegs = 100 << 20

func (controller *Controller) PutSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")

	err := r.ParseMultipartForm(hundredMegs)
	if err != nil {
		message := fmt.Sprintf("Failed to parse form: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	formFile, handler, err := r.FormFile("file")
	if err != nil {
		message := fmt.Sprintf("Failed to retrieve file: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}
	defer helpers.CloseSafely(formFile)

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	folder := filepath.Join(controller.songsFolder, pathQueryParam)
	path := filepath.Join(folder, handler.Filename)

	diskFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, filePermissions)
	if err != nil {
		message := fmt.Sprintf("Failed to retrieve file: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}
	defer helpers.CloseSafely(diskFile)

	_, err = io.Copy(diskFile, formFile)
	if err != nil {
		message := fmt.Sprintf("Failed to save file: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	readDir(w, folder, pathQueryParam, "File uploaded!")
}
