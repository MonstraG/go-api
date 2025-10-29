package pages

import (
	"embed"
	"encoding/hex"
	"go-api/infrastructure/helpers"
	"go-api/infrastructure/myLog"
	"go-api/infrastructure/reqRes"
	"io"
	"net/http"
	"path/filepath"

	"golang.org/x/crypto/sha3"
)

//go:embed public/*
var publicFs embed.FS

func PublicHandler(w reqRes.MyResponseWriter, r *reqRes.MyRequest) {
	pathQueryParam := r.PathValue("path")
	filename := filepath.Join("public", pathQueryParam)
	http.ServeFileFS(w, &r.Request, publicFs, filename)
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
