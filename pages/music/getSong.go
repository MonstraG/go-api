package music

import (
	"go-server/setup/reqRes"
	"net/http"
	"path/filepath"
)

func GetSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join(r.AppConfig.SongsFolder, pathQueryParam)
	http.ServeFile(w, &r.Request, filename)
}
