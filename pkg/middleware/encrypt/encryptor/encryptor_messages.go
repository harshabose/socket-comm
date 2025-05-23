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

func NewEncryptedMessage(msg message.Message) (*EncryptedMessage, error) {
	em := &EncryptedMessage{
		Timestamp: time.Now(),
	}

	bmsg, err := interceptor.NewBaseMessage(msg.GetProtocol(), msg, em)
	if err != nil {
		return nil, err
	}

	em.BaseMessage = bmsg

	return em, nil
}

func (m *EncryptedMessage) WriteProcess(_i interceptor.Interceptor, connection interceptor.Connection) error {
	i, ok := _i.(interfaces.CanGetState)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	s, err := i.GetState(connection)
	if err != nil {
		return err
	}

	ss, ok := s.(interfaces.CanEncrypt)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	encmsg, err := ss.Encrypt(m)
	if err != nil {
		return err
	}

	msg, ok := encmsg.(*EncryptedMessage)
	if !ok {
		return encryptionerr.ErrInvalidInterceptor
	}

	m.NextPayload = msg.NextPayload

	return nil
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
