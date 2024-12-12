package login

import (
	"go-server/models"
	"go-server/setup/appConfig"
	"testing"
	"time"
)

// token is generated via https://jwt.io/#debugger-io
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MzQwMjExNzIsImlzcyI6ImdvLWFwaSIsInN1YiI6IkpvaG4ifQ.i_yYoxfV_7TYeWJKII26wUsaWnVwpdcCTlWKWgDva_U"

var user = models.User{Username: "John"}
var config = appConfig.AppConfig{JWTSecret: "secret"}

func TestCreateJwt(t *testing.T) {
	jwt, err := createJwt(user, config, preChosenTime)
	if err != nil {
		t.Fatal(err)
		return
	}

	want := token
	if jwt != want {
		t.Fatalf("want %s, got %s", want, jwt)
	}
}

// if changing this, make sure to remove the milliseconds as the issued at claim doesn't save that
func preChosenTime() time.Time {
	chosenTime, err := time.Parse(time.RFC3339, "2024-12-12T16:32:52Z")
	if err != nil {
		panic("Test has invalid date format\n" + err.Error())
	}
	return chosenTime
}

func TestValidateJwt(t *testing.T) {
	claims, ok := validateJWT(token, config)
	if !ok {
		t.Fatalf("Failed to validate JWT")
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
	wantSub := user.Username
	if gotSub != wantSub {
		t.Fatalf("Invalid sub, want %v, got %v", wantIssuedAt, gotIssuedAt)
	}
}
