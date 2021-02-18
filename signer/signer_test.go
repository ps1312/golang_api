package signer

import (
	"errors"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

const ErrJWTSigner = "Error using JWTSigner adapter"

type JWTSigner struct {
	key interface{}
}

func (js *JWTSigner) Sign() (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, nil)

	token, err := at.SignedString(js.key)

	if err != nil {
		return "", errors.New(ErrJWTSigner)
	}

	return token, nil
}

func TestJWTSigner(t *testing.T) {
	t.Run("Delivers error on sign failure", func(t *testing.T) {
		sut := JWTSigner{key: []string{}}

		_, got := sut.Sign()
		want := ErrJWTSigner

		if got.Error() != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Delivers token string on sign success", func(t *testing.T) {
		sut := JWTSigner{key: []byte("any secret key")}

		got, err := sut.Sign()

		if err != nil {
			t.Errorf("got %q, want nil", err)
		}

		if got == "" {
			t.Error("got '', want token")
		}
	})
}
