package security

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

// CreateToken create a JWT that represents a user's session
func CreateToken(userId int, secret string) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	exp := time.Now().Add(time.Hour * 24)
	token.Claims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(exp),
		Subject:   fmt.Sprintf("%d", userId),
	}

	val, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return val, nil
}

// ValidateToken validate a JWT token and return the claims
func ValidateToken(token, secret string) (*jwt.RegisteredClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return &jwt.RegisteredClaims{}, err
	}

	if !jwtToken.Valid {
		return &jwt.RegisteredClaims{}, errors.New("invalid token signature")
	}

	c, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return &jwt.RegisteredClaims{}, errors.New("invalid token structure")
	}

	if c.VerifyExpiresAt(time.Now(), true) == false {
		return &jwt.RegisteredClaims{}, errors.New("token has expired")
	}

	return c, nil
}
