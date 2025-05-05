package errors

import "errors"

var (
	ErrInterfaceMisMatch    = errors.New("unsatisfied interface triggered")
	ErrMessageForServerOnly = errors.New("message should only be processed by server")
	ErrMessageForClientOnly = errors.New("message should only be processed by client")

	ErrConnectionNotFound = errors.New("connection not registered")
	ErrConnectionExists   = errors.New("connection already exists")
	ErrInvalidInterceptor = errors.New("inappropriate interceptor for the payload")
	ErrRoomNotFound       = errors.New("room does not exists")
	ErrRoomAlreadyExists  = errors.New("room already exists")

	ErrUnknownClientIDState       = errors.New("client ID not known at the moment")
	ErrClientNotAllowedInRoom     = errors.New("client is not allowed in the room")
	ErrClientIsAlreadyParticipant = errors.New("client is already a participant in the room")
	ErrClientNotAParticipant      = errors.New("client is not a participant in the room at the moment")
	ErrWrongRoom                  = errors.New("operation not permitted as room id did not match")
)
