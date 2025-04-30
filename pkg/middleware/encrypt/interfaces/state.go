package interfaces

import (
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/config"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type ProtocolFactoryOption func(Protocol) error

type State interface {
	GenerateKeyExchangeSessionID() types.KeyExchangeSessionID
	GetKeyExchangeSessionID() types.KeyExchangeSessionID
	WriteMessage(interceptor.Message) error
	GetConfig() config.Config
}

type CanGetState interface {
	GetState(interceptor.Connection) (State, error)
}

type CanSetState interface {
	SetState(interceptor.Connection, State) error
}

type CanRemoveState interface {
	RemoveState(interceptor.Connection) error
}

type StateManager interface {
	CanGetState
	CanSetState
	CanRemoveState
	ForEach(func(interceptor.Connection, State) error) error
}
