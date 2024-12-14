package ytDlp

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func SanitizeUrl(url string) string {
	index := strings.Index(url, "&")
	if index != -1 {
		return url[:index]
	}
	return url
}

const ytDlpBinary = "yt-dlp"
const fileNamePattern = "%(id)s.%(ext)s"

func Download(url string) {
	cmd := exec.Command(ytDlpBinary, []string{url, "-x", "-o", fileNamePattern}...)
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Printf("Failed to start download command:\n%v\n", err)
	}

	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Printf("Error from wait on command:\n%v\n", err)
		}
	}()
}
