package pages

import (
	"go-server/setup/reqRes"
	"net/http"
	"path/filepath"
)

func PublicHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join("public", pathQueryParam)
	http.ServeFile(w, &r.Request, filename)
}
