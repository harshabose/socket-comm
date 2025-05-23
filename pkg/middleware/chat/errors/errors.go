package errors

import (
	"github.com/harshabose/socket-comm/pkg/interceptor"
)

var (
	ErrRoomNotFound      = interceptor.NewError("room does not exists")
	ErrRoomAlreadyExists = interceptor.NewError("room already exists")

	ErrUnknownClientIDState       = interceptor.NewError("client id not known at the moment")
	ErrClientNotAllowed           = interceptor.NewError("client is not allowed in the room")
	ErrClientIsAlreadyParticipant = interceptor.NewError("client is already a participant in the room")
	ErrClientNotAParticipant      = interceptor.NewError("client is not a participant in the room at the moment")
	ErrWrongRoom                  = interceptor.NewError("operation not permitted as room id did not match")
)
