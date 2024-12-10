package login

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-server/setup"
	"log"
)

func createJwt(config setup.AppConfig) (string, error) {
	key := []byte(config.JWTSecret)
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	return jwtToken.SignedString(key)
}

func validateJWT(config setup.AppConfig, tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return config, nil
	})

	if err != nil {
		log.Fatalf("Failed to parse token:\n%v", err)
	}

	return token.Valid
}
