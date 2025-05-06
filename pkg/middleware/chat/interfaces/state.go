package interfaces

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type CanSetClientID interface {
	SetClientID(id types.ClientID) error
}

type State interface {
	GetClientID() (types.ClientID, error)
	Ctx() context.Context
}

type CanGetState interface {
	GetState(interceptor.Connection) (State, error)
}

type CanSetState interface {
	SetState(interceptor.Connection, State) error
}

type CanRemoveState interface {
	RemoveState(connection interceptor.Connection) error
}

type CanWriteMessage interface {
	Write(message.Message) error
}

type StateManager interface {
	CanGetState
	CanSetState
	CanRemoveState
	ForEach(func(interceptor.Connection, State) error) error
}
