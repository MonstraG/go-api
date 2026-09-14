package pages

import (
	"embed"
	"encoding/hex"
	"go-api/infrastructure/helpers"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/reqRes"
	"go-api/infrastructure/version"
	"io"
	"net/http"
	"path"
	"path/filepath"

	"golang.org/x/crypto/sha3"
)

//go:embed public/*
var publicFs embed.FS

func GetPublicFile(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join("public", pathQueryParam)

	f, err := publicFs.Open(filename)
	if err != nil {
		http.NotFound(w, &r.Request)
		return
	}
	defer helpers.CloseSafely(f)

	rs, ok := f.(io.ReadSeeker)
	if !ok {
		http.Error(w, "file is not seekable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=3600")

	http.ServeContent(
		w,
		&r.Request,
		path.Base(filename),
		version.AppBuildTime,
		rs,
	)
}

var StylesHash string

func init() {
	StylesHash = hashFile("public/styles.css")
}

var hasher = sha3.NewShake128()

func hashFile(filePath string) string {
	file, err := publicFs.Open(filePath)
	if err != nil {
		myLog.Fatal.Logf(err.Error())
	}

	defer helpers.CloseSafely(file)

	_, err = io.Copy(hasher, file)
	if err != nil {
		myLog.Fatal.Logf(err.Error())
	}

	sum := hasher.Sum(nil)
	shortHash := hex.EncodeToString(sum)[:8]
	return shortHash
}
