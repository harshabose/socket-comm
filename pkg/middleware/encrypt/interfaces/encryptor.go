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

type CanEncrypt interface {
	// Encrypt encrypts a message between sender and receiver
	Encrypt(message message.Message) (message.Message, error)
}

type CanDecrypt interface {
	// Decrypt decrypts an encrypted message in-place
	Decrypt(message message.Message) (message.Message, error)
}

// Encryptor defines the interface for message encryption and decryption
type Encryptor interface {
	KeySetter
	CanEncrypt
	CanDecrypt

	// SetSessionID sets the session identifier for this encryption session
	SetSessionID(id types.EncryptionSessionID)

	// Ready checks if the encryptor is properly initialized and ready to use
	Ready() bool

	io.Closer
}
