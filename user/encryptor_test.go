package user

import (
	"testing"
)

func TestBCryptEncrypter(t *testing.T) {
	t.Run("Deliver error on bcrypt failure", func(t *testing.T) {
		sut := BCryptEncrypter{}
		_, err := sut.encrypt("test", 99)

		if err == nil {
			t.Errorf("got nil, want failure")
		}
	})

	t.Run("Delivers crypted password", func(t *testing.T) {
		sut := BCryptEncrypter{}
		cryptedPassword, err := sut.encrypt("test", 10)

		if err != nil {
			t.Errorf("got failure, want nil")
		}

		if cryptedPassword == "" {
			t.Errorf("got empty, want hashed password")
		}
	})
}
