package myJwt

import (
	"go-api/models"
	"go-api/setup/appConfig"
	"strings"
	"testing"
	"time"
)

// if changing this, make sure to remove the milliseconds as the issued at claim doesn't save that
func preChosenTime() time.Time {
	chosenTime, err := time.Parse(time.RFC3339, "2024-12-12T16:32:52Z")
	if err != nil {
		panic("Test has invalid date format\n" + err.Error())
	}
	return chosenTime
}

var config = appConfig.AppConfig{JWTSecret: "random-32-bit-secret-for-testing"}

var myJwt = CreateMyJwt(config, preChosenTime)

// token is generated via https://jwt.io/#debugger-io
// (JWT Encoder tab)
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MzQwMjExNzIsImlzcyI6ImdvLWFwaSIsIm5hbWUiOiJKb2huIiwic3ViIjoiYWQxZTg5OTQtNTc2MC00OWM1LWI3MzEtNWI4YmE2MmQyM2ZlIn0.Bi9yRd8Fj4nZ0lmX3yx4S6v-k5xkfpR5omtKhJmDPG0"

// uuid.NewString()
var userId = "ad1e8994-5760-49c5-b731-5b8ba62d23fe"
var user = models.User{Username: "John", ID: userId}

func TestCreateJwt(t *testing.T) {
	jwt, err := myJwt.CreateJwt(user)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
		return
	}

	want := token
	if !strings.HasPrefix(jwt, want) {
		t.Fatalf("want \n%s,\ngot \n%s\n", want, jwt)
	}
}

func TestValidateJwt(t *testing.T) {
	claims, err := myJwt.ValidateJWT(token)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	gotIssuer, err := claims.GetIssuer()
	if err != nil {
		t.Fatalf("Failed to get issuer, %v", err)
	}
	wantIssuer := issuer
	if gotIssuer != wantIssuer {
		t.Fatalf("Invalid issuer, want %s, got %s", wantIssuer, gotIssuer)
	}

	gotIssuedAt, err := claims.GetIssuedAt()
	if err != nil {
		t.Fatalf("Failed to get issuer, %v", err)
	}
	wantIssuedAt := preChosenTime().UTC()
	if gotIssuedAt.UTC() != wantIssuedAt {
		t.Fatalf("Invalid issuedAt, want %v, got %v", wantIssuedAt.UTC(), gotIssuedAt.UTC())
	}

	gotSub, err := claims.GetSubject()
	if err != nil {
		t.Fatalf("Failed to get sub, %v", err)
	}
	wantSub := user.ID
	if gotSub != wantSub {
		t.Fatalf("Invalid sub, want %v, got %v", wantIssuedAt, gotIssuedAt)
	}
}
