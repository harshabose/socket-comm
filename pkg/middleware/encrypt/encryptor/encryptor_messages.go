package encryptor

import (
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type EncryptedMessage struct {
	interceptor.BaseMessage
	Nonce     types.Nonce
	Timestamp time.Time
	SessionID types.EncryptionSessionID
}

func NewEncryptedMessage(encrypted message.Payload, protocol message.Protocol, nonce types.Nonce, id types.EncryptionSessionID) (*EncryptedMessage, error) {
	msg := &EncryptedMessage{
		Nonce:     nonce,
		Timestamp: time.Now(),
		SessionID: id,
	}

	bmsg, err := interceptor.NewBaseMessage(protocol, encrypted, msg)
	if err != nil {
		return nil, err
	}

	msg.BaseMessage = bmsg

	return msg, nil
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

	decmsg, err := ss.Decrypt(m)
	if err != nil {
		return err
	}

	msg, ok := decmsg.(*EncryptedMessage)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor // JUST TO BE SURE
	}

	m.NextPayload = msg.NextPayload // JUST MAKING SURE

	return nil
}
