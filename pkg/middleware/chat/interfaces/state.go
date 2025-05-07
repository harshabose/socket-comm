package interfaces

import (
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
)

type CanGetState interface {
	GetState(interceptor.Connection) (*state.State, error)
}

type CanSetState interface {
	SetState(interceptor.Connection, *state.State) error
}

type CanRemoveState interface {
	RemoveState(connection interceptor.Connection) error
}

type CanWriteMessage interface {
	Write(message.Message) error
}
