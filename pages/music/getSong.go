package music

import (
	"go-api/setup/reqRes"
	"net/http"
	"path/filepath"
)

func (controller *Controller) GetSongHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join(controller.songsFolder, pathQueryParam)
	http.ServeFile(w, &r.Request, filename)
}
