package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string, salt string) string {
	bytes := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(bytes)
}

const saltLength = 32

func NewSalt() string {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		panic("they say it never happens")
	}
	return base64.RawStdEncoding.EncodeToString(salt)
}
