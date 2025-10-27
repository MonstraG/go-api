package helpers

import (
	"errors"
	"fmt"
	"go-api/infrastructure/reqRes"
	"io/fs"
	"net/http"
	"os"
)

func ReadFolder(w reqRes.MyResponseWriter, fileSystemFolder string) (bool, []os.DirEntry) {
	dirAsFile, err := os.Open(fileSystemFolder)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			message := fmt.Sprintf("%s does not exist", fileSystemFolder)
			w.Error(message, http.StatusBadRequest)
			return false, nil
		}

		message := fmt.Sprintf("Failure to open folder '%s': \n%v", fileSystemFolder, err)
		w.Error(message, http.StatusInternalServerError)
		return false, nil
	}

	dirEntries, err := dirAsFile.ReadDir(-1)
	if err != nil {
		// here I would check if this error is "NotADir" or something, but there seems to
		// not be a consistent / sane way to do this?
		// https://github.com/golang/go/issues/46734

		if errors.Is(err, fs.ErrNotExist) {
			message := fmt.Sprintf("Failure to read folder '%s': \n%v", fileSystemFolder, err)
			w.Error(message, http.StatusInternalServerError)
			return false, nil
		}

		message := fmt.Sprintf("Failure to read folder '%s': \n%v", fileSystemFolder, err)
		w.Error(message, http.StatusInternalServerError)
		return false, nil
	}

	err = dirAsFile.Close()
	if err != nil {
		message := fmt.Sprintf("Failure to close folder '%s': \n%v", fileSystemFolder, err)
		w.Error(message, http.StatusInternalServerError)
		return false, nil
	}

	return true, dirEntries
}
