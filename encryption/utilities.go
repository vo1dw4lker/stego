package encryption

import (
	crnd "crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/pbkdf2"
)

func RandomBytes(bytes int) []byte {
	rnd := make([]byte, bytes)
	_, err := crnd.Read(rnd)
	if err != nil {
		panic(err)
	}
	return rnd
}

// PBKDF2 is a password-based key derivation function.
// Use nil for salt to generate a new one.
// Returns the derived key and used salt.
func PBKDF2(password string, salt []byte) (key, genSalt []byte) {
	if salt != nil {
		genSalt = salt
	} else {
		genSalt = RandomBytes(8)
	}
	key = pbkdf2.Key([]byte(password), genSalt, 600000, 32, sha256.New)
	return
}
