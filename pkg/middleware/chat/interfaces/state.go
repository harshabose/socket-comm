package interfaces

import (
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/config"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/types"
)

type State interface {
	GetClientID() (types.ClientID, error)
	GetConfig() config.Config
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
