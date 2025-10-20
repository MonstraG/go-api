package fileExplorer

import (
	"fmt"
	"go-api/infrastructure/helpers"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/reqRes"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const hundredMegs = 100 << 20

func (controller *Controller) PutFile(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")

	err := r.ParseMultipartForm(hundredMegs)
	if err != nil {
		message := fmt.Sprintf("Failed to parse form: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		w.Error("No files uploaded", http.StatusBadRequest)
		return
	}

	folder := filepath.Join(controller.explorerRoot, pathQueryParam)

	for _, handler := range files {
		formFile, err := handler.Open()
		if err != nil {
			message := fmt.Sprintf("Failed to retrieve file: \n%v", err)
			w.Error(message, http.StatusInternalServerError)
		}

		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		path := filepath.Join(folder, handler.Filename)

		diskFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, filePermissions)
		if err != nil {
			message := fmt.Sprintf("Failed to retrieve file: \n%v", err)
			myLog.Error.Logf(message)
			renderExplorer(w, folder, pathQueryParam, message)
			return
		}

		_, err = io.Copy(diskFile, formFile)
		if err != nil {
			message := fmt.Sprintf("Failed to save file: \n%v", err)
			myLog.Error.Logf(message)
			renderExplorer(w, folder, pathQueryParam, message)
			return
		}

		helpers.CloseSafely(formFile)
		helpers.CloseSafely(diskFile)
	}

	renderExplorer(w, folder, pathQueryParam, "File(s) uploaded!")
}
