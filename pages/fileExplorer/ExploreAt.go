package fileExplorer

import (
	"fmt"
	"go-api/infrastructure/helpers"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/reqRes"
	"go-api/pages"
	"mime"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

var fileExplorerTemplate = pages.ParsePartial("fileExplorer/fileExplorerPartial.gohtml")

type FilesData struct {
	Items         []FileItem
	Path          string
	ResultMessage string
}

type FileItem struct {
	IsDir  bool
	IsSong bool
	IsGoUp bool
	Name   string
	Path   string
	Size   string
}

func (controller *Controller) ExploreAt(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	folder := filepath.Join(controller.explorerRoot, pathQueryParam)

	renderExplorer(w, folder, pathQueryParam, "")
}

func renderExplorer(w reqRes.MyResponseWriter, fileSystemFolder string, queryFolder string, resultMessage string) {
	ok, dirEntries := helpers.ReadFolder(w, fileSystemFolder)
	if !ok {
		return
	}

	var templatePageData = FilesData{
		Items:         make([]FileItem, len(dirEntries)),
		Path:          queryFolder,
		ResultMessage: resultMessage,
	}

	for index, file := range dirEntries {
		isDir := file.IsDir()
		fileName := file.Name()
		templatePageData.Items[index] = FileItem{
			IsDir:  isDir,
			IsSong: !isDir && IsSong(fileName),
			Name:   fileName,
			Path:   path.Join(queryFolder, fileName),
			Size:   getFormattedFileSize(file),
		}
	}

	slices.SortFunc(templatePageData.Items, func(a, b FileItem) int {
		if a.IsDir && b.IsDir {
			return strings.Compare(a.Name, b.Name)
		}
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
		templatePageData.Items = slices.Insert(templatePageData.Items, 0, FileItem{
			IsDir:  true,
			IsSong: false,
			IsGoUp: true,
			Name:   "..",
			Path:   path.Join(queryFolder, ".."),
		})
	}

	w.RenderTemplate(fileExplorerTemplate, templatePageData)
}

func IsSong(fileName string) bool {
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
		myLog.Info.Logf("Failure to get file info for '%s': \n\t%v", dirEntry.Name(), err)
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
