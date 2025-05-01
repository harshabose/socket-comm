package encrypt

import (
	"context"
	"time"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/config"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/encryptionerr"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/keyexchange"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/keyprovider"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/state"
	"github.com/harshabose/socket-comm/pkg/middleware/encrypt/types"
)

type Interceptor struct {
	interceptor.NoOpInterceptor
	nonceValidator     NonceValidator
	keyExchangeManager interfaces.KeyExchangeManager
	keyProvider        keyprovider.KeyProvider
	stateManager       interfaces.StateManager
	config             config.Config
}

func (i *Interceptor) BindSocketConnection(connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) (interceptor.Writer, interceptor.Reader, error) {
	ctx, cancel := context.WithCancel(i.Ctx)

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

	ctx, cancel := context.WithTimeout(i.Ctx, 10*time.Second)
	defer cancel()

	waiter := keyexchange.NewSessionStateTargetWaiter(ctx, types.SessionStateCompleted)

	if err := i.Process(waiter, s); err != nil {
		return err
	}

	return i.keyExchangeManager.Finalise(s)
}

func (i *Interceptor) InterceptSocketWriter(writer interceptor.Writer) interceptor.Writer {
	return interceptor.WriterFunc(func(conn interceptor.Connection, msg message.Message) error {
		s, err := i.GetState(conn)
		if err != nil {
			return err
		}

		ss, ok := s.(interfaces.CanEncrypt)
		if !ok {
			return err
		}

		encrypted, err := ss.Encrypt(msg)
		if err != nil {
			if !s.GetConfig().RequireEncryption {
				return writer.Write(conn, msg)
			}
			return err
		}

		return writer.Write(conn, encrypted)
	})
}

func (i *Interceptor) InterceptSocketReader(reader interceptor.Reader) interceptor.Reader {
	return interceptor.ReaderFunc(func(conn interceptor.Connection) (message.Message, error) {
		msg, err := reader.Read(conn)
		if err != nil {
			return msg, err
		}

		s, err := i.GetState(conn)
		if err != nil {
			return nil, err
		}

		ss, ok := s.(interfaces.CanDecrypt)
		if !ok {
			return msg, err
		}

		m, err := ss.Decrypt(msg)
		if err != nil {
			return nil, err
		}

		return m, nil
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
