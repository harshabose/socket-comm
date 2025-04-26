package keyprovider

import "crypto/ed25519"

// KeyProvider interface for secure access to cryptographic keys
type KeyProvider interface {
	GetSigningKey() ed25519.PrivateKey
	GetVerificationKey() ed25519.PublicKey
	Close() error
}
