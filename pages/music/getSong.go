package music

import (
	"go-api/infrastructure/reqRes"
	"net/http"
	"path/filepath"
)

func (controller *Controller) GetSongHandler(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join(controller.songsFolder, pathQueryParam)
	http.ServeFile(w, &r.Request, filename)
}
