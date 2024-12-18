package pages

import (
	"go-server/setup/reqRes"
	"path/filepath"
)

func PublicHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join("public", pathQueryParam)
	ServeFile(w, r, filename)
}
