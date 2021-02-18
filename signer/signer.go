package signer

// A Signer may sign a user with a token
type Signer interface {
	Sign(name string, email string, expiredAt int64) (string, error)
}
