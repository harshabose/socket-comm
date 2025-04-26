package keyexchange

import (
	"fmt"

	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/keyprovider"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/state"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type SessionState int

const (
	SessionStateInitial SessionState = iota
	SessionStateInProgress
	SessionStateCompleted
	SessionStateError
)

type MessageProcessor interface {
	Process(Protocol, *state.State) error
}

type Protocol interface {
	Init(s *state.State) error
	GetKeys() (encKey types.Key, decKey types.Key, err error)
	Process(MessageProcessor, *state.State) error
	GetState() SessionState
	IsComplete() bool
}

type ProtocolFactoryOption func(Protocol) error

func WithKeySignature(keyProvider keyprovider.KeyProvider) ProtocolFactoryOption {
	return func(protocol Protocol) error {
		curveProtocol, ok := protocol.(*Curve25519Protocol)
		if !ok {
			return fmt.Errorf("WithKeySignature only supports Curve25519Protocol")
		}

		curveProtocol.options.SigningKey = keyProvider.GetSigningKey()
		curveProtocol.options.VerificationKey = keyProvider.GetVerificationKey()
		curveProtocol.options.RequireSignature = true

		return nil
	}
}

type ProtocolFactory func(options ...ProtocolFactoryOption) (Protocol, error)
