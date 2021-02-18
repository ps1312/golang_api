package signer

import (
	"errors"
	"testing"
)

const ErrJWTSigner = "Error using JWTSigner adapter"

type JWTSigner struct{}

func (js *JWTSigner) Sign() (string, error) {
	return "", errors.New(ErrJWTSigner)
}

func TestJWTSigner(t *testing.T) {
	sut := JWTSigner{}

	_, got := sut.Sign()
	want := ErrJWTSigner

	if got.Error() != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
