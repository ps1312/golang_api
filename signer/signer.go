package signer

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// A Signer may sign a user with a token
type Signer interface {
	Sign(name string, email string, expiredAt int64) (string, error)
}

// ErrJWTSigner error const
const ErrJWTSigner = "Error using JWTSigner adapter"

// A JWTSigner Signer implementation
type JWTSigner struct {
	Key interface{}
}

// Sign function implementation
func (js *JWTSigner) Sign(name string, email string, expiredAt int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["name"] = name
	claims["email"] = email
	claims["exp"] = expiredAt

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString(js.Key)

	if err != nil {
		return "", errors.New(ErrJWTSigner)
	}

	return token, nil
}
