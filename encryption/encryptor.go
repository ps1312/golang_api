package encryption

import "golang.org/x/crypto/bcrypt"

// Encrypter password interface
type Encrypter interface {
	Encrypt(password string, cost int) (string, error)
}

// BCryptEncrypter implementation
type BCryptEncrypter struct{}

func (bc *BCryptEncrypter) Encrypt(password string, cost int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash), err
}
