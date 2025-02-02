package music

import (
	"errors"
	"fmt"
	"go-server/setup/reqRes"
	"html/template"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

var songsTemplate = template.Must(template.ParseFiles("pages/music/songsPartial.gohtml"))

type SongsData struct {
	Items           []SongItem
	HasUpNavigation bool
	UpPath          string
	Path            string
}

type SongItem struct {
	IsDir  bool
	IsSong bool
	Name   string
	Path   string
}

func GetSongs(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	folder := filepath.Join(r.AppConfig.SongsFolder, pathQueryParam)

	readDir(w, folder, pathQueryParam)
}

func readDir(w reqRes.MyWriter, folder string, query string) {
	dirAsFile, err := os.Open(folder)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			message := fmt.Sprintf("%s does not exist", folder)
			log.Printf(message)
			http.Error(w, message, http.StatusBadRequest)
			return
		}

		message := fmt.Sprintf("Failure to open folder '%s': \n%v", folder, err)
		log.Printf(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	fileInfo, err := dirAsFile.ReadDir(-1)
	if err != nil {
		// here I would check if this error is "NotADir" or something, but there seems to
		// not be a consistent / sane way to do this?
		// https://github.com/golang/go/issues/46734

		if errors.Is(err, fs.ErrNotExist) {
			message := fmt.Sprintf("Failure to read folde 2r '%s': \n%v", folder, err)
			log.Printf(message)
			http.Error(w, message, http.StatusInternalServerError)
		}

		message := fmt.Sprintf("Failure to read folder '%s': \n%v", folder, err)
		log.Printf(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	err = dirAsFile.Close()
	if err != nil {
		message := fmt.Sprintf("Failure to close folder '%s': \n%v", folder, err)
		log.Printf(message)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	var templatePageData = SongsData{
		Items:           make([]SongItem, len(fileInfo), len(fileInfo)),
		HasUpNavigation: query != "",
		UpPath:          path.Join(query, ".."),
		Path:            query,
	}

	for index, file := range fileInfo {
		isDir := file.IsDir()
		fileName := file.Name()
		templatePageData.Items[index] = SongItem{
			IsDir:  isDir,
			IsSong: !isDir && isSong(fileName),
			Name:   fileName,
			Path:   path.Join(query, fileName),
		}
	}

	slices.SortFunc(templatePageData.Items, func(a, b SongItem) int {
		if a.IsDir {
			return -1
		}
		if b.IsDir {
			return 1
		}
		return strings.Compare(a.Name, b.Name)
	})

	err = songsTemplate.Execute(w, templatePageData)
	if err != nil {
		message := fmt.Sprintf("Failed to render template: \n%v", err)
		log.Println(message)
		http.Error(w, message, http.StatusInternalServerError)
	}
}

func isSong(fileName string) bool {
	var extension = path.Ext(fileName)
	mimeType := mime.TypeByExtension(extension)
	return strings.HasPrefix(mimeType, "audio/")
}
