package encrypt

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/config"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptor"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/keyexchange"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/keyprovider"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/state"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type Interceptor struct {
	interceptor.NoOpInterceptor
	localMessageRegistry message.Registry
	nonceValidator       NonceValidator
	keyExchangeManager   interfaces.KeyExchangeManager
	keyProvider          keyprovider.KeyProvider
	stateManager         interfaces.StateManager
	config               config.Config
}

func (i *Interceptor) BindSocketConnection(connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) (interceptor.Writer, interceptor.Reader, error) {
	ctx, cancel := context.WithCancel(i.Ctx())

	newState, err := state.NewState(ctx, cancel, i.config, connection, writer, reader)
	if err != nil {
		return nil, nil, err
	}

	if err := i.stateManager.SetState(connection, newState); err != nil {
		return nil, nil, err
	}

	return writer, reader, nil
}

func (i *Interceptor) Init(connection interceptor.Connection) error {
	s, err := i.stateManager.GetState(connection)
	if err != nil {
		return err
	}

	if err := i.keyExchangeManager.Init(s, keyexchange.WithKeySignature(i.keyProvider)); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(i.Ctx(), 10*time.Second)
	defer cancel()

	waiter := keyexchange.NewSessionStateTargetWaiter(ctx, types.SessionStateCompleted)

	if err := i.Process(waiter, s); err != nil {
		return err
	}

	return i.keyExchangeManager.Finalise(s)
}

func (i *Interceptor) InterceptSocketWriter(writer interceptor.Writer) interceptor.Writer {
	return interceptor.WriterFunc(func(conn interceptor.Connection, msg message.Message) error {
		var m *encryptor.EncryptedMessage

		if !i.localMessageRegistry.Check(msg.GetProtocol()) {
			m, err := encryptor.NewEncryptedMessage(msg)
			if err != nil {
				return err
			}

			msg = m
		}

		if err := m.WriteProcess(i, conn); err != nil {
			return err
		}

		return writer.Write(conn, m)
	})
}

func (i *Interceptor) InterceptSocketReader(reader interceptor.Reader) interceptor.Reader {
	return interceptor.ReaderFunc(func(conn interceptor.Connection) (message.Message, error) {
		msg, err := reader.Read(conn)
		if err != nil {
			return msg, err
		}

		if !i.localMessageRegistry.Check(msg.GetProtocol()) {
			if !i.config.RequireEncryption {
				return msg, nil
			}
			return nil, encryptionerr.ErrInvalidInterceptor
		}

		m, ok := msg.(interceptor.Message)
		if !ok {
			return nil, encryptionerr.ErrInvalidInterceptor
		}

		if err := m.ReadProcess(i, conn); err != nil {
			return nil, err
		}

		return m.GetNext(i.GetMessageRegistry())
	})
}

func (i *Interceptor) UnBindSocketConnection(connection interceptor.Connection) {
	// TODO: Implement full closing
}

func (i *Interceptor) Close() error {
	// TODO: Use UnBindSocketConnection to close all
	// TODO: Close interceptor
	return nil
}

func (i *Interceptor) GetState(connection interceptor.Connection) (interfaces.State, error) {
	return i.stateManager.GetState(connection)
}

func (i *Interceptor) Process(msg interfaces.CanProcess, state interfaces.State) error {
	processor, ok := i.keyExchangeManager.(interfaces.ProtocolProcessor)
	if !ok {
		return encryptionerr.ErrInvalidMessageType
	}

	return processor.Process(msg, state)
}
