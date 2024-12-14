package music

import (
	"fmt"
	"go-server/pages/music/ytDlp"
	"go-server/setup/reqRes"
	"log"
	"net/http"
)

func PostHandler(w reqRes.MyWriter, r *reqRes.MyRequest) {
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		log.Printf("Failed to parse form:\n%v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for key, value := range r.Form {
		fmt.Printf("%s: %s\n", key, value)
	}

	songUrl := r.Form.Get("songUrl")
	if songUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.WriteSilent([]byte("songUrl missing"))
		return
	}

	sanitizedUrl := ytDlp.SanitizeUrl(songUrl)

	ytDlp.Download(sanitizedUrl)
}
