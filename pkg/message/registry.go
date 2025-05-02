// Package message provides a type-safe, extensible message system for WebSocket communication.
package message

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Envelope is a lightweight struct used to extract just the protocol
// information from a raw message without full deserialization.
// This enables protocol-based routing of incoming messages.
type Envelope struct {
	Protocol Protocol `json:"protocol"`
}

// Registry defines the interface for message type registration and instantiation.
// It provides a centralized mechanism for creating, deserializing, and inspecting messages.
type Registry interface {
	// Register adds a message factory for a specific protocol
	// Returns an error if the protocol is already registered
	Register(Protocol, Factory) error

	Check(protocol Protocol) bool

	// Create instantiates a new message for the given protocol
	// Returns an error if the protocol is not registered
	Create(Protocol) (Message, error)

	// Unmarshal creates and deserializes a message for a protocol
	// The provided data is parsed into the appropriate message type
	Unmarshal(Protocol, Payload) (Message, error)

	// UnmarshalRaw deserializes a message when the protocol is unknown
	// It first inspects the envelope to determine the protocol, then unmarshals accordingly
	UnmarshalRaw(json.RawMessage) (Message, error)
}

// Factory defines the interface for creating new message instances.
// Each message type should have a corresponding factory.
type Factory interface {
	// Create instantiates a new instance of a message type
	Create() (Message, error)
}

// FactoryFunc is a function type that implements the Factory interface.
// It allows simple functions to be used as factories without creating a separate type.
type FactoryFunc func() Message

// Create implements the Factory interface for FactoryFunc
func (f FactoryFunc) Create() Message {
	return f()
}

// DefaultRegistry provides a thread-safe implementation of the Registry interface.
// It maintains a map of protocols to their corresponding factories.
type DefaultRegistry struct {
	factories map[Protocol]Factory
	mux       sync.RWMutex // Protects concurrent access to factories
}

// NewRegistry creates a new message registry
// The returned registry is ready to use but contains no registered message types
func NewRegistry() *DefaultRegistry {
	return &DefaultRegistry{
		factories: make(map[Protocol]Factory),
	}
}

func (r *DefaultRegistry) Check(protocol Protocol) bool {
	r.mux.RLock()
	defer r.mux.RUnlock()

	_, exists := r.factories[protocol]
	return exists
}

// Register adds a message factory for a protocol
// Returns an error if the protocol is already registered
// This method is thread-safe
func (r *DefaultRegistry) Register(protocol Protocol, factory Factory) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, exists := r.factories[protocol]; exists {
		return fmt.Errorf("protocol %s is already registered", protocol)
	}

	r.factories[protocol] = factory
	return nil
}

// Create instantiates a new message for a protocol
// Returns an error if the protocol is not registered
// This method is thread-safe
func (r *DefaultRegistry) Create(protocol Protocol) (Message, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	factory, exists := r.factories[protocol]
	if !exists {
		return nil, ErrNoProtocolMatch
	}

	msg, err := factory.Create()
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// Unmarshal creates and deserializes a message for a protocol
// The provided data is parsed into the appropriate message type
// This method is thread-safe
func (r *DefaultRegistry) Unmarshal(protocol Protocol, data Payload) (Message, error) {
	msg, err := r.Create(protocol)
	if err != nil {
		return nil, err
	}

	if err := msg.Unmarshal(data); err != nil {
		return nil, err
	}

	return msg, nil
}

// UnmarshalRaw deserializes a message when the protocol is unknown
// It first extracts just the protocol from the data, then creates and unmarshals the appropriate message type
// This method is particularly useful for handling incoming WebSocket messages
// This method is thread-safe
func (r *DefaultRegistry) UnmarshalRaw(data Payload) (Message, error) {
	var envelope Envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("failed to extract protocol: %w", err)
	}

	if envelope.Protocol == "" {
		return nil, ErrInvalidMessageData
	}

	return r.Unmarshal(envelope.Protocol, data)
}
