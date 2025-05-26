package pages

import (
	"go-api/infrastructure/reqRes"
	"net/http"
	"path/filepath"
)

func PublicHandler(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join("public", pathQueryParam)
	http.ServeFile(w, &r.Request, filename)
}
