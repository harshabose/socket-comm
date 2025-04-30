package interfaces

import (
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type ProtocolProcessor interface {
	Process(CanProcess, State) error
}

type KeyExchangeManager interface {
	Init(state State, options ...ProtocolFactoryOption) error
}

type CanGetSessionState interface {
	GetState() types.SessionState
}

type Protocol interface {
	Init(s State) error
	IsComplete() bool
}

type CanProcess interface {
	Process(Protocol, State) error
}
