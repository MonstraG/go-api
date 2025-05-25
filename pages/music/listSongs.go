package music

import (
	"errors"
	"fmt"
	"go-api/setup/myLog"
	"go-api/setup/reqRes"
	"html/template"
	"io/fs"
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
	Items         []SongItem
	Path          string
	ResultMessage string
}

type SongItem struct {
	IsDir  bool
	IsSong bool
	IsGoUp bool
	Name   string
	Path   string
	Size   string
}

func (controller *Controller) GetSongs(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	folder := filepath.Join(controller.songsFolder, pathQueryParam)

	readDir(w, folder, pathQueryParam, "")
}

func readDir(w reqRes.MyWriter, fileSystemFolder string, queryFolder string, resultMessage string) {
	dirAsFile, err := os.Open(fileSystemFolder)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			message := fmt.Sprintf("%s does not exist", fileSystemFolder)
			w.Error(message, http.StatusBadRequest)
			return
		}

		message := fmt.Sprintf("Failure to open folder '%s': \n%v", fileSystemFolder, err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	dirEntries, err := dirAsFile.ReadDir(-1)
	if err != nil {
		// here I would check if this error is "NotADir" or something, but there seems to
		// not be a consistent / sane way to do this?
		// https://github.com/golang/go/issues/46734

		if errors.Is(err, fs.ErrNotExist) {
			message := fmt.Sprintf("Failure to read folder '%s': \n%v", fileSystemFolder, err)
			w.Error(message, http.StatusInternalServerError)
			return
		}

		message := fmt.Sprintf("Failure to read folder '%s': \n%v", fileSystemFolder, err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	err = dirAsFile.Close()
	if err != nil {
		message := fmt.Sprintf("Failure to close folder '%s': \n%v", fileSystemFolder, err)
		w.Error(message, http.StatusInternalServerError)
		return
	}

	var templatePageData = SongsData{
		Items:         make([]SongItem, len(dirEntries)),
		Path:          queryFolder,
		ResultMessage: resultMessage,
	}

	for index, file := range dirEntries {
		isDir := file.IsDir()
		fileName := file.Name()
		templatePageData.Items[index] = SongItem{
			IsDir:  isDir,
			IsSong: !isDir && isSong(fileName),
			Name:   fileName,
			Path:   path.Join(queryFolder, fileName),
			Size:   getFormattedFileSize(file),
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

	canGoUp := queryFolder != ""
	if canGoUp {
		templatePageData.Items = slices.Insert(templatePageData.Items, 0, SongItem{
			IsDir:  true,
			IsSong: false,
			IsGoUp: true,
			Name:   "..",
			Path:   path.Join(queryFolder, ".."),
		})
	}

	w.RenderTemplate(songsTemplate, templatePageData)
}

func isSong(fileName string) bool {
	var extension = path.Ext(fileName)
	mimeType := mime.TypeByExtension(extension)
	return strings.HasPrefix(mimeType, "audio/")
}

var sizes = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}

func getFormattedFileSize(dirEntry os.DirEntry) string {
	if dirEntry.IsDir() {
		return ""
	}

	fileInfo, err := dirEntry.Info()
	if err != nil {
		myLog.Info.Logf("Failure to get file info for '%s': \n%v", dirEntry.Name(), err)
		return ""
	}

	fileSize := fileInfo.Size()
	return formatFileSize(float64(fileSize))
}

// adjusted from https://ahmadrosid.com/cheatsheet/go/FormatFileSize
func formatFileSize(size float64) string {
	const base = 1024.0
	unitsLimit := len(sizes)
	i := 0
	for size >= base && i < unitsLimit {
		size = size / base
		i++
	}

	f := "%.0f %s"
	if i > 1 {
		f = "%.2f %s"
	}

	return fmt.Sprintf(f, size, sizes[i])
}
