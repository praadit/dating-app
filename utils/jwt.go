package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/praadit/dating-apps/config"
)

type Claims struct {
	jwt.StandardClaims
}

func GenerateLoginJwt(userID int) (string, time.Time, error) {
	expires := time.Now().Local().Add(time.Hour * 24)
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: expires.Unix(),
		Subject:   fmt.Sprint(userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(config.Config.JwtSign))
	return tokenString, expires, err
}

func ParseJWT(tokenString string) (*jwt.StandardClaims, error) {
	claims := &jwt.StandardClaims{}
	if tokenString == "" {
		return nil, ERR_INVALID_TOKEN()
	}

	var signingKey = []byte(config.Config.JwtSign)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) { return signingKey, nil })

	if err != nil {
		log.Println(err)
		return nil, ERR_FAILED_JWT_PARSE()
	}

	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, ERR_JWT_EXPIRED()
			}
		}
		return nil, ERR_INVALID_TOKEN()
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok {
		return claims, nil
	}

	return nil, ERR_INVALID_TOKEN()
}
