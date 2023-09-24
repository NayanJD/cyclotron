package cryptus

import (
	"testing"
)

func BenchmarkBcryptCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BcryptCompare("$2a$08$SEvf4099ZQW2jZYeFXhgkO9b1tc2aLS/MfvLCHWG/3UlOWGUNJTJ6", "abcdefgh")
	}
}
