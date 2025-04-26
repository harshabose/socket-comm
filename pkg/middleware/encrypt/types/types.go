package types

// Crypto-related type definitions for improved type safety
type (
	PrivateKey           [32]byte
	PublicKey            [32]byte
	Salt                 [16]byte
	SessionID            [16]byte
	Nonce                [12]byte
	Key                  [32]byte
	KeyExchangeProtocol  string
	KeyExchangeSessionID string
)

// EncryptionProtocol defines the capabilities of a specific protocol version
type EncryptionProtocol struct {
	Version             uint8               `json:"version"`
	CipherSuite         string              `json:"cipher_suite"`
	KeyExchangeProtocol KeyExchangeProtocol `json:"key_exchange"`
	Authenticated       bool                `json:"authenticated"`
	SupportedExtensions map[string]string   `json:"supported_extensions,omitempty"`
}

// Protocol versions
var (
	// ProtocolV1 is the original protocol
	ProtocolV1 = EncryptionProtocol{
		Version:             1,
		CipherSuite:         "AES-256-GCM",
		KeyExchangeProtocol: "curve25519-ed25519",
		Authenticated:       true,
	}

	// ProtocolV2 adds support for key rotation
	ProtocolV2 = EncryptionProtocol{
		Version:             2,
		CipherSuite:         "AES-256-GCM",
		KeyExchangeProtocol: "curve25519-ed25519",
		Authenticated:       true,
		SupportedExtensions: map[string]string{
			"key_rotation":    "supported",
			"forward_secrecy": "enabled",
		},
	}
)

// IsZero is a generic function to check if a value is the zero value for its type
func IsZero[T comparable](value T) bool {
	var zero T
	return value == zero
}
