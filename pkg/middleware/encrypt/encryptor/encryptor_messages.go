package encryptor

import (
	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
)

type EncryptedMessage struct {
	interceptor.BaseMessage
}

func (m *EncryptedMessage) ReadProcess(_i interceptor.Interceptor, conn interceptor.Connection) error {
	i, ok := _i.(interfaces.CanGetState)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	s, err := i.GetState(conn)
	if err != nil {
		return err
	}

	ss, ok := s.(interfaces.CanDecrypt)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	msg, err := ss.Decrypt(m)
	if err != nil {
		return err
	}

	// TODO: message.Registry is not implemented yet
	decrytpedMsg, err := message.Registry().Unmarshal(m.NextProtocol, m.NextPayload)
	if err != nil {
		return err
	}

	return nil
}
