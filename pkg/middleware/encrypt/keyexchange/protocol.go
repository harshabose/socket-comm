package keyexchange

import (
	"fmt"

	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/keyprovider"
)

func WithKeySignature(keyProvider keyprovider.KeyProvider) interfaces.ProtocolFactoryOption {
	return func(protocol interfaces.Protocol) error {
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

type ProtocolFactory func(options ...interfaces.ProtocolFactoryOption) (interfaces.Protocol, error)
