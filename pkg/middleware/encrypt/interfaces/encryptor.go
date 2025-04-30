package interfaces

import (
	"io"

	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type KeySetter interface {
	// SetKeys configures the encryption and decryption keys
	SetKeys(encryptorKey, decryptorKey types.Key) error
}

type KeyGetter interface {
	GetKeys() (encKey, decKey types.Key, err error)
}

// Encryptor defines the interface for message encryption and decryption
type Encryptor interface {
	KeySetter

	// SetSessionID sets the session identifier for this encryption session
	SetSessionID(id types.EncryptionSessionID)

	// Encrypt encrypts a message between sender and receiver
	Encrypt(senderID, receiverID string, message message.Message) (message.Message, error)

	// Decrypt decrypts an encrypted message in-place
	Decrypt(message message.Message) error

	// Ready checks if the encryptor is properly initialized and ready to use
	Ready() bool

	io.Closer
}
