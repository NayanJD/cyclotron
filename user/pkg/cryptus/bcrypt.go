package cryptus

import (
	"golang.org/x/crypto/bcrypt"
)

func BcryptCompare(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
