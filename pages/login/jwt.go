package login

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-server/models"
	"go-server/setup/appConfig"
	"log"
	"time"
)

const issuer = "go-api"

func createJwt(user models.User, config appConfig.AppConfig, now func() time.Time) (string, error) {
	key := []byte(config.JWTSecret)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": issuer,
		"sub": user.Username,
		"iat": now().Unix(),
	})

	return jwtToken.SignedString(key)
}

func validateJWT(tokenString string, config appConfig.AppConfig) (jwt.Claims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.JWTSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer(issuer), jwt.WithIssuedAt())

	if err != nil {
		log.Printf("Failed to parse token:\n%v", err)
		return nil, false
	}

	return token.Claims, token.Valid
}
