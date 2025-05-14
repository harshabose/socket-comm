package chat

import (
	"context"

	"github.com/harshabose/socket-comm/pkg/interceptor"
	"github.com/harshabose/socket-comm/pkg/message"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/errors"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/interfaces"
	"github.com/harshabose/socket-comm/pkg/middleware/chat/state"
)

type commonInterceptor struct {
	interceptor.NoOpInterceptor
	readProcessMessages  message.Registry
	writeProcessMessages message.Registry
	states               *state.Manager
}

func (i *commonInterceptor) BindSocketConnection(connection interceptor.Connection, writer interceptor.Writer, reader interceptor.Reader) (interceptor.Writer, interceptor.Reader, error) {
	ctx, cancel := context.WithCancel(i.Ctx())

	if err := i.states.SetState(connection, state.NewState(ctx, cancel, connection, writer, reader)); err != nil {
		cancel()
		return nil, nil, err
	}

	return writer, reader, nil
}

func (i *commonInterceptor) InterceptSocketWriter(writer interceptor.Writer) interceptor.Writer {
	return interceptor.WriterFunc(func(ctx context.Context, connection interceptor.Connection, msg message.Message) error {
		if msg == nil {
			return nil
		}

		if !i.writeProcessMessages.Check(msg.GetProtocol()) {
			return writer.Write(ctx, connection, msg)
		}

		m, ok := msg.(interceptor.Message)
		if !ok {
			return writer.Write(ctx, connection, msg)
		}

		next, err := m.GetNext(i.writeProcessMessages)
		if err != nil {
			return writer.Write(ctx, connection, msg)
		}

		if err := m.WriteProcess(ctx, i, connection); err != nil {
			return writer.Write(ctx, connection, next)
		}

		return writer.Write(ctx, connection, next)
	})
}

func (i *commonInterceptor) InterceptSocketReader(reader interceptor.Reader) interceptor.Reader {
	return interceptor.ReaderFunc(func(ctx context.Context, connection interceptor.Connection) (message.Message, error) {
		msg, err := reader.Read(ctx, connection)
		if err != nil {
			return msg, err
		}

		if msg == nil {
			return nil, nil
		}

		if !i.readProcessMessages.Check(msg.GetProtocol()) {
			return msg, nil
		}

		m, ok := msg.(interceptor.Message)
		if !ok {
			return msg, errors.ErrInterfaceMisMatch
		}

		next, err := m.GetNext(i.readProcessMessages)
		if err != nil {
			return msg, nil
		}

		if err := m.ReadProcess(ctx, i, connection); err != nil {
			return next, nil
		}

		return next, nil
	})
}

func (i *commonInterceptor) UnBindSocketConnection(connection interceptor.Connection) {

}

func (i *commonInterceptor) Close() error {
	// TODO: FIGURE OUR GOOD CLOSE STRATEGY
	return nil
}

func (i *commonInterceptor) GetState(connection interceptor.Connection) (*state.State, error) {
	return i.states.GetState(connection)
}

func (i *commonInterceptor) Process(ctx context.Context, process interfaces.CanBeProcessed, state *state.State) error {
	return process.Process(ctx, i, state)
}

func (i *commonInterceptor) ProcessBackground(ctx context.Context, process interfaces.CanBeProcessedBackground, state *state.State) interfaces.CanBeProcessedBackground {
	return process.ProcessBackground(ctx, i, state)
}
