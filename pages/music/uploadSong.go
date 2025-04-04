package music

import (
	"fmt"
	"go-server/setup/reqRes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const hundredMegs = 100 << 20

func PutSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")

	err := r.ParseMultipartForm(hundredMegs)
	if err != nil {
		message := fmt.Sprintf("Failed to parse form: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		message := fmt.Sprintf("Failed to retrieve file: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			message := fmt.Sprintf("Failed to close file: \n%v", err)
			log.Println(message)
		}
	}(file)

	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	folder := filepath.Join(r.AppConfig.SongsFolder, pathQueryParam)
	path := filepath.Join(folder, handler.Filename)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		message := fmt.Sprintf("Failed to retrieve file: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			message := fmt.Sprintf("Failed to close file: \n%v", err)
			log.Println(message)
		}
	}(f)

	_, err = io.Copy(f, file)
	if err != nil {
		message := fmt.Sprintf("Failed to save file: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
	}

	readDir(w, folder, pathQueryParam, "File uploaded!")
}
