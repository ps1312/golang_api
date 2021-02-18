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

func (js *JWTSigner) Sign(name string, email string, expiredAt int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["name"] = name
	claims["email"] = email
	claims["exp"] = expiredAt

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := at.SignedString(js.key)

	if err != nil {
		return "", errors.New(ErrJWTSigner)
	}

	return token, nil
}

func TestJWTSigner(t *testing.T) {
	t.Run("Delivers error on sign failure", func(t *testing.T) {
		sut := JWTSigner{key: []string{}}

		_, got := sut.Sign("", "", 0)
		want := ErrJWTSigner

		if got.Error() != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Delivers token string on sign success", func(t *testing.T) {
		sut := JWTSigner{key: []byte("any secret key")}

		got, err := sut.Sign("any-name", "any-email@mail.com", 1234)
		// fixed token for user with above credentials for HSA256 and `any secret key` secret and expired at 1234
		want := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFueS1lbWFpbEBtYWlsLmNvbSIsImV4cCI6MTIzNCwibmFtZSI6ImFueS1uYW1lIn0.qjAqj2FLf3yigHsQK13MgpwK8z6VwulpwZQ2IYvfpUY"

		if err != nil {
			t.Errorf("got %q, want nil", err)
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
