package encryptor

import (
	"io"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

// Encryptor defines the interface for message encryption and decryption
type Encryptor interface {
	// SetKeys configures the encryption and decryption keys
	SetKeys(encryptKey, decryptKey types.Key) error

	// SetSessionID sets the session identifier for this encryption session
	SetSessionID(id types.SessionID)

	// Encrypt encrypts a message between sender and receiver
	Encrypt(senderID, receiverID string, message message.Message) (message.Message, error)

	// Decrypt decrypts an encrypted message in-place
	Decrypt(message message.Message) error

	// Ready checks if the encryptor is properly initialized and ready to use
	Ready() bool

	io.Closer
}
