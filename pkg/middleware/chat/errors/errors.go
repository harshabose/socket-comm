package errors

import "errors"

var (
	ErrContextCancelled        = errors.New("context cancelled")
	ErrInterfaceMisMatch       = errors.New("unsatisfied interface triggered")
	ErrMessageForServerOnly    = errors.New("message should only be processed by server")
	ErrMessageForClientOnly    = errors.New("message should only be processed by client")
	ErrClientIDNotConsistent   = errors.New("client id is not consistent throughout the connection")
	ErrProcessExecutionStopped = errors.New("process execution stopped manually")

	ErrConnectionNotFound = errors.New("connection not registered")
	ErrConnectionExists   = errors.New("connection already exists")
	ErrInvalidInterceptor = errors.New("inappropriate interceptor for the payload")
	ErrRoomNotFound       = errors.New("room does not exists")
	ErrRoomAlreadyExists  = errors.New("room already exists")

	ErrUnknownClientIDState       = errors.New("client id not known at the moment")
	ErrClientNotAllowed           = errors.New("client is not allowed in the room")
	ErrClientIsAlreadyParticipant = errors.New("client is already a participant in the room")
	ErrClientNotAParticipant      = errors.New("client is not a participant in the room at the moment")
	ErrWrongRoom                  = errors.New("operation not permitted as room id did not match")
)

func New(text string) error {
	return errors.New(text)
}
