package encrypt

import (
	"time"

	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

// NonceValidator provides protection against replay attacks
type NonceValidator interface {
	// Validate checks if a nonce is valid and hasn't been seen before
	Validate(nonce []byte, sessionID types.EncryptionSessionID) error

	// Cleanup removes expired nonces
	Cleanup(before time.Time)
}
