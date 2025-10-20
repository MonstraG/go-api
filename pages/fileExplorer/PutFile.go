package fileExplorer

import (
	"errors"
	"fmt"
	"go-api/infrastructure/helpers"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/reqRes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const hundredMegs = 100 << 20

func (controller *Controller) PutFile(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")

	folder := filepath.Join(controller.explorerRoot, pathQueryParam)

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

	for _, handler := range files {
		err := writeFile(folder, handler)
		if err != nil {
			myLog.Error.Logf(err.Error())
			renderExplorer(w, folder, pathQueryParam, err.Error())
			return
		}
	}

	renderExplorer(w, folder, pathQueryParam, "File(s) uploaded!")
}

func writeFile(folder string, handler *multipart.FileHeader) error {
	formFile, err := handler.Open()
	defer helpers.CloseSafely(formFile)
	if err != nil {
		message := fmt.Sprintf("Failed to open file handler: \n%v", err)
		return errors.New(message)
	}

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	path := filepath.Join(folder, handler.Filename)

	diskFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, filePermissions)
	defer helpers.CloseSafely(diskFile)
	if err != nil {
		message := fmt.Sprintf("Failed to open file on disk: \n%v", err)
		return errors.New(message)
	}

	_, err = io.Copy(diskFile, formFile)
	if err != nil {
		message := fmt.Sprintf("Failed to copy file: \n%v", err)
		return errors.New(message)

	}

	return nil
}
