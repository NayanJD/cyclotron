package cryptus

import (
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type TestJwtClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

var testClaims = UserJwtClaims{
	UserId: "abcdef-abcdef-abcdef-abcdef-abcdef",
}

var expirationTime = time.Now().Add(900 * time.Second)

func BenchmarkJwt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewJwtToken(testClaims, "abcd", expirationTime)
	}
}
