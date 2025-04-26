// Package interceptor provides a middleware system for processing WebSocket messages.
// It builds on the message package to add processing capabilities to messages.
package interceptor

import "github.com/harshabose/socket-comm/pkg/message"

// Message extends the base message.Message interface with processing capabilities.
// This interface allows messages to be processed by interceptors in the communication chain.
// Types implementing this interface can define custom behavior for how they interact
// with specific interceptors.
type Message interface {
	// Message Embed the base Message interface
	message.Message

	// WriteProcess handles interceptor processing for outgoing messages.
	// This method is called when a message is being written to a connection.
	//
	// The implementation should handle any message-specific processing required
	// for the given interceptor type. For example, an encryption message would
	// encrypt data when this method is called.
	//
	// Parameters:
	//   - interceptor: The interceptor that should process this message
	//   - connection: The network connection associated with this message
	//
	// Returns an error if processing fails
	WriteProcess(Interceptor, Connection) error

	// ReadProcess handles interceptor processing for incoming messages.
	// This method is called when a message is being read from a connection.
	//
	// The implementation should handle any message-specific processing required
	// for the given interceptor type. For example, an encryption message would
	// decrypt data when this method is called.
	//
	// Parameters:
	//   - interceptor: The interceptor that should process this message
	//   - connection: The network connection associated with this message
	//
	// Returns an error if processing fails
	ReadProcess(Interceptor, Connection) error
}

// BaseMessage provides a default implementation of the Message interface.
// It embeds message.BaseMessage to inherit its functionality and adds
// a no-op Process method that can be overridden by specific message types.
//
// Custom interceptor message types should embed this struct and override
// the Process method with their specific processing logic.
type BaseMessage struct {
	// Embed the base message implementation
	message.BaseMessage
}

// NewBaseMessage creates a properly initialized interceptor BaseMessage for the key exchange module
func NewBaseMessage(protocol message.Protocol, sender message.Sender, receiver message.Receiver) BaseMessage {
	return BaseMessage{
		BaseMessage: message.BaseMessage{
			CurrentProtocol: protocol,
			CurrentHeader:   message.NewV1Header(sender, receiver),
			NextProtocol:    message.NoneProtocol,
		},
	}
}

// WriteProcess handles interceptor processing for outgoing messages.
// This method is called when a message is being written to a connection.
// It should be overridden by specific message types to implement
// their custom outgoing message processing logic.
//
// Parameters:
//   - interceptor: The interceptor that should process this message
//   - connection: The network connection associated with this message
//
// Returns nil by default, indicating no processing was performed
func (m *BaseMessage) WriteProcess(_ Interceptor, _ Connection) error {
	// Default implementation does nothing
	// Derived-types should override this method with specific processing logic
	return nil
}

// ReadProcess handles interceptor processing for incoming messages.
// This method is called when a message is being read from a connection.
// It should be overridden by specific message types to implement
// their custom incoming message processing logic.
//
// Parameters:
//   - interceptor: The interceptor that should process this message
//   - connection: The network connection associated with this message
//
// Returns nil by default, indicating no processing was performed
func (m *BaseMessage) ReadProcess(_ Interceptor, _ Connection) error {
	// Default implementation does nothing
	// Derived-types should override this method with specific processing logic
	return nil
}
