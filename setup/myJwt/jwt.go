package myJwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-server/models"
	"go-server/setup/appConfig"
	"time"
)

const issuer = "go-api"
const Cookie = "jwtToken"
const MaxAge = 3600 * 24

type MyJwt struct {
	now func() time.Time
}

// Singleton exists only to set "default" now time
// which, in turn, exists only to support passing custom time in tests :(
var Singleton = MyJwt{
	now: time.Now,
}

func (myJwt *MyJwt) CreateJwt(user models.User, config appConfig.AppConfig) (string, error) {
	key := []byte(config.JWTSecret)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": issuer,
		"sub": user.Username,
		"iat": myJwt.now().Unix(),
	})

	return jwtToken.SignedString(key)
}

func (myJwt *MyJwt) ValidateJWT(tokenString string, config appConfig.AppConfig) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.JWTSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer(issuer), jwt.WithIssuedAt())

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	iat, err := token.Claims.GetIssuedAt()
	if err != nil {
		return nil, ErrIssuedAtMissing
	}

	expired := iat.Time.Add(MaxAge * time.Second).Before(myJwt.now())
	if expired {
		return nil, ErrTokenExpired
	}

	return token.Claims, nil
}

var (
	ErrTokenInvalid    = errors.New("token is invalid")
	ErrIssuedAtMissing = errors.New("issued at time is missing")
	ErrTokenExpired    = errors.New("token is expired")
)
