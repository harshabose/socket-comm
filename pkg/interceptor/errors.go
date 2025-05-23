package interceptor

import "errors"

// COMMON ERRORS

var (
	ErrContextCancelled        = errors.New("context cancelled")
	ErrInterfaceMisMatch       = errors.New("unsatisfied interface triggered")
	ErrProcessExecutionStopped = errors.New("process execution stopped manually")
	ErrClientIDNotConsistent   = errors.New("client id is not consistent throughout the connection")

	ErrConnectionNotFound = errors.New("connection not registered")
	ErrConnectionExists   = errors.New("connection already exists")
	ErrInvalidInterceptor = errors.New("inappropriate interceptor for the message")
)

func NewError(text string) error {
	return errors.New(text)
}
