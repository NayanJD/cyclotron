package cryptus

import (
	"math/rand"
	"testing"
)

var salt []byte
var msg string

func init() {
	salt = make([]byte, 5)
	rand.Read(salt)
	msg = "abcdefgh"
}

func BenchmarkScryptGetPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ScriptPassword(msg, salt)
	}
}
