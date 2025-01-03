package pages

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-server/setup/reqRes"
	"net/http"
	"os"
)

func ServeFile(w reqRes.MyWriter, r *reqRes.MyRequest, filename string) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		message := fmt.Sprintf("Failed to stat file %s:\n%v\n", filename, err)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	if fileInfo.IsDir() {
		message := fmt.Sprintf("Failed to stat file: %s, it's a directory", filename)
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		message := fmt.Sprintf("Failed to read file %s:\n%v\n", filename, err)
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	eTag := calculateETag(file)
	w.Header().Set("ETag", eTag)

	readSeeker := bytes.NewReader(file)

	// this handles ETag matches inside
	http.ServeContent(w, &r.Request, filename, fileInfo.ModTime(), readSeeker)
}

var hasher = sha256.New()

// calculateETag generates a SHA-256 hash of the content and adds `W/` prefix to a hash to indicate weak comparison
func calculateETag(content []byte) string {
	hasher.Reset()
	hasher.Write(content)
	hash := hex.EncodeToString(hasher.Sum(nil))
	return fmt.Sprintf("W/\"%s\"", hash)
}
