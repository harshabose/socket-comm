// Package message provides functionality for message handling, registration, and unmarshaling
// in a protocol-based messaging system.
package message

import (
	"encoding/json"
	"fmt"
	"sync"
)

// envelope is an internal structure used to extract the protocol information
// from a raw JSON message. It serves as a container for the protocol field
// during the unmarshaling process.
type envelope struct {
	Protocol Protocol `json:"protocol"`
}

// Registry defines the interface for a message registry system that manages
// different message protocols. It provides methods for registering message factories,
// checking protocol existence, and unmarshaling messages from different formats.
type Registry interface {
	// Register adds a new protocol to the registry.
	Register(Protocol, Factory) error

	// Check function verifies if a protocol is registered in the registry.
	Check(protocol Protocol) bool

	// Unmarshal creates and populates a message of the specified protocol using the provided payload.
	Unmarshal(Protocol, Payload) (Message, error)

	// UnmarshalRaw extracts the protocol from a raw JSON message and unmarshals it accordingly.
	UnmarshalRaw(Payload) (Message, error)
}

// Factory defines an interface for creating new message instances.
// Implementations of this interface are responsible for instantiating
// specific message types based on their protocol.
type Factory interface {
	// Create instantiates a new message instance.
	Create() (Message, error)
}

// EmptyFactoryFunc is a function type that implements the Factory interface.
// It allows simple functions to be used as factories without creating a dedicated struct.
type EmptyFactoryFunc func() (Message, error)

// Create implements the Factory interface for EmptyFactoryFunc.
// It simply calls the function itself to create a new message instance.
func (f EmptyFactoryFunc) Create() (Message, error) {
	return f()
}

// DefaultRegistry is the standard implementation of the Registry interface.
// It maintains a thread-safe map of protocol-to-factory mappings.
type DefaultRegistry struct {
	// factories field maps protocol identifiers to their corresponding message factories
	factories map[Protocol]Factory
	// mux provides thread-safety for concurrent access to the factories map
	mux sync.RWMutex
}

// NewDefaultRegistry creates and initializes a new DefaultRegistry instance.
// It initializes the factories map and returns a ready-to-use registry.
func NewDefaultRegistry() *DefaultRegistry {
	return &DefaultRegistry{
		factories: make(map[Protocol]Factory),
	}
}

// Register adds a new protocol and its associated factory to the registry.
// It returns an error if the protocol is already registered.
// This method is thread-safe.
func (r *DefaultRegistry) Register(protocol Protocol, factory Factory) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if _, exists := r.factories[protocol]; exists {
		return fmt.Errorf("protocol %s is already registered", protocol)
	}

	r.factories[protocol] = factory
	return nil
}

// Check verifies if a protocol is registered in the registry.
// Returns true if the protocol exists, false otherwise.
// This method is thread-safe and uses a read lock for better concurrency.
func (r *DefaultRegistry) Check(protocol Protocol) bool {
	r.mux.RLock()
	defer r.mux.RUnlock()

	_, exists := r.factories[protocol]
	return exists
}

// Unmarshal creates a new message instance for the specified protocol and
// populates it with the provided data. It returns an error if:
// - The protocol is not registered
// - The factory fails to create a message instance
// - The message fails to unmarshal the data
// This method is thread-safe.
func (r *DefaultRegistry) Unmarshal(protocol Protocol, data Payload) (Message, error) {
	if !r.Check(protocol) {
		return nil, ErrNoProtocolMatch
	}

	msg, err := r.factories[protocol].Create()
	if err != nil {
		return nil, err
	}

	if err := msg.Unmarshal(data); err != nil {
		return nil, err
	}

	return msg, nil
}

// UnmarshalRaw extracts the protocol from a raw JSON message and then
// unmarshal the message using the appropriate factory. It returns an error if:
// - The JSON cannot be parsed
// - The protocol field is missing or empty
// - The Unmarshal method fails
// This method is particularly useful when receiving messages from external sources
// where the protocol is not known in advance.
func (r *DefaultRegistry) UnmarshalRaw(data Payload) (Message, error) {
	var envelope envelope
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("failed to extract protocol: %w", err)
	}

	if envelope.Protocol == "" {
		return nil, ErrInvalidMessageData
	}

	return r.Unmarshal(envelope.Protocol, data)
}
