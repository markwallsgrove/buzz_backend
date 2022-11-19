package security

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/gob"
	"errors"

	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/argon2"
)

// CreateRandomPassword create a random password that contains random
// numbers, letters, & symbols.
func CreateRandomPassword() (string, error) {
	return password.Generate(21, 5, 5, false, false)
}

// PasswordHash all the details and the resulting hash are stored
// within the struct. These details can be stored and then later
// used to generate a new hash from a given password and the results
// can be compared against the original hash.
type PasswordHash struct {
	Hash    []byte
	Salt    []byte
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// CreatePasswordHash convert a password string into a hash. This is to ensure
// that the password is secure when within the database's memory & at rest. If
// the database data was stolen this should ensure the passwords will not be
// used for a impersonation attack. The passwords could also be used as a repeat
// attack on other services.
func CreatePasswordHash(password string) ([]byte, error) {
	// See https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
	// for configuration guidance
	var time uint32 = 2
	var mb uint32 = 1024
	var memory uint32 = 15 * mb
	var threads uint8 = 1
	var keyLen uint32 = 32

	salt, err := generateRandomBytes(16)
	if err != nil {
		return []byte{}, err
	}

	passwordHash := createArgon2IdKeyHash(
		password,
		time,
		memory,
		threads,
		keyLen,
		salt,
	)

	var encodedPasswordHash bytes.Buffer
	encoder := gob.NewEncoder(&encodedPasswordHash)

	if err := encoder.Encode(passwordHash); err != nil {
		return []byte{}, err
	}

	return encodedPasswordHash.Bytes(), nil
}

// createArgon2IdKeyHash create a secure hash of the password
func createArgon2IdKeyHash(
	password string,
	time uint32,
	memory uint32,
	threads uint8,
	keyLen uint32,
	salt []byte,
) *PasswordHash {
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		keyLen,
	)

	// Store all the details within the PasswordHash
	// to be able to recreate the hash for verification.
	return &PasswordHash{
		Hash:    hash,
		Salt:    salt,
		Time:    time,
		Memory:  memory,
		Threads: threads,
		KeyLen:  keyLen,
	}
}

// ErrInvalidPassword after hashing the password the result was not
// the same as the original hash
var ErrInvalidPassword = errors.New("password does not match hash")

// VerifyPasswordHash verify a password against the original hash
func VerifyPasswordHash(password string, hash []byte) error {
	passwordHash := &PasswordHash{}
	decoder := gob.NewDecoder(bytes.NewBuffer(hash))
	err := decoder.Decode(&passwordHash)

	if err != nil {
		return err
	}

	passwordHashToCheck := createArgon2IdKeyHash(
		password,
		passwordHash.Time,
		passwordHash.Memory,
		passwordHash.Threads,
		passwordHash.KeyLen,
		passwordHash.Salt,
	)

	if subtle.ConstantTimeCompare(passwordHash.Hash, passwordHashToCheck.Hash) == 0 {
		return ErrInvalidPassword
	}

	return nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type LoginResponse struct {
	Result string
}
