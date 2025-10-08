package fileExplorer

import (
	"errors"
	"fmt"
	"go-api/infrastructure/reqRes"
	"net/http"
	"os"
	"path/filepath"
)

func (controller *Controller) DeleteFile(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	pathToRemove := filepath.Join(controller.explorerRoot, pathQueryParam)

	stat, err := os.Stat(pathToRemove)
	if errors.Is(err, os.ErrNotExist) {
		message := fmt.Sprintf("Path not found: \n%v", err)
		w.Error(message, http.StatusBadRequest)
	}

	isDir := stat.IsDir()
	if isDir {
		err = os.RemoveAll(pathToRemove)
	} else {
		err = os.Remove(pathToRemove)
	}
	if err != nil {
		message := fmt.Sprintf("Failed to remove: \n%v", err)
		w.Error(message, http.StatusInternalServerError)
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

	renderExplorer(w, fileSystemFolder, queryFolder, resultMessage)
}
