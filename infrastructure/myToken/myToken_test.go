package myToken

import (
	"crypto/aes"
	"go-api/infrastructure/models"
	"testing"
	"time"
	"uuid"
)

// if changing this, make sure to remove the milliseconds as the issued at claim doesn't save that
func preChosenTime() time.Time {
	chosenTime, err := time.Parse(time.RFC3339, "2024-12-12T16:32:52Z")
	if err != nil {
		panic("Test has invalid date format\n" + err.Error())
	}
	return chosenTime
}

const goldenToken = "AAAAAAAAAAAAAAAAAAAAADiFcNT2VxQ3Hzgd677dOpvfxMf9aUeAq/hixphoTGMiWyJHvE9CaTwtd6KSfYSTnuz9149bOaKh3mVdJbVSQgg="
const secret = "random-32-bit-secret-for-testing"

var cipherBlock, cipherErr = aes.NewCipher([]byte(secret))

type EmptyReader struct{}

func (EmptyReader EmptyReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

const practicallyForever = time.Hour * 24 * 365 * 100

var service = Service{
	now:         preChosenTime,
	cipherBlock: cipherBlock,
	ivReader:    EmptyReader{},
	maxAge:      practicallyForever,
}

var userId = uuid.MustParse("01a0200b-3103-779d-a68c-d72cb7ca1e71")
var user = models.User{Username: "John", ID: userId}

func TestCreateToken(t *testing.T) {
	if cipherErr != nil {
		t.Fatalf("Failed to create cipher: %v", cipherErr)
		return
	}

	tokenPayload := service.NewTokenPayload(user)
	createdCookie, err := service.CreateCookie(tokenPayload)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
		return
	}

	if createdCookie != goldenToken {
		t.Fatalf("want \n%s,\ngot \n%s\n", goldenToken, createdCookie)
	}
}

func TestParseToken(t *testing.T) {
	if cipherErr != nil {
		t.Fatalf("Failed to create cipher: %v", cipherErr)
		return
	}

	payload, err := service.ParseCookie(goldenToken)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	gotIssuedAt := time.Unix(payload.Iat, 0)
	wantIssuedAt := preChosenTime()
	if !gotIssuedAt.Equal(wantIssuedAt) {
		t.Fatalf("Invalid issuedAt, want %v, got %v", wantIssuedAt.UTC(), gotIssuedAt.UTC())
	}

	wantSub := user.ID
	if payload.Sub != wantSub {
		t.Fatalf("Invalid sub, want %v, got %v", wantIssuedAt, gotIssuedAt)
	}
}
