package signer

import (
	"testing"
)

func TestJWTSigner(t *testing.T) {
	t.Run("Delivers error on sign failure", func(t *testing.T) {
		sut := JWTSigner{Key: []string{}}

		_, got := sut.Sign("", "", 0)
		want := ErrJWTSigner

		if got.Error() != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Delivers token string on sign success", func(t *testing.T) {
		sut := JWTSigner{Key: []byte("any secret key")}

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
