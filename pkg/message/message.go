// Package message provides a type-safe, extensible message system for WebSocket communication.
// It implements a nested message structure that allows for message interception and transformation
// through an interceptor chain pattern.
package message

import (
	"encoding/json"
)

// Type aliases for improved readability and type safety
type (
	// Protocol identifies the message type or format
	Protocol string

	// Payload contains the serialized data of the message
	Payload json.RawMessage

	// Sender identifies the source of the message
	Sender string

	// Receiver identifies the intended recipient of the message
	Receiver string

	// Version specifies the message protocol version
	Version string
)

// Protocol constants
const (
	// NoneProtocol indicates no nested message exists
	NoneProtocol Protocol = "none"

	// Version1 is the current message protocol version
	Version1 Version = "v1.0"

	// UnknownReceiver receiverID initialising
	UnknownReceiver Receiver = "unknown"
)

// Message defines the interface that all message types must implement.
// It provides methods for protocol identification, serialization, and
// message nesting/unwrapping.
type Message interface {
	// GetProtocol returns the protocol identifier for this message
	GetProtocol() Protocol

	// GetNext retrieves the next message in the chain, if one exists
	// Returns nil, nil if there is no next message
	GetNext(Registry) (Message, error)

	// Marshal serializes the message to JSON format
	Marshal() ([]byte, error)

	// Unmarshal deserializes the message from JSON format
	Unmarshal([]byte) error
}

// Header contains common metadata for all messages
type Header struct {
	Sender   Sender   `json:"sender"`   // Sender identifies the message source
	Receiver Receiver `json:"receiver"` // Receiver identifies the intended recipient
	Version  Version  `json:"version"`  // Version specifies the protocol version
}

// NewV1Header creates a new header with Version1
// This is a convenience constructor for common header creation
func NewV1Header(sender Sender, receiver Receiver) Header {
	return Header{
		Sender:   sender,
		Receiver: receiver,
		Version:  Version1,
	}
}

// BaseMessage provides a foundation for all message types.
// It implements the Message interface and manages message nesting.
// Custom message types should embed this struct to inherit its functionality.
type BaseMessage struct {
	// CURRENT MESSAGE PROCESSOR
	CurrentProtocol Protocol `json:"protocol"` // CurrentProtocol identifies this message's type
	CurrentHeader   Header   `json:"header"`   // CurrentHeader contains metadata for this message
	// CURRENT OTHER FIELDS...

	// NEXT MESSAGE PROCESSOR
	NextPayload  json.RawMessage `json:"next,omitempty"` // NextPayload contains the serialized next message in the chain
	NextProtocol Protocol        `json:"next_protocol"`  // NextProtocol identifies the type of the next message. NoneProtocol indicates end of chain
}

// GetProtocol returns this message's protocol identifier
func (m *BaseMessage) GetProtocol() Protocol {
	return m.CurrentProtocol
}

// GetNext retrieves the next message in the chain, if one exists.
// Returns nil, nil if NextProtocol is NoneProtocol.
// Uses the provided Registry to create and unmarshal the next message.
func (m *BaseMessage) GetNext(registry Registry) (Message, error) {
	if m.NextProtocol == NoneProtocol {
		return nil, nil
	}

	if m.NextPayload == nil {
		return nil, ErrNoPayload
	}

	return registry.Unmarshal(m.NextProtocol, m.NextPayload)
}

// Marshal serializes the message to JSON format
func (m *BaseMessage) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Unmarshal deserializes the message from JSON format
func (m *BaseMessage) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}
