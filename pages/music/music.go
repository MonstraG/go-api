package music

import (
	"go-server/pages"
	"go-server/setup/reqRes"
	"path/filepath"
)

func GetSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join(r.AppConfig.SongsFolder, pathQueryParam)
	pages.ServeFile(w, r, filename)
}
