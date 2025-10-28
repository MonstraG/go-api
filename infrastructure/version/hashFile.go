package version

import (
	"encoding/hex"
	"go-api/infrastructure/helpers"
	"go-api/infrastructure/myLog"
	"io"
	"os"

	"golang.org/x/crypto/sha3"
)

var StylesHash string

func init() {
	StylesHash = hashFile("public/styles.css")
}

var hasher = sha3.NewShake128()

func hashFile(filePath string) string {
	file, err := os.Open(filePath)
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
