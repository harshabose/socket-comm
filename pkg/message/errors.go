// Package message provides a type-safe, extensible message system for WebSocket communication.
package message

import "errors"

// Error definitions for the message package
var (
	// ErrNoProtocolMatch is returned when attempting to create or unmarshal
	// a message with a protocol that is not registered
	ErrNoProtocolMatch = errors.New("no protocol in the registry")

	// ErrNoPayload is returned when a message has a non-none NextProtocol
	// but is missing the corresponding NextPayload data
	ErrNoPayload = errors.New("protocol is not none but payload is nil")

	// ErrInvalidMessageData is returned when raw message data cannot be
	// properly identified or does not contain a valid protocol field
	ErrInvalidMessageData = errors.New("invalid message data")
)
