package pages

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go-server/setup/reqRes"
	"log"
	"net/http"
	"os"
)

func ServeFile(w reqRes.MyWriter, r *reqRes.MyRequest, filename string) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			message := fmt.Sprintf("Failed to stat file %s:\n%v", filename, err)
			log.Println(message)
			w.Error(http.StatusBadRequest, "File not found")
			return
		}

		w.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to stat file %s:\n%v", filename, err))
		return
	}

	if fileInfo.IsDir() {
		w.Error(http.StatusBadRequest, fmt.Sprintf("Failed to get file: %s, it's a directory", filename))
		return
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		w.Error(http.StatusInternalServerError, fmt.Sprintf("Failed to read file %s:\n%v", filename, err))
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
