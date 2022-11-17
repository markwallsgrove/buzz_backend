package security

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type claims struct {
	*jwt.RegisteredClaims
	UserInfo interface{}
}

// CreateToken create a JWT that represents a user's session
func CreateToken(email string, secret string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	exp := time.Now().Add(time.Hour * 24)
	token.Claims = &claims{
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   email,
		},
	}

	val, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return val, nil
}
