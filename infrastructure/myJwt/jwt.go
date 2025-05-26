package myJwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-api/infrastructure/appConfig"
	"go-api/infrastructure/models"
	"time"
)

const issuer = "go-api"
const Cookie = "jwtToken"
const MaxAge = 3600 * 24

type Service struct {
	now       func() time.Time
	secretKey []byte
}

type MyCustomClaims struct {
	Username string `json:"name"`
	jwt.RegisteredClaims
}

func CreateMyJwt(config appConfig.AppConfig, now func() time.Time) Service {
	return Service{
		now:       now,
		secretKey: []byte(config.JWTSecret),
	}
}

func (myJwt *Service) CreateJwt(user models.User) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  issuer,
		"sub":  user.ID,
		"iat":  myJwt.now().Unix(),
		"name": user.Username,
	})

	return jwtToken.SignedString(myJwt.secretKey)
}

var validationMethod = []string{"HS256"}

func (myJwt *Service) ValidateJWT(tokenString string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return myJwt.secretKey, nil
	}, jwt.WithValidMethods(validationMethod), jwt.WithIssuer(issuer), jwt.WithIssuedAt(), jwt.WithLeeway(5*time.Second))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return nil, ErrTokenClaimCastFailed
	}

	return claims, nil
}

var (
	ErrTokenInvalid         = errors.New("token is invalid")
	ErrTokenClaimCastFailed = errors.New("failed to cast")
)
