package users

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// Encrypter password interface
type Encrypter interface {
	encrypt(password string) (string, error)
}

type BCryptEncrypter struct{}

func (bc *BCryptEncrypter) encrypt(password string, cost int) (string, error) {
	_, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return "", err
}

func TestBCryptEncrypter(t *testing.T) {
	t.Run("Deliver error on bcrypt failure", func(t *testing.T) {
		sut := BCryptEncrypter{}
		_, err := sut.encrypt("test", 99)

		if err == nil {
			t.Errorf("got nil, want failure")
		}
	})
}
