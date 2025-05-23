package music

import (
	"fmt"
	"go-server/setup/myLog"
	"go-server/setup/reqRes"
	"io"
	"mime/multipart"
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

	file, handler, err := r.FormFile("file")
	if err != nil {
		message := fmt.Sprintf("Failed to retrieve file: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			myLog.Info.Logf("Failed to close file: \n%v", err)
		}
	}(file)

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	folder := filepath.Join(controller.songsFolder, pathQueryParam)
	path := filepath.Join(folder, handler.Filename)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		message := fmt.Sprintf("Failed to retrieve file: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			myLog.Info.Logf("Failed to close file: \n%v", err)
		}
	}(f)

	_, err = io.Copy(f, file)
	if err != nil {
		message := fmt.Sprintf("Failed to save file: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	readDir(w, folder, pathQueryParam, "File uploaded!")
}
