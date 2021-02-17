package signer

// A Signer may sign a user with a token
type Signer interface {
	Sign() error
}
