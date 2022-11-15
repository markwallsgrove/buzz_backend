package models

import (
	"crypto/rand"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/argon2"
)

// CreateRandomPassword create a random password that contains random
// numbers, letters, & symbols.
func CreateRandomPassword() (string, error) {
	return password.Generate(21, 5, 5, false, false)
}

// CreatePasswordHash convert a password string into a hash. This is to ensure
// that the password is secure when within the database's memory & at rest. If
// the database data was stolen this should ensure the passwords will not be
// used for a impersonation attack. The passwords could also be used as a repeat
// attack on other services.
func CreatePasswordHash(password string) (string, error) {
	salt, err := generateRandomBytes(16)
	if err != nil {
		return "", err
	}

	// See https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
	// for configuration guidance
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		2,
		15*1024,
		1,
		32,
	)

	return string(hash), nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
