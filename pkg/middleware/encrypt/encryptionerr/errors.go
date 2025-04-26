package encryptionerr

import "errors"

// Common error definitions for the encryption interceptor
var (
	// General errors
	ErrConnectionNotFound = errors.New("connection not registered")
	ErrConnectionExists   = errors.New("connection already exists")
	ErrInvalidInterceptor = errors.New("inappropriate interceptor for the payload")

	// Key exchange errors
	ErrInvalidMessageType  = errors.New("message does not implement required key exchange MessageProcessor interface")
	ErrKeyExchangeTimeout  = errors.New("key exchange timed out")
	ErrInvalidSignature    = errors.New("signature verification failed")
	ErrProtocolNotFound    = errors.New("key exchange protocol not found")
	ErrExchangeInProgress  = errors.New("key exchange already in progress")
	ErrSessionNotFound     = errors.New("key exchange session not found")
	ErrInvalidSessionState = errors.New("key exchange state is not valid")
	ErrExchangeNotComplete = errors.New("key exchange not completed")

	// Encryption errors
	ErrEncryptionNotReady = errors.New("encryption not ready")
	ErrInvalidKey         = errors.New("invalid encryption key")
	ErrInvalidNonce       = errors.New("invalid nonce")
	ErrNonceReused        = errors.New("nonce has been used before")

	// Configuration errors
	ErrInvalidConfig   = errors.New("invalid configuration")
	ErrInvalidProvider = errors.New("invalid key provider")

	// Security errors
	ErrInvalidServerRequest = errors.New("invalid request to server")
	ErrSecurityViolation    = errors.New("security violation detected")
)
