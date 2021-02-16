package users

import "golang.org/x/crypto/bcrypt"

// Encrypter password interface
type Encrypter interface {
	encrypt(password string) (string, error)
}

// BCryptEncrypter implementation
type BCryptEncrypter struct{}

func (bc *BCryptEncrypter) encrypt(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash), err
}
