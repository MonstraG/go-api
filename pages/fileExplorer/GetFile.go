package fileExplorer

import (
	"go-api/infrastructure/reqRes"
	"net/http"
	"path/filepath"
)

func (controller *Controller) GetFile(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join(controller.explorerRoot, pathQueryParam)
	http.ServeFile(w, &r.Request, filename)
}
