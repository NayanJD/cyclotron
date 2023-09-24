package cryptus

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type UserJwtClaims struct {
	UserId  string `json:"user_id"`
	TokenId string `json:"token_id"`
	jwt.RegisteredClaims
}

func NewJwtToken(claims UserJwtClaims, secret string, expirationTime time.Time) (string, error) {

	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}
