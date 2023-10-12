package cryptus

import (
	"golang.org/x/crypto/scrypt"
)

func ScriptPassword(psw string, salt []byte) string {
	dk, _ := scrypt.Key([]byte(psw), salt, 1<<14, 8, 1, 32)
	return string(dk)
}
